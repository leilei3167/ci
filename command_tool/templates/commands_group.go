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

func (g CommandGroups) Has(c *cobra.Command) bool {
	for _, group := range g {
		for _, command := range group.Commands {
			if command == c {
				return true
			}
		}
	}
	return false
}

// AddAdditionalCommands 向已有的一个分组中额外添加一个分组,这个分组的成员是cmds有而g没有的命令.
func AddAdditionalCommands(g CommandGroups, message string, cmds []*cobra.Command) CommandGroups {
	group := CommandGroup{Name: message}
	for _, c := range cmds {
		// Don't show commands that have no short description
		if !g.Has(c) && len(c.Short) != 0 {
			group.Commands = append(group.Commands, c)
		}
	}
	if len(group.Commands) == 0 {
		return g
	}
	return append(g, group)
}
