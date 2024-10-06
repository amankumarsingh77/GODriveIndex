package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"net/url"
)

var (
	CryptoBaseKey = []byte("3225f86e99e205347b4310e437253bfd")
	EncryptIV     = []byte{247, 254, 106, 195, 32, 148, 131, 244, 222, 133, 26, 182, 20, 138, 215, 81}
)

func EncryptString(plaintext string) (string, error) {
	block, err := aes.NewCipher(CryptoBaseKey)
	if err != nil {
		return "", err
	}
	paddedPlaintext := pad([]byte(plaintext), aes.BlockSize)
	ciphertext := make([]byte, len(paddedPlaintext))
	mode := cipher.NewCBCEncrypter(block, EncryptIV)
	mode.CryptBlocks(ciphertext, paddedPlaintext)
	encoded := base64.RawURLEncoding.EncodeToString(ciphertext)
	return url.QueryEscape(encoded), nil
}

func DecryptString(ciphertext string) (string, error) {
	unescaped, err := url.QueryUnescape(ciphertext)
	if err != nil {
		return "", err
	}
	decodedCiphertext, err := base64.RawURLEncoding.DecodeString(unescaped)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(CryptoBaseKey)
	if err != nil {
		return "", err
	}
	mode := cipher.NewCBCDecrypter(block, EncryptIV)
	mode.CryptBlocks(decodedCiphertext, decodedCiphertext)
	return string(unpad(decodedCiphertext)), nil
}

func pad(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

func unpad(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])
	return src[:(length - unpadding)]
}

func EncryptForURL(input string) string {
	encrypted, _ := EncryptString(input)
	return encrypted
}
