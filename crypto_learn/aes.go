package main

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"log"
)

// 对称密钥
const key = "hellcxzcadwasdad"

var iv = []byte("123451helloworld")

func main() {
	// 待加密的明文
	plaintext := []byte("abcdbcjdhwjdhchd")
	// 加密
	ciphertext, err := encrypt(plaintext)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Ciphertext: %x\n", ciphertext)

	// 解密
	decrypted, err := decrypt(ciphertext)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Decrypted: %s\n", decrypted)
}

func encrypt(plaintext []byte) ([]byte, error) {
	aesCipher, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}
	// 密文,分组大小,加上一个初始向量的预留空间(为分组长度)
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	// 分组模式,使用CBC(使用向量)
	mode := cipher.NewCBCEncrypter(aesCipher, iv)

	// 注意,此处没有考虑明文的大小,不足分组长度是需要填充算法的

	// 对明文进行加密(去除向量头部)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)
	// 初始向量放在头部
	copy(ciphertext[:aes.BlockSize], iv)
	return ciphertext, nil
}

func decrypt(ciphertext []byte) ([]byte, error) {
	aesCipher, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}
	// 获取分组大小
	blockSize := aes.BlockSize
	// 分组模式,使用CBC(使用向量)
	mode := cipher.NewCBCDecrypter(aesCipher, iv)
	// 密文,分组大小
	plaintext := make([]byte, len(ciphertext)-blockSize)
	// 对密文进行解密
	mode.CryptBlocks(plaintext, ciphertext[blockSize:])
	return plaintext, nil
}
