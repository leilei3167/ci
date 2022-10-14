# cobra,pflag,viper构建命令行应用

## 需求点

1. 命令行工具的命令应该分组
2. 命令行工具应该提供flag的管理(全局flag,各个子命令的flag)
3. 子命令方便拓展
4. 自定义的help输出界面(模板输出,添加新的子命令时也能自适应)
5. 提供配置的灵活管理(flag数量巨大时需要使用配置文件提供默认的配置,通过flag显式指定的值优先级更高)

## 实现

以下`XXX`代表着二进制文件的名称

### NewXXXCommand 创建程序根命令

大致可以分为5步

1. 创建`cobra.Command`,并指定必要的字段,如`Use`,`Short`,`Long`,以及执行的动作`Run`,根命令建议运行函数直接设置为执行`HelpFunc`,这里可以指定一些钩子函数(如`PersistentPreRunE`,`PersistentPostRunE`),如在开始前执行pprof的初始化,在执行结束后进行pprof文件的清理
2. 初始化根命令的flagSet,`globalFlags := cmds.PersistentFlags()`,将此命令行工具所需要的通用的**config(应该有一个Config结构体来进行管理)**加入到这个flagSet中,这个flagSet将成为所有子命令可访问的全局flag
3. 配置文件的处理,大型的ctl工具的flag会非常之多,每一个都通过flag来指定是不可能的,因此一个ctl必须要支持配置文件,在配置文件中指定一些默认值,使得ctl在不指定多数flag的情况下也能正常运行
    - 使用viper的`_ = viper.BindPFlags(cmds.PersistentFlags())`,将之前定义的全局flag加入到viper中,再使用`cobra.OnInitialize(func() {}`来指定程序运行时的回调函数,回调函数的逻辑基本可以概括为:判断flag是否已经显式指定了配置文件->如果没有则在一些默认的地址搜索配置文件(指定一个默认的配置文件名)->调用`viper.ReadInConfig()`得到配置,程序中就可以使用`viper.Get`来获取某个配置了
    - 在viper没有找到配置文件时,仅仅打印错误以提示即可,不要终止程序
4. 构建子命令,并加入到根命令中
    - 构建SDK客户端接口(如果需要)
    - 子命令的管理交给`CommandGroup`这个结构来管理,其包含2个字段,一个是Name,分组的名字,一个是`Commands`包含这个分组下的所有命令
    - 对于根命令来说,命令分组是一个`CommandGroup`的切片
    - 每一个子命令都应该独立成包,暴露一个NewXXXCommand的方法(完成子命令的各种创建,flag管理等等)
    - 添加一些额外的命令,比如options命令(大型ctl都会将flag隐藏,执行options时才会显示)
    - 如果子命令需要特殊的输出模板,可在子命令提供的NewXXXCommand中指定(一般命令的格式和大部分不同才需要,如options命令仅需打印flag的情况即可)

5. 修改替换默认的ctl输出形式(重点)

这样做的好处就是,目录组织有条理,将同一类命令分组,每一个分组下的不同命令都是单独的包,并且构建一个命令的逻辑相同,便于后续拓展新的命令,还可以使用模板来生成命令代码脚手架,新命令就只需填充修改即可

最关键的是其中的第5步,将会用到**模板语法**来替换默认的帮助页面

### 替换根命令中的默认帮助页面的流程

默认的页面将会将所有的flag,命令,不做区分的打印出来,对于大型ctl来讲将会非常的冗杂,因此替换默认的显示是非常有必要的,示例中通过`func ActsAsRootCommand(cmd *cobra.Command, filters []string, groups ...CommandGroup)`函数来实现对一个根命令的修改,其接受一个根命令,一个filter(用于过滤一些不需要直接显示的命令),一个或多个命令分组(一个分组中有多个命令),修改过程如下:

1. 对于命令的格式化抽象为一个`templater`来负责处理,首先完成对他的初始化

```go
type templater struct {
    RootCmd       *cobra.Command 
    UsageTemplate string    //使用页面模板
    HelpTemplate  string    //帮助信息模板
    CommandGroups          // 各个命令分组
    Filtered      []string // 需要隐藏的命令
}
```

2. 指定Usage和Help的模板语法

[模板语法的使用](https://juejin.cn/post/6844903762901860360)

模板语法的使用流程概括:`tt := template.New("usage")`创建一个模板  
->`tt.Funcs(t.templateFuncs(exposedFlags...))`定义模板内可能需要的各种函数  
->`template.Must(tt.Parse(t.UsageTemplate))` 解析模板
->`tt.Execute(os.Stdout, c)`执行模板,需要指定输出的目的地,以及一个参数(参数在模板中会以 . 指代)

`MainUsageTemplate()`和`MainHelpTemplate()`分别定义了页面的打印逻辑

3. 替换根命令的页面默认函数

```go
// 修改flag错误时执行的操作,修改help命令的操作,修改Usage页面的展示方式
 cmd.SetFlagErrorFunc(templater.RootCmd.FlagErrorFunc())
 cmd.SetUsageFunc(templater.UsageFunc())
 cmd.SetHelpFunc(templater.HelpFunc())
```

而一个`templater.UsageFunc()`,就是对MainUsageTemplate的执行过程

### 实现模板输出的细节

关键点就在于模板的定义以及模板函数的定义

以`MainHelpTemplate`为例,其内容为`{{with or .Long .Short}}{{. | trim}}{{end}}{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}`

>or是一个函数,其其接受2个参数(此处参数为.Long和.Short)
>or 函数会返回第一个非空的参数(此处即,如果设置了Long,则将long的参数返回)
> with语句会判断`or .Long .Short`这一个pipeline产生的是否是空值,如果是,不会进行操作;**否则将会执行到{{end}}之间的内容,并且将产生的值赋值给 .**,而 . 又将自己的值传递给了 trime函数做参数
> if:如果pipeline产生的值为empty，不产生输出，否则输出T1执行结果。不改变dot的值

简单来说,会调用传入的`cobra.Command`的Long和Short(or 会先判断Long),如果其中一个有值,则会传入给trim函数(解析模板时传入的),实现Long的打印;之后又会判断其是否Runable或有子命令,如果有,则调用Usage方法,打印用法信息

最关键的是`MainUsageTemplate`的构建,由于帮助页面较大,可以将其分解为几个部分来进行模板组合

```go
func MainUsageTemplate() string {
	// Usage的内容比较长,因此拆分为几个部分来组装
	sections := []string{
		"\n\n",
		SectionVars,//声明模板中会使用的变量
		SectionAliases,
		SectionExamples,
		SectionSubCommands,
		SectionFlags,
		SectionTipsHelp,
		SectionTipsGlobalOptions,
	}
	return strings.TrimRightFunc(strings.Join(sections, ""), unicode.IsSpace)
}
```

1. SectionVars:这一部分主要是调用函数,将结果存到变量中
2. SectionAliases:判断命令是否有别名,如果有,则打印
3. SectionExamples:判断命令是否指定了用例,有则打印
4. SectionSubCommands:判断命令是否有子命令,有则打印(调用`cmdGroupsString`)
5. **SctionFlag**:判断是否有flag(此处flag经过flagsNotIntersected过滤,一定是全局flag中没有的部分),有则调用`flagsUsages`进行格式化打印(要确保对其等)
6. SectionTipsHelp: 在页面底部添加一条tips(如果该命令有子命令才执行打印)
7. SectionTipsGlobalOptions: 页面底部添加一条提示(提示如何查看全局的flag)

模板要使用的函数都在`helper.go`文件中
