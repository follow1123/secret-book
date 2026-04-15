package cmd

import (
	"fmt"
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
		bm, err := bookmanager.New(bookmanager.DefaultSecretsFile())
		if err != nil {
			return err
		}
		table := tablewriter.NewTable(os.Stdout, tablewriter.WithEastAsian(tw.Off))
		tableData := [][]string{}

		id := args[0]
		secretMap := bm.GetByIdPerfix(id)

		if len(secretMap) == 0 {
			return fmt.Errorf("no secret id %s", id)
		} else if len(secretMap) > 1 {
			return fmt.Errorf("duplicated id prefix %s", id)
		} else {
			for _, secret := range secretMap {
				tableData = append(tableData, []string{"ID", secret.Id})
				tableData = append(tableData, []string{"PLATFORM", secret.Platform})
				tableData = append(tableData, []string{"ACCOUNT", secret.Account})
				tableData = append(tableData, []string{"PASSWORD", secret.Password})
				tableData = append(tableData, []string{"REMARK", secret.Remark})
				tableData = append(tableData, []string{"CREATE TIME", secret.CreateTime})
			}
		}

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
