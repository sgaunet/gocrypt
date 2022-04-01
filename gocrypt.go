package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	var inputFile, outputFile string
	var keyFile string
	var helpOption, encryptOption, decryptOption bool
	var err error

	flag.StringVar(&inputFile, "i", "", "File to encrypt/decrypt")
	flag.StringVar(&outputFile, "o", "", "Output file")
	flag.StringVar(&keyFile, "k", "", "Name of file containing the key")
	flag.BoolVar(&encryptOption, "e", false, "encrypt option")
	flag.BoolVar(&decryptOption, "d", false, "decrypt option")
	flag.BoolVar(&helpOption, "h", false, "Print help")
	flag.Parse()

	if inputFile == "" || outputFile == "" || keyFile == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	if encryptOption && decryptOption {
		fmt.Fprintf(os.Stderr, "cannot get two options encrypt and decrypt at the same time")
		os.Exit(1)
	}

	if helpOption {
		flag.PrintDefaults()
		os.Exit(0)
	}

	if !isFileExists(inputFile) {
		fmt.Fprintf(os.Stderr, "File %s does not exist\n", inputFile)
		os.Exit(1)
	}

	if isFileExists(outputFile) {
		err = os.Remove(outputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot remove file %s\n", outputFile)
			os.Exit(1)
		}

	}

	key, err := getKey(keyFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err.Error())
		os.Exit(1)
	}

	if encryptOption {
		err = encryptFile(key, inputFile, outputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err.Error())
			os.Exit(1)
		}
	}
	if decryptOption {
		err = decryptFile(key, inputFile, outputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err.Error())
			os.Exit(1)
		}
	}
}

func getKey(keyFilename string) ([]byte, error) {
	if !isFileExists(keyFilename) {
		return nil, errors.New("file " + keyFilename + " does not exist.")
	}

	key, err := ioutil.ReadFile(keyFilename)
	if err != nil {
		return nil, err
	}
	if len(key) != 16 && len(key) != 32 && len(key) != 64 {
		return nil, errors.New("length of key should be 16, 32 or 64 if you want to ctypt in AES-128, AES-256 or AES-512")
	}

	return key, err
}

func encryptFile(key []byte, inputFile string, outputFile string) error {
	// Reading plaintext file
	plainText, err := ioutil.ReadFile(inputFile)
	if err != nil {
		return err
	}
	// Creating block of algorithm
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	// Creating GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return errors.New("cipher GCM err: " + err.Error())
	}
	// Generating random nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		log.Fatalf("nonce  err: %v", err.Error())
	}

	cipherText := gcm.Seal(nonce, nonce, plainText, nil)

	// Writing ciphertext file
	err = ioutil.WriteFile(outputFile, cipherText, 0600)
	if err != nil {
		return errors.New("write file err: " + err.Error())
	}
	return nil
}

func decryptFile(key []byte, inputFile string, outputFile string) error {
	// Reading ciphertext file
	cipherText, err := ioutil.ReadFile(inputFile)
	if err != nil {
		return err
	}

	// Creating block of algorithm
	block, err := aes.NewCipher(key)
	if err != nil {
		return errors.New("cipher err: " + err.Error())
	}

	// Creating GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return errors.New("cipher GCM err: " + err.Error())
	}

	// Deattached nonce and decrypt
	nonce := cipherText[:gcm.NonceSize()]
	cipherText = cipherText[gcm.NonceSize():]
	plainText, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return errors.New("decrypt file err: " + err.Error())
	}

	// Writing decryption content
	err = ioutil.WriteFile(outputFile, plainText, 0600)
	if err != nil {
		return errors.New("write file err: " + err.Error())
	}
	return nil
}
