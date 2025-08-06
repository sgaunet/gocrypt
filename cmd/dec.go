// Package cmd contains the command-line interface for gocrypt.
package cmd

import (
	"fmt"
	"os"

	"github.com/sgaunet/gocrypt/internal/aes"
	"github.com/spf13/cobra"
)

// setCmd represents the set command.
var decCmd = &cobra.Command{
	Use:   "dec",
	Short: "decrypt file in AES 128/256",
	Long:  `decrypt file in AES 128/256`,
	Run: func(cmd *cobra.Command, _ []string) {
		var (
			tmpFile               *os.File
			overwriteOriginalFile bool
			key                   []byte
			err                   error
		)

		if inputFile == "" {
			fmt.Fprintf(os.Stderr, "inputfile not specified\n")
			_ = cmd.Help()
			os.Exit(1)
		}

		if !isFileExists(inputFile) {
			fmt.Fprintf(os.Stderr, "File %s does not exist\n", inputFile)
			os.Exit(1)
		}

		if outputFile == "" {
			overwriteOriginalFile = true
			tmpFile, err = os.CreateTemp("/tmp", "gocrypt")
			if err != nil {
				fmt.Fprintf(os.Stderr, "Cannot create temp file\n")
				os.Exit(1)
			}
			outputFile = tmpFile.Name()
			_ = tmpFile.Close()
		}

		if isFileExists(outputFile) && !overwriteOriginalFile {
			err = os.Remove(outputFile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Cannot remove file %s\n", outputFile)
				os.Exit(1)
			}
		}

		keyFromEnv := os.Getenv("GOCRYPT_KEY")
		key = []byte(keyFromEnv)
		if len(key) == 0 {
			key, err = aes.GetKeyFromFile(keyFile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err.Error())
				os.Exit(1)
			}
		}

		// Open input and output files for streaming
		inF, err := os.Open(inputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot open input file: %v\n", err)
			os.Exit(1)
		}
		defer func() { _ = inF.Close() }()
		outF, err := os.Create(outputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot create output file: %v\n", err)
			os.Exit(1)
		}
		defer func() { _ = outF.Close() }()

		err = aes.DecryptFile(key, inF, outF)
		if err != nil {
			fmt.Fprintf(os.Stderr, "DecryptFile failed: %v\n", err)
			os.Exit(1)
		}

		if overwriteOriginalFile {
			// rename tmp file to original file
			err = os.Rename(outputFile, inputFile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err.Error())
				os.Exit(1)
			}
		}

		if rmOption && !overwriteOriginalFile {
			err = os.Remove(inputFile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err.Error())
				os.Exit(1)
			}
		}
	},
}
