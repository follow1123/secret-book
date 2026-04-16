package cmd

import (
	"github.com/follow1123/secret-book/bookmanager"
	"github.com/spf13/cobra"
)

var (
	updateFlagPlatform string
	updateFlagAccount  string
	updateFlagRemark   string
	updateFlagPassword bool
)

var updateCmd = &cobra.Command{
	Use:          "update id",
	Aliases:      []string{"rm"},
	Short:        "update secret",
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true, // 关闭错误时的帮助信息
	GroupID:      cmdGrpDefault,
	RunE: func(cmd *cobra.Command, args []string) error {
		bm, err := bookmanager.New(bookmanager.DefaultSecretsFile())
		if err != nil {
			return err
		}
		id := args[0]
		secret := bookmanager.Secret{
			Platform: updateFlagPlatform,
			Account:  updateFlagAccount,
			Remark:   updateFlagRemark,
		}
		if updateFlagPassword {
			password, err := readPassword()
			if err != nil {
				return err
			}
			secret.Password = password
		}

		if err := bm.UpdateByIdPrefix(id, secret); err != nil {
			return err
		}
		if err := bm.Save(); err != nil {
			return err
		}

		return nil
	},
}

func init() {

	updateCmd.Flags().StringVarP(&updateFlagPlatform, "platform", "p", "", "update platform")
	updateCmd.Flags().StringVarP(&updateFlagAccount, "account", "a", "", "update account")
	updateCmd.Flags().StringVarP(&updateFlagRemark, "remark", "r", "", "update remark")
	updateCmd.Flags().BoolVarP(&updateFlagPassword, "password", "P", false, "update password")

	rootCmd.AddCommand(updateCmd)
}
