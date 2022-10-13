package templates

import (
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

// Templater 负责将命令格式化输出.
type Templater struct {
	RootCmd       *cobra.Command
	UsageTemplate string
	HelpTemplate  string
	CommandGroups          // 各个命令分组
	Filtered      []string // 需要隐藏的命令
}

func ActsAsRootCommand(cmd *cobra.Command, filters []string, groups ...CommandGroup) {
	if cmd == nil {
		panic("nil root command")
	}

	templater := &Templater{
		RootCmd:       cmd,
		CommandGroups: groups,
		Filtered:      filters,
		UsageTemplate: "",                 // 模板语法
		HelpTemplate:  MainHelpTemplate(), // 模板语法
	}

	// 修改cmd的Usage函数等
	// cmd.SetFlagErrorFunc() //flag错误时,执行

	// cmd.SetUsageFunc()
	cmd.SetHelpFunc(templater.HelpFunc())
}

func (t *Templater) HelpFunc() func(*cobra.Command, []string) {
	return func(c *cobra.Command, s []string) {
		tt := template.New("help")
		tt.Funcs(t.templateFuncs())             // 传入模板中调用的函数
		template.Must(tt.Parse(t.HelpTemplate)) // 解析模板
		// TODO:适应terminal输出
		err := tt.Execute(c.OutOrStdout(), c)
		if err != nil {
			c.Println(err)
		}
	}
}

func (t *Templater) templateFuncs() template.FuncMap {
	return template.FuncMap{
		"trim": strings.TrimSpace,
	}
}
