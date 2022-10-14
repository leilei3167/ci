package templates

import (
	"fmt"
	"os"
	"text/template"

	"github.com/spf13/cobra"
)

// templater 负责将命令格式化输出.
type templater struct {
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

	templater := &templater{
		RootCmd:       cmd,
		CommandGroups: groups,
		Filtered:      filters,
		UsageTemplate: MainUsageTemplate(), // 模板语法
		HelpTemplate:  MainHelpTemplate(),  // 模板语法
	}

	// 修改flag错误时执行的操作,修改help命令的操作,修改Usage页面的展示方式
	cmd.SetFlagErrorFunc(templater.RootCmd.FlagErrorFunc())
	cmd.SetUsageFunc(templater.UsageFunc())
	cmd.SetHelpFunc(templater.HelpFunc())
}

func (t *templater) FlagErrorFunc(exposedFlags ...string) func(*cobra.Command, error) error {
	return func(c *cobra.Command, err error) error {
		c.SilenceUsage = true // 当flag错误时,不要直接打印Usage页面
		switch c.CalledAs() {
		case "options":
			return fmt.Errorf("%s\nrun '%s' without flags", err, c.CommandPath())
		default:
			return fmt.Errorf("%s\nsee '%s --help' for usage", err, c.CommandPath())
		}
	}
}

func (t *templater) UsageFunc(exposedFlags ...string) func(*cobra.Command) error {
	return func(c *cobra.Command) error {
		tt := template.New("usage")                // 创建模板引擎
		tt.Funcs(t.templateFuncs(exposedFlags...)) // 传入模板需要的函数,一般不会在模板中直接定义函数
		template.Must(tt.Parse(t.UsageTemplate))   // 解析指定的模板格式(Must是一个包装,确保执行成功)
		// out := term.NewResponsiveWriter(c.OutOrStderr())
		return tt.Execute(os.Stdout, c)
	}
}

func (t *templater) HelpFunc() func(*cobra.Command, []string) {
	return func(c *cobra.Command, s []string) {
		tt := template.New("help")
		tt.Funcs(t.templateFuncs())             // 传入模板中调用的函数
		template.Must(tt.Parse(t.HelpTemplate)) // 解析模板
		// TODO:适应terminal输出
		err := tt.Execute(c.OutOrStdout(), c) // 将cobra.Command作为模板的参数
		if err != nil {
			c.Println(err)
		}
	}
}

// UseOptionsTemplates.
func UseOptionsTemplates(cmd *cobra.Command) {
	templater := &templater{
		UsageTemplate: OptionsUsageTemplate(),
		HelpTemplate:  OptionsHelpTemplate(), // 空字符
	}
	cmd.SetUsageFunc(templater.UsageFunc())
	cmd.SetHelpFunc(templater.HelpFunc())
}
