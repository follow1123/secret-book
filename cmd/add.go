package cmd

import (
	"fmt"

	"github.com/follow1123/secret-book/bookmanager"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add platform account [remark]",
	Short: "add secret",
	// Args:         cobra.ExactArgs(2),
	SilenceUsage: true, // 关闭错误时的帮助信息
	GroupID:      cmdGrpDefault,
	RunE: func(cmd *cobra.Command, args []string) error {
		passwd, err := readPassword("Enter Book Password: ")
		if err != nil {
			return fmt.Errorf("read password error:\n\t%w", err)
		}
		if secretsFile == "" {
			secretsFile = bookmanager.DefaultSecretsFile()
		}
		bm, err := bookmanager.New(secretsFile, passwd)
		if err != nil {
			return err
		}
		secret := bookmanager.Secret{
			Platform: args[0],
			Account:  args[1],
		}

		password, err := readPassword("Enter Password:")
		if err != nil {
			return err
		}
		secret.Password = password

		if len(args) > 2 {
			secret.Remark = args[2]
		}

		if err := bm.Add(secret); err != nil {
			return err
		}
		if err := bm.Save(); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
