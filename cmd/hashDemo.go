package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

/*
	This file demo'es the use of hash functions in Go
 */
func main() {
	secret := []byte("the shared secret key here")
	message := []byte("the message to hash here")

	hash := hmac.New(sha256.New, secret)
	hash.Write(message)

	// Encode to hex
	fmt.Println(hex.EncodeToString(hash.Sum(nil)))

	/*message = []byte("Another message")

	hash = hmac.New(sha256.New, secret)
	hash.Write(message)

	// Encode to hex
	fmt.Println(hex.EncodeToString(hash.Sum(nil)))*/

	// In case you need base64 encoding
	//base64.StdEncoding.EncodeToString(hash.Sum(nil))
}
