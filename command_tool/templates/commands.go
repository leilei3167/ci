package templates

import "github.com/spf13/cobra"

type CommandGroup struct { // 一个命令分组
	Name     string           // 分组名称
	Commands []*cobra.Command // 子命令列表
}

type CommandGroups []CommandGroup // 一个程序有多个命令分组

func (g CommandGroups) Add(c *cobra.Command) {
	for _, group := range g {
		c.AddCommand(group.Commands...)
	}
}
