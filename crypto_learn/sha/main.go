package main

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
)

func main() {
	a := "jhdiasji哈哈哈哈哈哈哈"
	sum := sha256.Sum256([]byte(a))
	// 十六进制显示
	fmt.Printf("%x\n", sum)
	fmt.Println(generateRandom(17))
}

func generateRandom(n int) string {
	if n <= 16 {
		n = 16
	}
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x", b)
}
