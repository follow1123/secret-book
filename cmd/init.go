package cmd

import (
	"fmt"
	"os"

	"github.com/follow1123/secret-book/bookmanager"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:          "init",
	Short:        "init secret file",
	SilenceUsage: true, // 关闭错误时的帮助信息
	GroupID:      cmdGrpDefault,
	RunE: func(cmd *cobra.Command, args []string) error {
		secretPath := bookmanager.DefaultSecretsFile()

		_, err := os.Stat(secretPath)
		notExists := false
		if err != nil {
			if os.IsNotExist(err) {
				notExists = true
			} else {
				return fmt.Errorf("check book path %s error:\n\t%w", secretPath, err)
			}
		}

		if notExists {
			password, err := readPassword()
			if err != nil {
				return err
			}

			key, err := bookmanager.GenerateKey(password)
			if err != nil {
				return fmt.Errorf("generate key error:\n\t%w", err)
			}
			encryptedSecret, err := bookmanager.Encrypt([]byte("{}"), key)
			if err != nil {
				return fmt.Errorf("encrypt secret error:\n\t%w", err)
			}
			if err := os.WriteFile(secretPath, encryptedSecret, 0664); err != nil {
				return fmt.Errorf("save secret error:\n\t%w", err)
			}

		} else {
			cmd.Printf("secret file is already init at %s\n", secretPath)

		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
