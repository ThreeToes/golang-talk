package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"log"
)

func main() {
	key := "SuperSecretKey"
	payload, _ := base64.StdEncoding.DecodeString(
		"WtHbeerkPQoESMS50WJ9pW2f+nMHaebFcx0ThGVWhPBhPgiK/e0V35T0X/gQec5+zQQ3vNlp17Si9JvvLJw=")
	key256 := sha256.Sum256([]byte(key))
	cphr, _ := aes.NewCipher(key256[:])
	gcm, _ := cipher.NewGCM(cphr)
	if len(payload) < gcm.NonceSize() {
		log.Fatalf("Nonce too short for payload")
	}
	nonce, encryptedText := payload[:gcm.NonceSize()], payload[gcm.NonceSize():]
	plainText, _ := gcm.Open(nil, nonce, encryptedText, nil)
	log.Printf("Got plaintext '%s'", plainText)
}
