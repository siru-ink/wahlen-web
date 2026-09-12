package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
)

var encryptionKey []byte

func init() {
	keyEnvVar := os.Getenv("SERVER_ENCRYPTION_KEY")
	if keyEnvVar == "" {
		log.Fatal("SERVER_ENCRYPTION_KEY environmental variable not set")
	}

	decodedKey, err := hex.DecodeString(keyEnvVar)
	if err != nil {
		log.Fatal("SERVER_ENCRYPTION_KEY contained malformed hex characters")
	}
	if len(decodedKey) != 32 {
		log.Fatal("SERVER_ENCRYPTION_KEY was not 32 bytes long")
	}

	encryptionKey = decodedKey
}

func Encrypt(text string) (encryptedMessage string, err error) {
	plaintext := []byte(text)

	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return "", err
	}

	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)

	return fmt.Sprintf("%x", ciphertext), nil
}

func Decrypt(text string) (decryptedMessage string, err error) {
	encryptedMessage, err := hex.DecodeString(text)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceLength := aesGCM.NonceSize()

	nonce, ciphertext := encryptedMessage[:nonceLength], encryptedMessage[nonceLength:]

	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s", plaintext), nil
}
