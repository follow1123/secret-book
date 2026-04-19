package cmd

import (
	"os"

	"github.com/follow1123/secret-book/bookmanager"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/tw"
	"github.com/spf13/cobra"
)

var (
	historyFlagPlatform string
	historyFlagAccount  string
	historyFlagRemark   string
	historyFlagPassword string
)

var historyCmd = &cobra.Command{
	Use:          "history",
	Short:        "secrets history",
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

		tableData = append(tableData, []string{"Id", "Platform", "Account", "Password", "Remark", "Create Time", "Operation Time", "Operation Type"})
		historySecrets := bm.ListHistory(bookmanager.Secret{
			Platform: historyFlagPlatform,
			Account:  historyFlagAccount,
			Password: historyFlagPassword,
			Remark:   historyFlagRemark,
		})
		for _, hs := range historySecrets {
			tableData = append(tableData, []string{hs.Id, hs.Platform, hs.Account, hs.Password, addNewlinesEveryNChars(hs.Remark, 10), hs.CreateTime, hs.OperationTime, string(hs.OperationType)})
		}

		table.Header(tableData[0])
		if err := table.Bulk(tableData[1:]); err != nil {
			return err
		}

		if err := table.Render(); err != nil {
			return err
		}

		return nil
	},
}

func init() {

	historyCmd.Flags().StringVarP(&historyFlagPlatform, "platform", "p", "", "history platform")
	historyCmd.Flags().StringVarP(&historyFlagAccount, "account", "a", "", "history account")
	historyCmd.Flags().StringVarP(&historyFlagRemark, "remark", "r", "", "history remark")
	historyCmd.Flags().StringVarP(&historyFlagPassword, "password", "P", "", "history password")

	rootCmd.AddCommand(historyCmd)
}
