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
	Short: "Tool to encrypt/decrypt file in AES128/256/512.",
	Long:  `Tool to encrypt/decrypt file in AES128/256/512.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
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
}
