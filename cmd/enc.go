package cmd

import (
	"fmt"
	"os"

	"github.com/sgaunet/gocrypt/internal/aes"
	"github.com/spf13/cobra"
)

// setCmd represents the set command
var encCmd = &cobra.Command{
	Use:   "enc",
	Short: "encrypt file in AES 128/256/512",
	Long:  `encrypt file in AES 128/256/512`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error

		if inputFile == "" {
			fmt.Fprintf(os.Stderr, "inputfile not specified\n")
			cmd.Help()
			os.Exit(1)
		}

		if !isFileExists(inputFile) {
			fmt.Fprintf(os.Stderr, "File %s does not exist\n", inputFile)
			os.Exit(1)
		}

		if outputFile == "" {
			outputFile = fmt.Sprintf("%s.enc", inputFile)
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

		err = aes.EncryptFile(key, inputFile, outputFile)
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
