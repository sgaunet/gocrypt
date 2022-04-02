package main

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

// (128, 192, or 256 bits // 8 bits= one character)
const KEY_LENGTH_AES128 = 16
const KEY_LENGTH_AES256 = 24
const KEY_LENGTH_AES512 = 32

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

	if (encryptOption && decryptOption) || (!encryptOption && !decryptOption) {
		fmt.Fprintf(os.Stderr, "choose only one option between encrypt and decrypt")
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
	if len(key) != KEY_LENGTH_AES128 && len(key) != KEY_LENGTH_AES256 && len(key) != KEY_LENGTH_AES512 {
		return nil, errors.New("length of key should be 16, 24 or 32 characters if you want to respectively encrypt in AES-128, AES-256 or AES-512")
	}

	return key, err
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
	stream := cipher.NewOFB(block, iv[:])
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
	stream := cipher.NewOFB(block, iv[:])
	cipherReader := &cipher.StreamReader{S: stream, R: reader}
	if _, err = io.Copy(f, cipherReader); err != nil {
		return err
	}

	return nil
}
