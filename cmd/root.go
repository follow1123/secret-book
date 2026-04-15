package cmd

import (
	"os"

	"github.com/spf13/cobra"
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
