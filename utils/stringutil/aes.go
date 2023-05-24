package stringutil

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"github.com/forgoer/openssl"
)

func AesCBCEncrypt(message, key, iv string) string {
	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()
	data, err := openssl.AesCBCEncrypt([]byte(message), []byte(key), []byte(iv), openssl.PKCS7_PADDING)
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(data)
}

func AesCBCDecrypt(ciphertext, key, iv string) string {
	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()
	src, _ := base64.StdEncoding.DecodeString(ciphertext)
	dst, err := openssl.AesCBCDecrypt(src, []byte(key), []byte(iv), openssl.PKCS7_PADDING)
	if err != nil {
		return ""
	}
	return string(dst)
}

func AesEncrypt(message string, key string) string {
	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()
	// 转成字节数组
	origData := []byte(message)
	k := []byte(key)
	// 分组秘钥
	block, _ := aes.NewCipher(k)
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 补全码
	origData = PKCS7Padding(origData, blockSize)
	// 加密模式
	blockMode := cipher.NewCBCEncrypter(block, k[:blockSize])
	// 创建数组
	encrypted := make([]byte, len(origData))
	// 加密
	blockMode.CryptBlocks(encrypted, origData)
	return base64.StdEncoding.EncodeToString(encrypted)
}

func AesDecrypt(ciphertext string, key string) string {
	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()
	// 转成字节数组
	encryptedByte, _ := base64.StdEncoding.DecodeString(ciphertext)
	k := []byte(key)
	// 分组秘钥
	block, _ := aes.NewCipher(k)
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 加密模式
	blockMode := cipher.NewCBCDecrypter(block, k[:blockSize])
	// 创建数组
	orig := make([]byte, len(encryptedByte))
	// 解密
	blockMode.CryptBlocks(orig, encryptedByte)
	// 去补全码
	orig = PKCS7UnPadding(orig)
	return string(orig)
}

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padText...)
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unPadding := int(origData[length-1])
	return origData[:(length - unPadding)]
}
