package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var inputFile, outputFile string
var keyFile string
var rmOption bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gocrypt",
	Short: "Tool to encrypt/decrypt file in AES128/256.",
	Long: `Tool to encrypt/decrypt file in AES128/256. 
	
The key can be provided in a file or in the GOCRYPT_KEY environment variable.
The length of the key string must be 12 (AES128), 32 (AES256).`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	encCmd.Flags().StringVar(&inputFile, "i", "", "file to encrypt")
	encCmd.Flags().StringVar(&outputFile, "o", "", "output file")
	encCmd.Flags().StringVar(&keyFile, "k", "", "file containing the key to encrypt (or set GOCRYPT_KEY env variable)")
	encCmd.Flags().BoolVar(&rmOption, "del", false, "delete source file after encryption")
	rootCmd.AddCommand(encCmd)

	decCmd.Flags().StringVar(&inputFile, "i", "", "file to decrypt")
	decCmd.Flags().StringVar(&outputFile, "o", "", "output file")
	decCmd.Flags().StringVar(&keyFile, "k", "", "file containing the key to decrypt (or set GOCRYPT_KEY env variable)")
	decCmd.Flags().BoolVar(&rmOption, "del", false, "delete source file after decryption")
	rootCmd.AddCommand(decCmd)
	rootCmd.AddCommand(genCmd)
}
