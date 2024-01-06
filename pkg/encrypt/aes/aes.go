package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"github.com/lanyulei/toolkit/logger"
	"io"
)

/*
  @Author : lanyulei
  @Desc :
*/

// Encrypt 使用AES加密数据
func Encrypt(plaintext, key []byte) (ciphertext []byte, err error) {
	var (
		block cipher.Block
	)

	block, err = aes.NewCipher(key)
	if err != nil {
		logger.Errorf("aes.NewCipher error: %v", err)
		return
	}

	ciphertext = make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]

	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		logger.Errorf("io.ReadFull error: %v", err)
		return
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)

	return
}

// Decrypt 使用AES解密数据
func Decrypt(ciphertext, key []byte) (result []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		logger.Errorf("aes.NewCipher error: %v", err)
		return
	}

	if len(ciphertext) < aes.BlockSize {
		logger.Errorf("密文太短")
		return
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)

	result = unfill(ciphertext)
	return
}

func unfill(data []byte) []byte {
	length := len(data)
	unpadding := int(data[length-1])
	return data[:(length - unpadding)]
}
