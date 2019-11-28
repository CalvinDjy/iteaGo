package algorithm

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

var secretKey = "MYgGnQE2jD4FADSFFDSEWsd451Dadfgh"

func Aes128Encrypt(encodeStr string, key string) (string, error) {
	k := getKey(key)
	encodeBytes := []byte(encodeStr)
	//根据key 生成密文
	block, err := aes.NewCipher(k)
	if err != nil {
		return "", err
	}

	blockSize := block.BlockSize()
	encodeBytes = padding(encodeBytes, blockSize)

	blockMode := cipher.NewCBCEncrypter(block, k)
	crypted := make([]byte, len(encodeBytes))
	blockMode.CryptBlocks(crypted, encodeBytes)

	return base64.StdEncoding.EncodeToString(crypted), nil
}

func Aes128Decrypt(decodeStr string, key string) (string, error) {
	k := getKey(key)
	//先解密base64
	decodeBytes, err := base64.StdEncoding.DecodeString(decodeStr)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(k)
	if err != nil {
		return "", err
	}
	blockMode := cipher.NewCBCDecrypter(block, k)
	origData := make([]byte, len(decodeBytes))

	blockMode.CryptBlocks(origData, decodeBytes)
	origData = unpadding(origData)
	return string(origData), nil
}

func padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	//填充
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)

	return append(ciphertext, padtext...)
}

func unpadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func getKey(key string) []byte {
	k := []byte(key)
	l := len(k)
	switch {
	case l == 0:
		return []byte(secretKey)
	case l >= 16:
		return k[:16]
	default:
		return append(k, []byte(secretKey)[:16-l]...)
	}
}
