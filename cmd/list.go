package cmd

import (
	"os"

	"github.com/follow1123/secret-book/bookmanager"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/tw"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:          "list [platform]",
	Short:        "list secrets",
	SilenceUsage: true, // 关闭错误时的帮助信息
	GroupID:      cmdGrpDefault,
	RunE: func(cmd *cobra.Command, args []string) error {
		bm, err := bookmanager.New(bookmanager.DefaultSecretsFile())
		if err != nil {
			return err
		}
		table := tablewriter.NewTable(os.Stdout, tablewriter.WithEastAsian(tw.Off))
		tableData := [][]string{}

		if len(args) == 0 {
			tableData = append(tableData, []string{"Platforms"})
			for _, platform := range bm.ListPlatforms() {
				tableData = append(tableData, []string{platform})
			}

		} else {
			tableData = append(tableData, []string{"Id", "Platform", "Account", "Remark", "Create Time"})
			platform := args[0]
			for _, secret := range bm.ListByPlatform(platform) {
				// todo 默认只显示固定长度文本的 remark 字段，需要使用 info 命令查看详细信息
				tableData = append(tableData, []string{secret.Id, secret.Platform, secret.Account, secret.Remark, secret.CreateTime})
			}
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
	rootCmd.AddCommand(listCmd)
}
