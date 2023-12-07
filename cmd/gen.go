package cmd

import (
	"fmt"
	"math/rand"

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
		fmt.Println(genRandomString(16))
	},
}

var genAES256Cmd = &cobra.Command{
	Use:   "aes256",
	Short: "generate AES 256 key",
	Long:  `generate AES 256 key`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(genRandomString(24))
	},
}

var genAES512Cmd = &cobra.Command{
	Use:   "aes512",
	Short: "generate AES 512 key",
	Long:  `generate AES 512 key`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(genRandomString(32))
	},
}

func init() {
	genCmd.AddCommand(genAES128Cmd)
	genCmd.AddCommand(genAES256Cmd)
	genCmd.AddCommand(genAES512Cmd)
}

// genRandomString generates a random string of length n
func genRandomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
