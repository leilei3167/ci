package cmd

import (
	"io"
	"os"

	"github.com/leilei3167/ci/command_tool/cmd/info"
	"github.com/leilei3167/ci/command_tool/pkg/config"
	"github.com/leilei3167/ci/command_tool/pkg/options"
	"github.com/leilei3167/ci/command_tool/templates"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewXXXComand 返回程序的根命令,直接执行会运行helpFunc.
func NewXXXComand(in io.Reader, out, err io.Writer) *cobra.Command {
	// 1.创建根命令
	cmds := &cobra.Command{
		Use:   os.Args[0],                         // 本命令行工具的名字,可设为os.Args[0]
		Short: "XXXctl controls the XXX platform", // 简短介绍
		Long: `XXXctl controls the XXX platform xxx xxx xxx xxx xxx
		And xx xxx xxx xxx`, // 长描述需要处理规范//TODO:处理rawString,规范输出
		Run: runHelp, // 根命令直接输出帮助函数
		// 运行前和运行后的钩子函数
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return initProfiling()
		}, // 最优先执行的函数,PersistentPreRun->PreRun->Run->PostRun->PersistentPostRun
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			return flushProfiling()
		},
	}

	// 2.初始化处理根命令的flagSet
	globalFlags := cmds.PersistentFlags()
	// globalFlags.SetNormalizeFunc(func(f *pflag.FlagSet, name string) pflag.NormalizedName { return pflag.NormalizedName("") }) //
	// 设置标志的转换函数,可用于纠错,或者提醒弃用等 globalFlags.SetNormalizeFunc(func(f *pflag.FlagSet, name string)
	// pflag.NormalizedName { return pflag.NormalizedName("") }) // 比如用户可能将-分隔符输入成_,可以通过此函数翻译

	// 3.向flagSet中添加选项
	// 3.1 通用的分析选项(pprof)
	addProfilingFlags(globalFlags)

	// 3.2添加工具所需的配置选项(Options需要同时支持从配置文件读取和从flag中指定)
	opts := options.NewConfigFlags(true).WithDeprecatedPasswordFlag()
	opts.AddFlags(globalFlags)

	// 4.配置文件的定位与解析
	// 4.1首先,将已加入cmd的flag的值,添加到Viper中;flag设置的值优先级更高
	_ = viper.BindPFlags(cmds.PersistentFlags())
	cobra.OnInitialize(func() { // 在最最开始,及Execute被调用时执行,读取配置文件的值,会被flag的值覆盖
		config.LoadConfig(viper.GetString(options.FlagConfig), "xxxc") // 传入默认配置路径
	})

	// 5. 可在此根据配置初始化SDK客户端等,或者工具所需的客户端等等

	// 6.添加子命令
	// 6.1 指定输出
	ioStreams := options.IOStreams{In: in, Out: out, ErrOut: err}

	// 6.2初始化子命令,并添加到根命令
	groups := templates.CommandGroups{
		{
			Name: "Basic Commands:",
			Commands: []*cobra.Command{
				// 此处添加各个分组下的命令,每一个子命令的构建方式应该都是相同的
				// 可自动生成
				info.NewInfoCommand(ioStreams),
			},
		},
		{
			Name:     "Identity and Access Management Commands:",
			Commands: []*cobra.Command{
				// 添加命令
			},
		},
		// 可以添加更多分组
	}
	groups.Add(cmds)

	// TODO:此处修改cmds的展示,替换默认的展示页面,实现分组展示的关键
	filters := []string{"options"} // 不显示的命令名称
	templates.ActsAsRootCommand(cmds, filters, groups...)

	// 添加一些公用的命令(可隐藏)
	cmds.AddCommand(options.NewCmdOptions(ioStreams.Out))
	return cmds
}

func runHelp(c *cobra.Command, args []string) {
	_ = c.Help()
}
