package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
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
	return EncryptVia(text, encryptionKey)
}

func Decrypt(text string) (decryptedMessage string, err error) {
	return DecryptVia(text, encryptionKey)
}

func EncryptVia(text string, key []byte) (ciphertext string, err error) {
	if text == "" {
		return "", errors.New("must provide encryption data")
	}
	if len(key) != 32 {
		return "", errors.New("provided key must be 32 bytes")
	}

	plaintext := []byte(text)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", errors.New("failed to create AES block cipher")
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", errors.New("failed to create GCM block counter")
	}

	nonce := make([]byte, aesGCM.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return "", errors.New("failed to read cryptographic randomness for nonce")
	}

	cipherOutput := aesGCM.Seal(nil, nonce, plaintext, nil)

	cipherValues := append(nonce, cipherOutput...)

	return fmt.Sprintf("%x", cipherValues), nil
}

func DecryptVia(ciphertext string, key []byte) (text string, err error) {
	if ciphertext == "" {
		return "", errors.New("must provide ciphertext")
	}
	if len(key) != 32 {
		return "", errors.New("provided key must be 32 bytes")
	}

	encryptedMessage, err := hex.DecodeString(ciphertext)
	if err != nil {
		return "", errors.New("invalid/undefined hex data in ciphertext")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", errors.New("failed to create AES block cipher")
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", errors.New("failed to create GCM block counter")
	}

	nonceLength := aesGCM.NonceSize()

	nonce, cipherOutput := encryptedMessage[:nonceLength], encryptedMessage[nonceLength:]

	plaintext, err := aesGCM.Open(nil, nonce, cipherOutput, nil)
	if err != nil {
		return "", errors.New("decryption failed")
	}

	return fmt.Sprintf("%s", plaintext), nil
}
