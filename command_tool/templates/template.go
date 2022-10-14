package templates

import (
	"strings"
	"unicode"
)

func MainHelpTemplate() string {
	// or是一个函数,其其接受2个参数(此处参数为.Long和.Short)
	// or 函数会返回第一个非空的参数(此处即,如果设置了Long,则将long的参数返回)
	// with语句会判断`or .Long .Short`这一个pipeline产生的是否是空值,如果是,不会进行操作;
	// 否则将会执行到{{end}}之间的内容,并且将产生的值赋值给 . ,而 . 又将自己的值传递给了 trime函数做参数
	// if:如果pipeline产生的值为empty，不产生输出，否则输出T1执行结果。不改变dot的值
	return `{{with or .Long .Short}}{{. | trim}}{{end}}{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}`
}

func MainUsageTemplate() string {
	// Usage的内容比较长,因此拆分为几个部分来组装
	sections := []string{
		"\n\n",
		SectionVars,
		SectionAliases,
		SectionExamples,
		SectionSubCommands,
		SectionFlags,
		SectionTipsHelp,
		SectionTipsGlobalOptions,
	}
	return strings.TrimRightFunc(strings.Join(sections, ""), unicode.IsSpace)
}

const (
	// SectionVars 定义模板内使用的变量
	// 调用外部函数,将 . 作为参数.
	// $isRootCmd判断当前才command是否是根命令
	// $rootCmd 返回了当前templates中的根命令
	//$visibleFlags保存了可见的flag(先将全局flag和命令的flag进行正交,全局flag没有的会显示).
	//$explicitlyExposedFlags 指定了需要特别展示的有哪些flag
	//$optionsCmdFor 从当前命令的根命令往下找,最近的具有options命令的路径
	//$usageLine 当前命令的usage line,如果当前命令有flag并且没有[options]则为其添加.
	SectionVars = `{{$isRootCmd := isRootCmd .}}` +
		`{{$rootCmd := rootCmd .}}` +
		`{{$visibleFlags := visibleFlags (flagsNotIntersected .LocalFlags .PersistentFlags)}}` +
		`{{$explicitlyExposedFlags := exposed .}}` +
		`{{$optionsCmdFor := optionsCmdFor .}}` +
		`{{$usageLine := usageLine .}}`

		// 这一部分会判断当前的命令是否有别名,如果有,则输出
		// 调用的是gt函数,cobra.Gt会比较第一个参数是否大于第二个参数,.Aliases
	SectionAliases = `{{if gt .Aliases 0}}Aliases:
{{.NameAndAliases}}

{{end}}`
	// 如果当前命令设置了用例,则打印.
	SectionExamples = `{{if .HasExample}}Examples:
{{trimRight .Example}}

{{end}}`
	// 当前命令有子命令时,按分组打印命令及其short介绍,分组之间用\n\n进行分隔.
	SectionSubCommands = `{{if .HasAvailableSubCommands}}{{cmdGroupsString .}}
	
{{end}}`
	// 打印命令的flag,
	// 1.会先从$visibleFlags  $explicitlyExposedFlags两个FlagSet中调用HasFlags方法判断是否有内容.
	// 2.如有内容则打印flag的使用页面.
	// TODO:{{if $visibleFlags.HasFlags}}马上{{end}}是什么意思.
	SectionFlags = `{{if or $visibleFlags.HasFlags $explicitlyExposedFlags.HasFlags}}Options:
{{if $visibleFlags.HasFlags}}{{trimRight (flagsUsages $visibleFlags)}}{{end}}{{if $explicitlyExposedFlags.HasFlags}}{{if $visibleFlags.HasFlags}}
{{end}}{{trimRight (flagsUsages $explicitlyExposedFlags)}}{{end}}
	
{{end}}`

	SectionUsage = `{{if and .Runnable (ne .UsaLine "")(ne .UseLine $rootCmd)}}Usage:
   {{$usageLine}}

{{end}}
	`

	SectionTipsHelp = `{{if .HasSubCommands}}Use "{{$rootCmd}} <command> --help" for more information about a given command.
{{end}}`

	SectionTipsGlobalOptions = `{{if $optionsCmdFor}}Use "{{$optionsCmdFor}}" for a list of global command-line options (applies to all commands).
{{end}}`
)

func OptionsHelpTemplate() string {
	return ""
}

// OptionsUsageTemplate 将会打印options命令的页面.
func OptionsUsageTemplate() string {
	return `{{ if .HasInheritedFlags}}The following options can be passed to any command:

{{flagsUsages .InheritedFlags}}{{end}}`
}
