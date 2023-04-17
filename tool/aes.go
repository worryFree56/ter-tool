package tool

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

// =================== CFB ======================
// 加密
func AesEncryptCFB(origData []byte, secret []byte) (encrypted []byte) {
	block, err := aes.NewCipher(secret)
	if err != nil {
		panic(err)
	}
	encrypted = make([]byte, aes.BlockSize+len(origData))
	iv := encrypted[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(encrypted[aes.BlockSize:], origData)
	return encrypted
}

// 解密
func AesDecryptCFB(encrypted []byte, secert []byte) (decrypted []byte) {
	block, _ := aes.NewCipher(secert)
	if len(encrypted) < aes.BlockSize {
		panic("ciphertext too short")
	}
	iv := encrypted[:aes.BlockSize]
	encrypted = encrypted[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(encrypted, encrypted)
	return encrypted
}
