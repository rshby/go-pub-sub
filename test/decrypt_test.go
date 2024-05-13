// You can edit this code!
// Click here and start typing.
package test

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"testing"
)

func TestDescript(t *testing.T) {
	privateKey := `YqiPRNISlDzhOVnvSK4QIbrVPTioz77F`

	chiperText := `8f2453276ebf2c16cb911583f6418c17c4c9aabf751e08ccd5e3caba2eed926b3dcd661661380c3d800cce49d8a56c3120fdbf2eb70a2b2a56cf789e064d36f14fe3ca1a877d4f30e4df00070683ef4bf6d06975e6646e2e4cd8e555587b0f8306b6e949b0471739540c3740bf279e4e2bf7ec1b4e9909ccf3a9c67acdf88b900144c0e96bcf39b031dcc2ecc0c17da1ad85056c4372456cba9753ff49d0696b860e68a263a3107aa863c32835925dcd109244e8b9df7b8c01006baee97194192329f24c9724b1d8454752cba17a063a5034b4ffc21c36f6966c71a06bcc57e4`
	ivKey := `cd388d2076a27f47bbf1d5efdb62ba52`

	cryptor := Cryptor(string(privateKey), ivKey)

	fmt.Println(cryptor.Decrypt(chiperText))
}

func Cryptor(privateKey, ivKey string) *AESCryptor {
	return NewAESCryptor(privateKey, ivKey, aes.BlockSize)
}

// AESCryptor aes cryptor
type AESCryptor struct {
	encryptionKey []byte
	iv            string
	blockSize     int
}

// NewAESCryptor create new AESCryptor.
// The key will be hashed using sha256 and will be used as the private key
func NewAESCryptor(key string, iv string, blockSize int) *AESCryptor {
	return &AESCryptor{
		encryptionKey: generateEncryptionKey(key),
		iv:            iv,
		blockSize:     blockSize,
	}
}

// Decrypt decrypt cipherText using CBC
func (aesCryptor *AESCryptor) Decrypt(cipherText string) (plainText string, err error) {
	ivKey, err := aesCryptor.generateIVKey(aesCryptor.iv)
	if err != nil {
		return
	}

	cipherTextDecoded, err := hex.DecodeString(cipherText)
	if err != nil {
		return
	}

	block, err := aes.NewCipher(aesCryptor.encryptionKey)
	if err != nil {
		return
	}

	mode := cipher.NewCBCDecrypter(block, ivKey)
	mode.CryptBlocks(cipherTextDecoded, cipherTextDecoded)

	return string(aesCryptor.pkcs5Unpadding(cipherTextDecoded)), nil
}

func (aesCryptor *AESCryptor) generateIVKey(iv string) (bIv []byte, err error) {
	if len(iv) > 0 {
		ivKey, err := hex.DecodeString(iv)
		if err != nil {
			return nil, errors.New("ErrDecodeIVKey")
		}

		return ivKey, nil
	}

	ivKey, err := GenerateRandomIVKey(aesCryptor.blockSize)
	if err != nil {
		return nil, errors.New("ErrGenerateIVKey")
	}

	return hex.DecodeString(ivKey)
}

func (aesCryptor *AESCryptor) pkcs5Padding(ciphertext []byte, blockSize int, after int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)

	return append(ciphertext, padtext...)
}

func (aesCryptor *AESCryptor) pkcs5Unpadding(src []byte) []byte {
	if len(src) == 0 {
		return nil
	}

	length := len(src)
	unpadding := int(src[length-1])
	cutLen := (length - unpadding)
	// check boundaries
	if cutLen < 0 || cutLen > length {
		return src
	}

	return src[:cutLen]
}

func generateEncryptionKey(key string) []byte {
	hash := sha256.New()
	hash.Write([]byte(key))
	return hash.Sum(nil)
}

// GenerateRandomIVKey generate random IV value
func GenerateRandomIVKey(blockSize int) (bIv string, err error) {
	bytes := make([]byte, blockSize)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}
