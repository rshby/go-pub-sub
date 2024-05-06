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

	chiperText := `5b08fee17686a2d3703ab8e3260f17019a323772dcc852d5ac1b5089249f76db5f23b6b3a15ea4b84f1012e16f4a4d5469907ab007c1cc836dfe48b3fbe4febf23bb5c9e9a2c4560020641a30d1f93dbd98e1417a19c6a9a4451252542f2f0a46c9e9b39eaf533271ec6d940e393c1380c4ab6d4f5747684bf405a8d6a99a36bbea899a802c1e65cdd98c1c8a9f93655259a282d8381809e82198cda0e7278287d4af86590554cfd4168f8d0d7f7b362793b36ccb52da8b8ff6fd8c27d6f382535dd303823a8dc1ba9ff7555de715e1b2c500672288fba7429fc04061906eaa3`
	ivKey := `bcb8c7643a5dd1f75c719b32078f7ab6`

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
