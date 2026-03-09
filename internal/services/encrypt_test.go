package services

import (
	"fmt"
	"os"
	"testing"

	"golang.org/x/crypto/argon2"
)

var (
	password = os.Getenv("PASSWORD")
	salt     = []byte(os.Getenv("SALT"))
)

func TestEncrypt(t *testing.T) {
	if password == "" {
		t.Skip("PASSWORD is not set")
	}

	key := argon2.IDKey([]byte(password), salt, 1, 164*10244, 4, 32)
	secret, _ := EncryptText("Hello Whit3roof", key)
	fmt.Println("Encrypted text: ", secret)

	original, _ := DecryptText(secret, key)
	fmt.Println("Original text: ", original)
}

func TestDecrypt(t *testing.T) {
	if password == "" {
		t.Skip("PASSWORD is not set")
	}

	key := argon2.IDKey([]byte(password), salt, 1, 164*10244, 4, 32)

	original, _ := DecryptText("BH0vEStZQfoBY7LxIY4WBaTp+ouhQmwLAa75tXz0NegvGEpQqKQ=", key)
	fmt.Println("Decrypted text: ", original)
}
