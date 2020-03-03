package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"godev/strutils"
	"io"
)

// Encrypt encrypts plaintext and encode it into base64 string
//	key length need to be 16bytes (AES-128), 24bytes(AES-192) or 32bytes(AES-256)
//	default is 16bytes (blockSize)
func Encrypt(plaintext, key []byte) (string, error) {
	data, err := aesCBCEncrypt(plaintext, key)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

// Decrypt decrypts cipher text and returns the original string
func Decrypt(cipherText string, key []byte) (string, error) {
	data, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}
	plaintext, err := aesCBCDecrypt(data, key)
	if err != nil {
		return "", err
	}
	return strutils.BytesToString(plaintext), nil
}

func aesCBCEncrypt(plaintext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	// padding
	blockSize := block.BlockSize()
	plaintext = pkcs7Padding(plaintext, blockSize)
	// cipher
	cipherText := make([]byte, blockSize+len(plaintext))
	// initial vector, need to be random
	// different initial vector with the same key will get different cipher text,
	// can treat it as one-time session
	// iv length should be the same with blockSize
	iv := cipherText[:blockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	// encrypt
	blockMode := cipher.NewCBCEncrypter(block, iv)
	blockMode.CryptBlocks(cipherText[blockSize:], plaintext)
	return cipherText, nil
}

func aesCBCDecrypt(cipherText, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	if len(cipherText) < blockSize {
		return nil, fmt.Errorf("cipher text length: %d is shorter than block size: %d", len(cipherText), blockSize)
	}
	// initial vector
	iv := cipherText[:blockSize]
	cipherText = cipherText[blockSize:]
	// CBC mode always works in whole blocks
	if len(cipherText)%blockSize != 0 {
		return nil, fmt.Errorf("cipher text length: %d is not a multiple of the block size: %d", len(cipherText), blockSize)
	}
	// decrypt
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(cipherText, cipherText)
	// unPadding
	return pkcs7UnPadding(cipherText)
}

func pkcs7Padding(plaintext []byte, blockSize int) []byte {
	// padding number
	padding := blockSize - len(plaintext)%blockSize
	// if padding <= 256 (2 ** 8), it can be put in 1 byte
	// blockSize is 16 (default) or 32, < 256, so 1 byte is enough
	// pkcs7 requires at least one byte to represent padding number
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(plaintext, padText...)
}

func pkcs7UnPadding(origData []byte) ([]byte, error) {
	length := len(origData)
	unPadding := int(origData[length-1])
	// wrong key
	if length-unPadding < 0 {
		return nil, fmt.Errorf("wrong unPadding length: length - unPadding = %d, should be wrong key", length-unPadding)
	}
	return origData[:(length - unPadding)], nil
}
