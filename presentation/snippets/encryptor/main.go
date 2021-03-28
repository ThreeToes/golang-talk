package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"log"
)

func main()  {
	key := "SuperSecretKey"
	payload := "This is something really important"
	// We need to make sure our key is 32 bytes long so we pick AES-256
	key256 := sha256.Sum256([]byte(key))
	cphr, _ := aes.NewCipher(key256[:])
	gcm, _ := cipher.NewGCM(cphr)
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return
	}
	b := gcm.Seal(nonce, nonce, []byte(payload), nil)
	b64str := base64.StdEncoding.EncodeToString(b)
	log.Printf("Encrypted payload: %s", b64str)
}

