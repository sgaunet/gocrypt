package cmd

import (
	"crypto/rand"
	"fmt"
	"math/big"

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
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(genRandomString(aes.KeyLenAES128))
	},
}

var genAES256Cmd = &cobra.Command{
	Use:   "aes256",
	Short: "generate AES 256 key",
	Long:  `generate AES 256 key`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(genRandomString(aes.KeyLenAES256))
	},
}

var genAES512Cmd = &cobra.Command{
	Use:   "aes512",
	Short: "generate AES 512 key",
	Long:  `generate AES 512 key`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(genRandomString(aes.KeyLenAES512))
	},
}

func init() {
	genCmd.AddCommand(genAES128Cmd)
	genCmd.AddCommand(genAES256Cmd)
	genCmd.AddCommand(genAES512Cmd)
}

// genRandomString generates a random string of length n
func genRandomString(n int) (string, error) {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	max := big.NewInt(int64(len(letters)))
	b := make([]rune, n)
	for i := range b {
		idx, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		b[i] = letters[int(idx.Int64())]
	}
	return string(b), nil
}
