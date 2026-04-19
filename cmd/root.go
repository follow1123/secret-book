package cmd

import (
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

const (
	Version = "0.1.0"
)

const (
	cmdGrpDefault = "available"
)

var rootFlagVersion bool

var secretsFile string

var rootCmd = &cobra.Command{
	Use:   "sbook",
	Short: "secret book",
	Run: func(cmd *cobra.Command, args []string) {
		if rootFlagVersion {
			cmd.Println(Version)
			return
		}
		cmd.Help()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

	rootCmd.PersistentFlags().StringVar(&secretsFile, "secrets-path", "", "target secrets file path")
	rootCmd.Flags().BoolVarP(&rootFlagVersion, "version", "v", false, "print version")

	rootCmd.AddGroup(&cobra.Group{ID: cmdGrpDefault, Title: "Available Commands"})
}

func readBookPassword() (string, error) {
	passwd := readPasswordFromEnv()
	if passwd == "" {
		p, err := readPassword("Enter Book Password: ")
		if err != nil {
			return "", fmt.Errorf("read password error:\n\t%w", err)
		}
		passwd = p
	}
	return passwd, nil
}

func readPasswordFromEnv() string {
	return strings.TrimSpace(os.Getenv("SECRET_BOOK_PASSWORD"))
}

func readPassword(prompt string) (string, error) {
	fmt.Print(prompt)
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", fmt.Errorf("read password error:\n\t%w", err)
	}
	// 清除当前行的提示信息
	fmt.Print("\r")
	fmt.Print(strings.Repeat(" ", 50)) // 覆盖多余的字符
	fmt.Print("\r")                    // 再次回到行首
	return strings.TrimSpace(string(bytePassword)), nil
}

func addNewlinesEveryNChars(s string, charLimit int) string {
	// 将字符串转换为 []rune，确保每个字符完整处理
	runes := []rune(s)
	var result []rune

	// 遍历并每charLimit个字符后添加换行符
	for i := range runes {
		// 每charLimit个字符插入一个换行符
		if i > 0 && i%charLimit == 0 {
			result = append(result, '\n') // 添加换行符
		}
		result = append(result, runes[i]) // 添加当前字符
	}

	// 返回转换后的字符串
	return string(result)
}
