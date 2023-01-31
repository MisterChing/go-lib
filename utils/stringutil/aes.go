package stringutil

import (
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
