package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

// (128, 192, or 256 bits // 8 bits= one character)

// KeyLenAES128 is the length of the key for AES128
const KeyLenAES128 = 16

// KeyLenAES256 is the length of the key for AES256
const KeyLenAES256 = 32

// GetKeyFromFile reads a key from the specified file, trims carriage returns and newlines,
// and checks that the key length is valid for AES-128 (16 bytes) or AES-256 (32 bytes).
// Returns the key as a byte slice or an error if the key is invalid or the file cannot be read.
func GetKeyFromFile(keyFilename string) ([]byte, error) {
	key, err := os.ReadFile(keyFilename)
	if err != nil {
		return nil, err
	}
	keyWithoutCR := strings.Trim(string(key), "\r\n")

	if len(keyWithoutCR) != KeyLenAES128 && len(keyWithoutCR) != KeyLenAES256 {
		errMsg := fmt.Sprintf("length of key should be %d (AES128), %d (AES256)", KeyLenAES128, KeyLenAES256)
		return nil, errors.New(errMsg)
	}

	return []byte(keyWithoutCR), err
}

// EncryptFile encrypts data from the provided io.Reader and writes the encrypted output to the provided io.Writer.
func EncryptFile(key []byte, reader io.Reader, writer io.Writer) error {
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return err
	}
	// Write the nonce at the beginning of the output
	if _, err := writer.Write(nonce); err != nil {
		return err
	}

	plaintext, err := io.ReadAll(reader)
	if err != nil {
		return err
	}
	ciphertext := gcm.Seal(nil, nonce, plaintext, nil)
	if _, err := writer.Write(ciphertext); err != nil {
		return err
	}
	return nil
}

// DecryptFile decrypts data from the provided io.Reader and writes the decrypted output to the provided io.Writer.
func DecryptFile(key []byte, reader io.Reader, writer io.Writer) error {
	block, err := aes.NewCipher(key)
	if err != nil {
		return errors.New("cipher err: " + err.Error())
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(reader, nonce); err != nil {
		return err
	}
	ciphertext, err := io.ReadAll(reader)
	if err != nil {
		return err
	}
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return err
	}
	if _, err := writer.Write(plaintext); err != nil {
		return err
	}
	return nil
}
