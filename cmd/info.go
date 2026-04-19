package cmd

import (
	"os"

	"github.com/follow1123/secret-book/bookmanager"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/tw"
	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:          "info id",
	Short:        "secret details",
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
		table := tablewriter.NewTable(os.Stdout, tablewriter.WithEastAsian(tw.Off))
		tableData := [][]string{}

		id := args[0]
		secret, err := bm.GetSecretByIdPerfix(id)
		if err != nil {
			return err
		}

		tableData = append(tableData, []string{"ID", secret.Id})
		tableData = append(tableData, []string{"PLATFORM", secret.Platform})
		tableData = append(tableData, []string{"ACCOUNT", secret.Account})
		tableData = append(tableData, []string{"PASSWORD", secret.Password})
		tableData = append(tableData, []string{"REMARK", addNewlinesEveryNChars(secret.Remark, 30)})
		tableData = append(tableData, []string{"CREATE TIME", secret.CreateTime})

		if err := table.Bulk(tableData); err != nil {
			return err
		}

		if err := table.Render(); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
