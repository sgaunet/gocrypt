package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/sgaunet/gocrypt/internal/aes"
	"github.com/spf13/cobra"
)

// setCmd represents the set command
var decCmd = &cobra.Command{
	Use:   "dec",
	Short: "decrypt file in AES 128/256/512",
	Long:  `decrypt file in AES 128/256/512`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		if !isFileExists(inputFile) {
			fmt.Fprintf(os.Stderr, "File %s does not exist\n", inputFile)
			os.Exit(1)
		}

		if outputFile == "" {
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

		err = aes.DecryptFile(key, inputFile, outputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err.Error())
			os.Exit(1)
		}

		if rmOption {
			err = os.Remove(inputFile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err.Error())
				os.Exit(1)
			}
		}
	},
}
