package aes

import (
	"bytes"
	"crypto/rand"
	"os"
	"testing"
)

func TestEncryptDecryptFile(t *testing.T) {
	key := make([]byte, KeyLenAES128)
	_, err := rand.Read(key)
	if err != nil {
		t.Fatalf("failed to generate random key: %v", err)
	}

	plaintext := []byte("this is a secret message for testing")

	inputFile, err := os.CreateTemp("", "gocrypt_input_*.txt")
	if err != nil {
		t.Fatalf("failed to create temp input file: %v", err)
	}
	defer os.Remove(inputFile.Name())
	if _, err := inputFile.Write(plaintext); err != nil {
		t.Fatalf("failed to write to input file: %v", err)
	}
	inputFile.Close()

	outputFile, err := os.CreateTemp("", "gocrypt_output_*.bin")
	if err != nil {
		t.Fatalf("failed to create temp output file: %v", err)
	}
	outputFile.Close()
	defer os.Remove(outputFile.Name())

	// Encrypt
	err = EncryptFile(key, inputFile.Name(), outputFile.Name())
	if err != nil {
		t.Fatalf("EncryptFile failed: %v", err)
	}

	// Decrypt
	decryptedFile, err := os.CreateTemp("", "gocrypt_decrypted_*.txt")
	if err != nil {
		t.Fatalf("failed to create temp decrypted file: %v", err)
	}
	decryptedFile.Close()
	defer os.Remove(decryptedFile.Name())

	err = DecryptFile(key, outputFile.Name(), decryptedFile.Name())
	if err != nil {
		t.Fatalf("DecryptFile failed: %v", err)
	}

	decrypted, err := os.ReadFile(decryptedFile.Name())
	if err != nil {
		t.Fatalf("failed to read decrypted file: %v", err)
	}

	if !bytes.Equal(plaintext, decrypted) {
		t.Fatalf("decrypted data does not match original plaintext\nGot: %s\nWant: %s", decrypted, plaintext)
	}
}

func TestGetKeyFromFileErrors(t *testing.T) {
	// Non-existent file
	_, err := GetKeyFromFile("/tmp/nonexistent_gocrypt_key.txt")
	if err == nil {
		t.Error("expected error for non-existent file, got nil")
	}

	// Invalid key length (too short)
	tmpFile, err := os.CreateTemp("", "gocrypt_key_short_*.txt")
	if err != nil {
		t.Fatalf("failed to create temp key file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	if _, err := tmpFile.Write([]byte("shortkey")); err != nil {
		t.Fatalf("failed to write short key: %v", err)
	}
	tmpFile.Close()
	_, err = GetKeyFromFile(tmpFile.Name())
	if err == nil {
		t.Error("expected error for invalid key length, got nil")
	}

	// Invalid key length (not 16 or 32)
	tmpFile2, err := os.CreateTemp("", "gocrypt_key_invalidlen_*.txt")
	if err != nil {
		t.Fatalf("failed to create temp key file: %v", err)
	}
	defer os.Remove(tmpFile2.Name())
	if _, err := tmpFile2.Write([]byte("123456789012345678901234567")); err != nil {
		t.Fatalf("failed to write invalid length key: %v", err)
	}
	tmpFile2.Close()
	_, err = GetKeyFromFile(tmpFile2.Name())
	if err == nil {
		t.Error("expected error for invalid key length, got nil")
	}
}
