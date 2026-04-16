package cmd

import (
	"github.com/follow1123/secret-book/bookmanager"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:          "delete id",
	Aliases:      []string{"rm"},
	Short:        "delete secret",
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true, // 关闭错误时的帮助信息
	GroupID:      cmdGrpDefault,
	RunE: func(cmd *cobra.Command, args []string) error {
		// todo 实现强制删除功能，直接删除，不保留到历史列表内
		// todo 添加 -h 选项，实现根据id删除历史列表内的数据
		bm, err := bookmanager.New(bookmanager.DefaultSecretsFile())
		if err != nil {
			return err
		}
		id := args[0]
		if err := bm.DeleteByIdPrefix(id); err != nil {
			return err
		}
		if err := bm.Save(); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
