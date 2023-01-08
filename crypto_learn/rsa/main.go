package main

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"log"
)

func main() {
	// 生成一对密钥(公私)
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal(err)
	}
	publicKey := privateKey.PublicKey

	// 待加密的明文
	plaintext := []byte("你好东澳岛草草草草草草草草草")

	// 使用公钥加密
	ciphertext, err := encrypt(publicKey, plaintext)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Ciphertext: %x\n", ciphertext)

	// 使用私钥解密
	decrypted, err := decrypt(*privateKey, ciphertext)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Decrypted: %s\n", decrypted)
}

func encrypt(publicKey rsa.PublicKey, plaintext []byte) ([]byte, error) {
	// 加密
	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, &publicKey, plaintext)
	if err != nil {
		return nil, err
	}
	return ciphertext, nil
}

func decrypt(privateKey rsa.PrivateKey, ciphertext []byte) ([]byte, error) {
	// 解密
	plaintext, err := rsa.DecryptPKCS1v15(rand.Reader, &privateKey, ciphertext)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}
