package cmd

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"

	"github.com/sgaunet/gocrypt/internal/aes"
	"github.com/spf13/cobra"
)

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "generate AES key",
	Long:  `generate AES key`,
}

var genAES128Cmd = &cobra.Command{
	Use:   "aes128",
	Short: "generate AES 128 key",
	Long:  `generate AES 128 key`,
	Run: func(_ *cobra.Command, _ []string) {
		res, err := genRandomString(aes.KeyLenAES128)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		fmt.Println(res)
	},
}

var genAES256Cmd = &cobra.Command{
	Use:   "aes256",
	Short: "generate AES 256 key",
	Long:  `generate AES 256 key`,
	Run: func(_ *cobra.Command, _ []string) {
		res, err := genRandomString(aes.KeyLenAES256)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		fmt.Println(res)
	},
}

func init() {
	genCmd.AddCommand(genAES128Cmd)
	genCmd.AddCommand(genAES256Cmd)
}

// genRandomString generates a random string of length n.
func genRandomString(n int) (string, error) {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	maxIdx := big.NewInt(int64(len(letters)))
	b := make([]rune, n)
	for i := range b {
		idx, err := rand.Int(rand.Reader, maxIdx)
		if err != nil {
			return "", fmt.Errorf("failed to generate random int: %w", err)
		}
		b[i] = letters[int(idx.Int64())]
	}
	return string(b), nil
}
