package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello world!\n")
	})
	log.Fatal(http.ListenAndServeTLS(
		"localhost:8081", "server.crt", "server.key",
		nil))
	// curl -v -k https://localhost:8081
}
