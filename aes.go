package main

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

// (128, 192, or 256 bits // 8 bits= one character)
const KEY_LENGTH_AES128 = 16
const KEY_LENGTH_AES256 = 24
const KEY_LENGTH_AES512 = 32

func getKey(keyFilename string) ([]byte, error) {
	if !isFileExists(keyFilename) {
		return nil, errors.New("file " + keyFilename + " does not exist.")
	}

	key, err := os.ReadFile(keyFilename)
	if err != nil {
		return nil, err
	}
	keyWithoutCR := strings.Trim(string(key), "\r\n")

	if len(keyWithoutCR) != KEY_LENGTH_AES128 && len(keyWithoutCR) != KEY_LENGTH_AES256 && len(keyWithoutCR) != KEY_LENGTH_AES512 {
		errMsg := fmt.Sprintf("length of key should be %d, %d or %d characters if you want to respectively encrypt in AES-128, AES-256 or AES-512", KEY_LENGTH_AES128, KEY_LENGTH_AES256, KEY_LENGTH_AES512)
		return nil, errors.New(errMsg)
	}

	return []byte(keyWithoutCR), err
}

func encryptFile(key []byte, inputFile string, outputFile string) error {
	// Creating block of algorithm
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	reader, err := os.Open(inputFile)
	if err != nil {
		return err
	}
	defer reader.Close()

	writer, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer writer.Close()

	iv := make([]byte, aes.BlockSize)
	stream := cipher.NewOFB(block, iv)
	cipherWriter := &cipher.StreamWriter{
		S: stream,
		W: writer,
	}
	if _, err = io.Copy(cipherWriter, reader); err != nil {
		return err
	}
	return nil
}

func decryptFile(key []byte, inputFile string, outputFile string) error {
	// Creating block of algorithm
	block, err := aes.NewCipher(key)
	if err != nil {
		return errors.New("cipher err: " + err.Error())
	}
	reader, err := os.Open(inputFile)
	if err != nil {
		return err
	}
	defer reader.Close()

	f, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer f.Close()

	iv := make([]byte, aes.BlockSize)
	stream := cipher.NewOFB(block, iv)
	cipherReader := &cipher.StreamReader{S: stream, R: reader}
	if _, err = io.Copy(f, cipherReader); err != nil {
		return err
	}

	return nil
}
