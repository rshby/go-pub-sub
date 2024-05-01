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

	chiperText := `ff5961672b325d04ed28ddbebe1fdd19373f56d84498fa0f2a5fdaf1c013a7f37199eaef743f8a8ef158054b7dde8fd83a333d1ea8d2f0e1c0f87cc5cdcd44d810143a948c312560a9de60d10b9a965a2716e4625615ebe04e95297f3db8499a4eafdd149bf044e2f97cb9331bb44a149403fd98c98aca0411fbb090458198c8eef8313ee7307e7e7b08cb3dec360a44f9b339966f51ec6e273c0184725993022eaf989ca55776be3071a24780e88f569f5f9fc34f7610a9caabbb5ce56cddefe77c60b16043c50e46988f44d44b106f3d7e9b0470a686ea9aa14fb5b66cfa243c3d95682809bd5251ae3db876bcdb70c4a514a1ec5fa089248b008e2135058bfb795c7c39b50d04f9c2f3d63d8a9373`
	ivKey := `d81c4519c52543d224d7579bdfe657d1`

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
