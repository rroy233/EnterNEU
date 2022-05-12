package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"github.com/rroy233/EnterNEU/logger"
	"runtime"

	"github.com/rroy233/EnterNEU/configs"
)

func AesCbcEncryptString(text string, key string) (string, error) {
	defer func() {
		if err := recover(); err != nil {
			switch err.(type) {
			case runtime.Error:
				logger.Info.Println("runtime err:", err, "Check that the key or text is correct")
			default:
				logger.Info.Println("error:", err)
			}
		}
	}()
	if len(key) != 32 {
		return "", errors.New("ErrKeyLength")
	}
	aesKey := []byte(key)
	plainText := []byte(text)
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return "", err
	}
	paddingText := pkcs5Padding(plainText, block.BlockSize())

	var iv []byte
	iv = []byte(configs.Get().General.AesIv)

	blockMode := cipher.NewCBCEncrypter(block, iv)
	cipherText := make([]byte, len(paddingText))
	blockMode.CryptBlocks(cipherText, paddingText)

	out := base64.URLEncoding.EncodeToString(cipherText)
	return out, nil
}

func AesCbcEncrypt(plainText []byte, key string) ([]byte, error) {
	defer func() {
		if err := recover(); err != nil {
			switch err.(type) {
			case runtime.Error:
				logger.Info.Println("runtime err:", err, "Check that the key or text is correct")
			default:
				logger.Info.Println("error:", err)
			}
		}
	}()
	if len(key) != 32 {
		return nil, errors.New("ErrKeyLength")
	}
	aesKey := []byte(key)
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}
	paddingText := pkcs5Padding(plainText, block.BlockSize())

	var iv []byte
	iv = []byte(configs.Get().General.AesIv)

	blockMode := cipher.NewCBCEncrypter(block, iv)
	cipherData := make([]byte, len(paddingText))
	blockMode.CryptBlocks(cipherData, paddingText)

	return cipherData, nil
}

func AesCbcDecryptString(text string, key string) (string, error) {
	if len(key) != 32 {
		return "", errors.New("ErrKeyLength")
	}
	aesKey := []byte(key)
	cipherText, err := base64.URLEncoding.DecodeString(text)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return "", err
	}

	defer func() {
		if err := recover(); err != nil {
			switch err.(type) {
			case runtime.Error:
				logger.Info.Println("runtime err:", err, "Check that the key or text is correct")
			default:
				logger.Info.Println("error:", err)
			}
		}
	}()
	var iv []byte
	iv = []byte(configs.Get().General.AesIv)
	blockMode := cipher.NewCBCDecrypter(block, iv)
	paddingText := make([]byte, len(cipherText))
	blockMode.CryptBlocks(paddingText, cipherText)

	plainText, err := pkcs5UnPadding(paddingText)
	if err != nil {
		return "", err
	}

	return string(plainText), nil
}

func AesCbcDecrypt(data []byte, key string) ([]byte, error) {
	if len(key) != 32 {
		return nil, errors.New("ErrKeyLength")
	}
	aesKey := []byte(key)

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := recover(); err != nil {
			switch err.(type) {
			case runtime.Error:
				logger.Info.Println("runtime err:", err, "Check that the key or text is correct")
			default:
				logger.Info.Println("error:", err)
			}
		}
	}()
	var iv []byte
	iv = []byte(configs.Get().General.AesIv)
	blockMode := cipher.NewCBCDecrypter(block, iv)
	paddingText := make([]byte, len(data))
	blockMode.CryptBlocks(paddingText, data)

	originalData, err := pkcs5UnPadding(paddingText)
	if err != nil {
		return nil, err
	}

	return originalData, nil
}

func pkcs5Padding(plainText []byte, blockSize int) []byte {
	padding := blockSize - (len(plainText) % blockSize)
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	newText := append(plainText, padText...)
	return newText
}

func pkcs5UnPadding(plainText []byte) ([]byte, error) {
	length := len(plainText)
	number := int(plainText[length-1])
	if number > length {
		return nil, errors.New("ErrPaddingSize")
	}
	return plainText[:length-number], nil
}
