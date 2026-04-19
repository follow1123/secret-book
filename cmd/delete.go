package cmd

import (
	"github.com/follow1123/secret-book/bookmanager"
	"github.com/spf13/cobra"
)

var (
	deleteFlagForce   bool
	deleteFlagHistory bool
)

var deleteCmd = &cobra.Command{
	Use:          "delete id",
	Aliases:      []string{"rm"},
	Short:        "delete secret",
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true, // 关闭错误时的帮助信息
	GroupID:      cmdGrpDefault,
	RunE: func(cmd *cobra.Command, args []string) error {
		passwd, err := readBookPassword()
		if err != nil {
			return err
		}
		if secretsFile == "" {
			secretsFile = bookmanager.DefaultSecretsFile()
		}
		bm, err := bookmanager.New(secretsFile, passwd)
		if err != nil {
			return err
		}
		id := args[0]
		if deleteFlagHistory {
			if err := bm.DeleteHistoryByIdPrefix(id); err != nil {
				return err
			}
		} else {
			if err := bm.DeleteByIdPrefix(id, !deleteFlagForce); err != nil {
				return err
			}
		}

		if err := bm.Save(); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	deleteCmd.Flags().BoolVarP(&deleteFlagForce, "force", "f", false, "force delete do not save to history")
	deleteCmd.Flags().BoolVarP(&deleteFlagHistory, "history", "H", false, "delete history by id")

	rootCmd.AddCommand(deleteCmd)
}
