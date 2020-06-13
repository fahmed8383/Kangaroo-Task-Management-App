package api

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

// Encrypt uses the encryption key to return an encrypted string
func Encrypt(key string, plaintext string) (string, error) {

	// create a block cipher using the encryption key
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "nil", err
	}

	// returns the cipher block wrapped in GCM with standard nonce length
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "nil", err
	}

	// create an empty array to store the nonce
	nonce := make([]byte, gcm.NonceSize())
	if err != nil {
		return "nil", err
	}

	// populated the nonce with a random number
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return "nil", err
	}

	// append the nonce before the encrypted data so it can be used later to decryption
	encrypted := gcm.Seal(nonce, nonce, []byte(plaintext), nil)

	return string(encrypted), nil
}

// Decrypt uses the encryption key to retun a decrypted string
func Decrypt(key string, encrypted string) (string, error) {

	// create a block cipher using the encryption key
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "nil", err
	}

	// returns the cipher block wrapped in GCM with standard nonce length
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "nil", err
	}

	// get the nonce appended to the encrypted string
	nonce := []byte(encrypted)[:gcm.NonceSize()]

	// get the text that needs to be decrypted from the encrypted string
	encryptedText := []byte(encrypted)[gcm.NonceSize():]

	// decrypt the gcm wrapped text
	plaintext, err := gcm.Open(nil, nonce, encryptedText, nil)
	if err != nil {
		return "nil", err
	}

	return string(plaintext), nil
}
