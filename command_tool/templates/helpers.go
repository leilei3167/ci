package templates

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
	"unicode"

	flag "github.com/spf13/pflag"

	"github.com/spf13/cobra"
)

func (t *templater) templateFuncs(exposedFlags ...string) template.FuncMap {
	return template.FuncMap{
		"trim":                strings.TrimSpace,
		"trimRight":           func(s string) string { return strings.TrimRightFunc(s, unicode.IsSpace) },
		"trimLeft":            func(s string) string { return strings.TrimLeftFunc(s, unicode.IsSpace) },
		"gt":                  cobra.Gt,
		"eq":                  cobra.Eq,
		"rpad":                rpad,
		"appendIfNotPresent":  appendIfNotPresent,  // 如果s中有该字符串,则不做处理,否则添加
		"flagsNotIntersected": flagsNotIntersected, // 找出两个flagset中不相交的部分flag,单独构成一个flag
		"visibleFlags":        visibleFlags,        // 将help flag隐藏处理,构成新的一个flagset
		"flagsUsages":         flagsUsages,         // 将一个flagSet中的所有flag进行格式化打印
		"cmdGroups":           t.cmdGroups,         // 处理命令的分组
		"cmdGroupsString":     t.cmdGroupsString,   // 将分组的命令全部打印
		"rootCmd":             t.rootCmdName,
		"isRootCmd":           t.isRootCmd,
		"optionsCmdFor":       t.optionsCmdFor, // 找到到该命令的调用链
		"usageLine":           t.usageLine,     // 打印当前命令的Usageline,如 cmd [flags] [options]
		"exposed": func(c *cobra.Command) *flag.FlagSet { // 指定只展示哪些flag
			exposed := flag.NewFlagSet("exposed", flag.ContinueOnError)
			if len(exposedFlags) > 0 {
				for _, name := range exposedFlags {
					if flag := c.Flags().Lookup(name); flag != nil {
						exposed.AddFlag(flag)
					}
				}
			}
			return exposed
		},
	}
}

func (t *templater) rootCmdName(c *cobra.Command) string { // 根命令名称
	return t.rootCmd(c).CommandPath()
}

func (t *templater) isRootCmd(c *cobra.Command) bool { // 当前命令是否是根命令
	return t.rootCmd(c) == c
}

func (t *templater) rootCmd(c *cobra.Command) *cobra.Command {
	if c != nil && !c.HasParent() { // 没有父命令的就是根命令
		return c
	}
	if t.RootCmd == nil {
		panic("nil root cmd")
	}
	return t.RootCmd
}

// 处理一个命令和其所有子命令,进行分组,对未分组的命令加入到Other Commands分组中.
func (t *templater) cmdGroups(c *cobra.Command, all []*cobra.Command) []CommandGroup {
	if len(t.CommandGroups) > 0 && c == t.RootCmd {
		// 如果已有分组,并且c就是根命令,则直接返回命令组,将额外的一些命令添加到一个Other Commands中
		return AddAdditionalCommands(t.CommandGroups, "Other Commands:", all)
	}

	// 非根命令或当前命令没有分组
	all = filter(all, "options") // 将options命令过滤
	return []CommandGroup{
		{
			Name:     "Available Commands:",
			Commands: all,
		},
	}
}

// 实现将给定的分组格式化成字符串.
func (t *templater) cmdGroupsString(c *cobra.Command) string {
	groups := []string{}
	for _, cmdGroup := range t.cmdGroups(c, c.Commands()) { // 将一个命令和他所有的子命令传入
		cmds := []string{cmdGroup.Name} // 分组名
		for _, cmd := range cmdGroup.Commands {
			if cmd.IsAvailableCommand() { // 分组名\n+(空行)命令名(左对齐)+命令短介绍
				cmds = append(cmds, "  "+rpad(cmd.Name(), cmd.NamePadding())+" "+cmd.Short)
			}
		}
		groups = append(groups, strings.Join(cmds, "\n")) // 每个命令间换行
	}
	return strings.Join(groups, "\n\n") // 分组间用两个空行
}

func rpad(s string, padding int) string {
	// %5s:最小宽度为5
	// %-5s:最小宽度为5,并左对齐
	// %.5s:最大宽度为5
	// %5.7:最小宽度为5,最大宽度为7
	// %-5.7:最小宽度为5,最大宽度为7,左对齐
	// %5.3s:宽度大于3,截断
	// %05s:如果宽度小于5,会补零
	template := fmt.Sprintf("%%-%ds", padding) // %-11s,代表宽度必须为11,不足的会在后面补齐
	// 可以使的短名称命令也对其
	return fmt.Sprintf(template, s)
}

func (t *templater) usageLine(c *cobra.Command) string {
	usage := c.UseLine()
	suffix := "[options]"
	if c.HasFlags() && !strings.Contains(usage, suffix) { // 如果有flag,并且其usage中不包含[options]后缀
		usage += " " + suffix // 则添加后缀
	}
	return usage
}

func flagsUsages(f *flag.FlagSet) string {
	x := new(bytes.Buffer)

	f.VisitAll(func(flag *flag.Flag) {
		if flag.Hidden {
			return
		}
		format := "--%s=%s: %s\n" // 名称=默认值  Usage

		if flag.Value.Type() == "string" {
			format = "--%s='%s': %s\n"
		}

		if len(flag.Shorthand) > 0 { // 如果有简写,则先打印简写
			format = "  -%s, " + format
		} else {
			format = "   %s   " + format
		}

		fmt.Fprintf(x, format, flag.Shorthand, flag.Name, flag.DefValue, flag.Usage)
	})

	return x.String()
}

func flagsNotIntersected(l *flag.FlagSet, r *flag.FlagSet) *flag.FlagSet {
	f := flag.NewFlagSet("notIntersected", flag.ContinueOnError)
	l.VisitAll(func(fg *flag.Flag) { // 遍历l,找出其在r中不存在的flag
		if r.Lookup(fg.Name) == nil {
			f.AddFlag(fg)
		}
	})
	return f
}

func filter(cmds []*cobra.Command, names ...string) []*cobra.Command {
	out := []*cobra.Command{}
	for _, c := range cmds { // 将每一个命令与filter中的名字匹配
		if c.Hidden {
			continue
		}
		skip := false
		for _, name := range names {
			if name == c.Name() {
				skip = true
				break // 符合filter列表,退将此命令标记为过滤
			}
		}
		if skip {
			continue
		}
		out = append(out, c)

	}
	return out
}

func appendIfNotPresent(s, stringToAppend string) string {
	if strings.Contains(s, stringToAppend) { // 如果s中有该字符串,则不做处理
		return s
	}
	return s + " " + stringToAppend // 否则添加到s后面
}

func visibleFlags(l *flag.FlagSet) *flag.FlagSet {
	hidden := "help"
	f := flag.NewFlagSet("visible", flag.ContinueOnError)
	l.VisitAll(func(fg *flag.Flag) {
		if fg.Name != hidden {
			f.AddFlag(fg)
		}
	})
	return f
}

func (t *templater) optionsCmdFor(c *cobra.Command) string {
	if !c.Runnable() {
		return ""
	}

	rootCmdStructure := t.parents(c)
	for i := len(rootCmdStructure) - 1; i >= 0; i-- { // 从后向前遍历(即从root命令开始)
		cmd := rootCmdStructure[i]
		if _, _, err := cmd.Find([]string{"options"}); err == nil {
			return cmd.CommandPath() + " options"
		}
	}

	return ""
}

func (t *templater) parents(c *cobra.Command) []*cobra.Command {
	// 当前命令向上级遍历,获取调用链
	parents := []*cobra.Command{c}

	for current := c; !t.isRootCmd(current) && current.HasParent(); {
		current = current.Parent()
		parents = append(parents, current)
	}
	return parents
}
