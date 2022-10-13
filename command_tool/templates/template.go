package templates

func MainHelpTemplate() string {
	// or是一个函数,其其接受2个参数(此处参数为.Long和.Short)
	// or 函数会返回第一个非空的参数(此处即,如果设置了Long,则将long的参数返回)
	// with语句会判断`or .Long .Short`这一个pipeline产生的是否是空值,如果是,不会进行操作;
	// 否则将会执行到{{end}}之间的内容,并且将产生的值赋值给 . ,而 . 又将自己的值传递给了 trime函数做参数
	// if:如果pipeline产生的值为empty，不产生输出，否则输出T1执行结果。不改变dot的值
	return `{{with or .Long .Short}}{{. | trim}}{{end}}{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}`
}
