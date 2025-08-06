// Package aes provides AES-GCM encryption and decryption functionality.
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
	"path/filepath"
	"strings"
)

// ErrInvalidKeyLength is returned when the key length is invalid.
var ErrInvalidKeyLength = errors.New("invalid key length")

// ErrCipherCreation is returned when cipher creation fails.
var ErrCipherCreation = errors.New("cipher creation failed")

// ErrInvalidChunkLength is returned when chunk length is invalid.
var ErrInvalidChunkLength = errors.New("invalid chunk length")

// (128, 192, or 256 bits // 8 bits= one character)

// KeyLenAES128 is the length of the key for AES128.
const KeyLenAES128 = 16

// KeyLenAES256 is the length of the key for AES256.
const KeyLenAES256 = 32

// GetKeyFromFile reads a key from the specified file, trims carriage returns and newlines,
// and checks that the key length is valid for AES-128 (16 bytes) or AES-256 (32 bytes).
// Returns the key as a byte slice or an error if the key is invalid or the file cannot be read.
func GetKeyFromFile(keyFilename string) ([]byte, error) {
	// Clean the path to prevent directory traversal
	cleanFilename := filepath.Clean(keyFilename)
	key, err := os.ReadFile(cleanFilename) // #nosec G304 - file path is cleaned
	if err != nil {
		return nil, fmt.Errorf("failed to read key file: %w", err)
	}
	keyWithoutCR := strings.Trim(string(key), "\r\n")

	if len(keyWithoutCR) != KeyLenAES128 && len(keyWithoutCR) != KeyLenAES256 {
		errMsg := fmt.Sprintf("length of key should be %d (AES128), %d (AES256)", KeyLenAES128, KeyLenAES256)
		return nil, fmt.Errorf("%w: %s", ErrInvalidKeyLength, errMsg)
	}

	return []byte(keyWithoutCR), nil
}

// EncryptFile encrypts data from the provided io.Reader and writes the encrypted output to the provided io.Writer.
// Each chunk is encrypted and authenticated independently with its own random nonce and tag.
func EncryptFile(key []byte, reader io.Reader, writer io.Writer) error {
	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrCipherCreation, err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("failed to create GCM: %w", err)
	}

	const chunkSize = 32 * 1024 // 32KB
	buf := make([]byte, chunkSize)
	for {
		n, readErr := io.ReadFull(reader, buf)
		if n > 0 {
			if err := encryptChunk(gcm, writer, buf[:n], n); err != nil {
				return err
			}
		}
		if errors.Is(readErr, io.EOF) || errors.Is(readErr, io.ErrUnexpectedEOF) {
			break
		}
		if readErr != nil {
			return fmt.Errorf("read error: %w", readErr)
		}
	}
	return nil
}

// encryptChunk encrypts a single chunk and writes it to the writer.
func encryptChunk(gcm cipher.AEAD, writer io.Writer, plaintext []byte, chunkLen int) error {
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return fmt.Errorf("failed to generate nonce: %w", err)
	}
	ciphertext := gcm.Seal(nil, nonce, plaintext, nil)
	
	// Write chunk length (uint32, big endian)
	var chunkLenBytes [4]byte
	if chunkLen > 0 && chunkLen <= 0x7FFFFFFF { // Check for safe conversion
		binary.BigEndian.PutUint32(chunkLenBytes[:], uint32(chunkLen))
	} else {
		return fmt.Errorf("%w: %d", ErrInvalidChunkLength, chunkLen)
	}
	if _, err := writer.Write(chunkLenBytes[:]); err != nil {
		return fmt.Errorf("failed to write chunk length: %w", err)
	}
	// Write nonce
	if _, err := writer.Write(nonce); err != nil {
		return fmt.Errorf("failed to write nonce: %w", err)
	}
	// Write ciphertext (includes tag)
	if _, err := writer.Write(ciphertext); err != nil {
		return fmt.Errorf("failed to write ciphertext: %w", err)
	}
	return nil
}

// DecryptFile decrypts data from the provided io.Reader and writes the decrypted output to the provided io.Writer.
// Each chunk is decrypted and authenticated independently.
func DecryptFile(key []byte, reader io.Reader, writer io.Writer) error {
	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrCipherCreation, err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("failed to create GCM: %w", err)
	}

	const nonceSize = 12 // GCM standard
	const tagSize = 16   // GCM standard
	for {
		var chunkLenBuf [4]byte
		_, err := io.ReadFull(reader, chunkLenBuf[:])
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read chunk length: %w", err)
		}
		chunkLen := binary.BigEndian.Uint32(chunkLenBuf[:])
		// Read nonce
		nonce := make([]byte, nonceSize)
		if _, err := io.ReadFull(reader, nonce); err != nil {
			return fmt.Errorf("failed to read nonce: %w", err)
		}
		// Read ciphertext (chunkLen + tagSize)
		ciphertext := make([]byte, chunkLen+uint32(tagSize))
		if _, err := io.ReadFull(reader, ciphertext); err != nil {
			return fmt.Errorf("failed to read ciphertext: %w", err)
		}
		plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
		if err != nil {
			return fmt.Errorf("failed to decrypt chunk: %w", err)
		}
		if _, err := writer.Write(plaintext); err != nil {
			return fmt.Errorf("failed to write plaintext: %w", err)
		}
	}
	return nil
}
