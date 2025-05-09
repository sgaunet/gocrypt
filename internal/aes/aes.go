package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/binary"
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
// Each chunk is encrypted and authenticated independently with its own random nonce and tag.
func EncryptFile(key []byte, reader io.Reader, writer io.Writer) error {
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	const chunkSize = 32 * 1024 // 32KB
	buf := make([]byte, chunkSize)
	for {
		n, readErr := io.ReadFull(reader, buf)
		if n > 0 {
			plaintext := buf[:n]
			nonce := make([]byte, gcm.NonceSize())
			if _, err := rand.Read(nonce); err != nil {
				return err
			}
			ciphertext := gcm.Seal(nil, nonce, plaintext, nil)
			// Write chunk length (uint32, big endian)
			var chunkLen [4]byte
			binary.BigEndian.PutUint32(chunkLen[:], uint32(n))
			if _, err := writer.Write(chunkLen[:]); err != nil {
				return err
			}
			// Write nonce
			if _, err := writer.Write(nonce); err != nil {
				return err
			}
			// Write ciphertext (includes tag)
			if _, err := writer.Write(ciphertext); err != nil {
				return err
			}
		}
		if readErr == io.EOF {
			break
		}
		if readErr == io.ErrUnexpectedEOF {
			break
		}
		if readErr != nil {
			return readErr
		}
	}
	return nil
}

// DecryptFile decrypts data from the provided io.Reader and writes the decrypted output to the provided io.Writer.
// Each chunk is decrypted and authenticated independently.
func DecryptFile(key []byte, reader io.Reader, writer io.Writer) error {
	block, err := aes.NewCipher(key)
	if err != nil {
		return errors.New("cipher err: " + err.Error())
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	const nonceSize = 12 // GCM standard
	const tagSize = 16   // GCM standard
	for {
		var chunkLenBuf [4]byte
		_, err := io.ReadFull(reader, chunkLenBuf[:])
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		chunkLen := binary.BigEndian.Uint32(chunkLenBuf[:])
		// Read nonce
		nonce := make([]byte, nonceSize)
		if _, err := io.ReadFull(reader, nonce); err != nil {
			return err
		}
		// Read ciphertext (chunkLen + tagSize)
		ciphertext := make([]byte, chunkLen+uint32(tagSize))
		if _, err := io.ReadFull(reader, ciphertext); err != nil {
			return err
		}
		plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
		if err != nil {
			return err
		}
		if _, err := writer.Write(plaintext); err != nil {
			return err
		}
	}
	return nil
}
