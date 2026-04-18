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
	rootCmd.Flags().BoolVarP(&rootFlagVersion, "version", "v", false, "print version")

	rootCmd.AddGroup(&cobra.Group{ID: cmdGrpDefault, Title: "Available Commands"})
}

func readPassword() (string, error) {
	fmt.Print("Enter Password: ")
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
