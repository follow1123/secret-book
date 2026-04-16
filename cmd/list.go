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
			secrets := bm.ListByPlatform(platform)
			for _, secret := range secrets {
				tableData = append(tableData, []string{secret.Id[:6], secret.Platform, secret.Account, truncateText(secret.Remark, 10), secret.CreateTime})
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

func truncateText(text string, maxLen int) string {
	// 获取字符串的 rune 切片（按字符切分）
	runes := []rune(text)

	// 如果字符长度小于 maxLen，直接返回原始字符串
	if len(runes) <= maxLen {
		return text
	}

	// 截取前 maxLen 个字符
	return string(runes[:maxLen])
}
