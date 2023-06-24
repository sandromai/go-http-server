package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
	"os"

	"github.com/sandromai/go-http-server/types"
)

func Encrypt(
	data string,
) (
	encryptedData string,
	appErr *types.AppError,
) {
	key := []byte(os.Getenv("ENCRYPTION_KEY"))

	block, err := aes.NewCipher(key)

	if err != nil {
		return "", &types.AppError{
			StatusCode: 500,
			Message:    "Error encrypting data.",
		}
	}

	rawData := []byte(data)
	cipherText := make([]byte, aes.BlockSize+len(rawData))
	iv := cipherText[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", &types.AppError{
			StatusCode: 500,
			Message:    "Error creating IV.",
		}
	}

	stream := cipher.NewCFBEncrypter(block, iv)

	stream.XORKeyStream(cipherText[aes.BlockSize:], rawData)

	return base64.URLEncoding.EncodeToString(cipherText), nil
}

func Decrypt(
	encryptedData string,
) (
	decryptedData string,
	appErr *types.AppError,
) {
	key := []byte(os.Getenv("ENCRYPTION_KEY"))

	block, err := aes.NewCipher(key)

	if err != nil {
		return "", &types.AppError{
			StatusCode: 500,
			Message:    "Error decrypting data.",
		}
	}

	cipherText, err := base64.URLEncoding.DecodeString(encryptedData)

	if err != nil {
		return "", &types.AppError{
			StatusCode: 500,
			Message:    "Could not decrypt data.",
		}
	}

	if len(cipherText) < aes.BlockSize {
		return "", &types.AppError{
			StatusCode: 500,
			Message:    "Invalid encoded data.",
		}
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	stream.XORKeyStream(cipherText, cipherText)

	return string(cipherText), nil
}
