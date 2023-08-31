package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/sgaunet/gocrypt/pkg/aes"
)

var version string = "development"

func printVersion() {
	fmt.Println(version)
}

func main() {
	var inputFile, outputFile string
	var keyFile string
	var helpOption, encryptOption, decryptOption, vOption, rmOption bool
	var err error

	flag.StringVar(&inputFile, "i", "", "File to encrypt/decrypt")
	flag.StringVar(&outputFile, "o", "", "Output file")
	flag.StringVar(&keyFile, "k", "", "Name of file containing the key")
	flag.BoolVar(&encryptOption, "e", false, "Encrypt option")
	flag.BoolVar(&decryptOption, "d", false, "Decrypt option")
	flag.BoolVar(&vOption, "v", false, "Get version")
	flag.BoolVar(&helpOption, "h", false, "Print help")
	flag.BoolVar(&rmOption, "rm", false, "remove input file after encryption/decryption")
	flag.Parse()

	if vOption {
		printVersion()
		os.Exit(0)
	}

	if inputFile == "" {
		fmt.Fprintf(os.Stderr, "error: option -i is mandatory\n\n")
		flag.PrintDefaults()
		os.Exit(1)
	}
	if keyFile == "" {
		fmt.Fprintf(os.Stderr, "error: option -k is mandatory\n\n")
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

	if outputFile == "" && encryptOption {
		outputFile = fmt.Sprintf("%s.enc", inputFile)
		fmt.Println("outputfile not specified, initialised to :", outputFile)
	}
	if outputFile == "" && decryptOption {
		outputFile = strings.Replace(inputFile, ".enc", "", 1)
		if inputFile == outputFile {
			outputFile = fmt.Sprintf("%s.ori", inputFile)
		}
		fmt.Println("outputfile not specified, initialised to :", outputFile)
	}

	if isFileExists(outputFile) {
		err = os.Remove(outputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot remove file %s\n", outputFile)
			os.Exit(1)
		}
	}

	key, err := aes.GetKey(keyFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err.Error())
		os.Exit(1)
	}

	if encryptOption {
		err = aes.EncryptFile(key, inputFile, outputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err.Error())
			os.Exit(1)
		}
	}
	if decryptOption {
		err = aes.DecryptFile(key, inputFile, outputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err.Error())
			os.Exit(1)
		}
	}

	if rmOption {
		err = os.Remove(inputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err.Error())
			os.Exit(1)
		}
	}
}

func isFileExists(file string) bool {
	f, err := os.Open(file)
	if os.IsNotExist(err) {
		return false
	}
	defer f.Close()
	i, _ := os.Stat(file)
	return !i.IsDir()
}
