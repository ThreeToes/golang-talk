package encrypter

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
)

func EncryptData(key string, payload []byte) ([]byte, error) {
	// We need to make sure our key is 32 bytes long so we pick AES-256
	key256 := sha256.Sum256([]byte(key))
	cphr, err := aes.NewCipher(key256[:])
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(cphr)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	b := gcm.Seal(nonce, nonce, payload, nil)
	return b, nil
}

func DecryptData(key string, payload []byte) ([]byte, error) {
	key256 := sha256.Sum256([]byte(key))
	cphr, err := aes.NewCipher(key256[:])
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(cphr)
	if err != nil {
		return nil, err
	}
	if len(payload) < gcm.NonceSize() {
		return nil, fmt.Errorf("payload too short for nonce")
	}
	nonce, encryptedText := payload[:gcm.NonceSize()], payload[gcm.NonceSize():]
	return gcm.Open(nil, nonce, encryptedText, nil)
}