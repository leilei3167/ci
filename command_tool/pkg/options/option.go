// Package options 提供通用的一些配置.
package options

import (
	"flag"
	"io"
	"sync"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/leilei3167/ci/command_tool/templates"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const ( // 此处定义每个flag的名称,分组
	FlagConfig      = "XXXconfig"
	FlagBearerToken = "user.token"
	FlagUsername    = "user.username"
	FlagPassword    = "user.password"

	FlagInsecure      = "server.insecure-skip-tls-verify"
	FlagTimeout       = "server.timeout"
	FlagMaxRetries    = "server.max-retries"
	FlagRetryInterval = "server.retry-interval"
)

// ConfigFlags 代表着命令行工具所需要的配置,全部采用指针字段,方便判断是否为空.
type ConfigFlags struct { // 如果是调用SDK的命令行工具,此处可按需要添加必要的Flag
	XXXConfig *string

	// Auth
	BearerToken *string
	Username    *string
	Password    *string
	// SecretID    *string
	// SecretKey   *string

	// TLS
	Insecure *bool
	// TLSServerName *string
	// CertFile      *string
	// KeyFile       *string
	// CAFile        *string

	// connection
	Timeout       *time.Duration
	MaxRetries    *int
	RetryInterval *time.Duration

	usePersistentConfig bool

	lock sync.Mutex
}

// NewConfigFlags returns ConfigFlags with default values set.
func NewConfigFlags(usePersistentConfig bool) *ConfigFlags {
	return &ConfigFlags{ // 添加一些默认的flag值
		XXXConfig:           pointer.ToString(""),
		BearerToken:         pointer.ToString(""),
		Insecure:            pointer.ToBool(false),
		Timeout:             pointer.ToDuration(30 * time.Second),
		MaxRetries:          pointer.ToInt(0),
		RetryInterval:       pointer.ToDuration(1 * time.Second),
		usePersistentConfig: usePersistentConfig,
	}
}

// 开关式的添加flag选项.
func (f *ConfigFlags) WithDeprecatedPasswordFlag() *ConfigFlags {
	f.Username = pointer.ToString("")
	f.Password = pointer.ToString("")
	return f
}

// AddFlags 添加所有的不为nil的选项到flagset中.
func (f *ConfigFlags) AddFlags(flags *pflag.FlagSet) {
	if f.XXXConfig != nil {
		flags.StringVar(f.XXXConfig, FlagConfig, *f.XXXConfig, "配置文件路径")
	}

	if f.BearerToken != nil {
		flags.StringVar(
			f.BearerToken,
			FlagBearerToken,
			*f.BearerToken,
			"Bearer token for authentication to the API server",
		)
	}

	if f.Username != nil {
		flags.StringVar(f.Username, FlagUsername, *f.Username, "Username for basic authentication to the API server")
	}

	if f.Password != nil {
		flags.StringVar(f.Password, FlagPassword, *f.Password, "Password for basic authentication to the API server")
	}

	if f.Insecure != nil {
		flags.BoolVar(f.Insecure, FlagInsecure, *f.Insecure, ""+
			"If true, the server's certificate will not be checked for validity. This will make your HTTPS connections insecure")
	}

	if f.Timeout != nil {
		flags.DurationVar(
			f.Timeout,
			FlagTimeout,
			*f.Timeout,
			"The length of time to wait before giving up on a single server request. Non-zero values should contain a corresponding time unit (e.g. 1s, 2m, 3h). A value of zero means don't timeout requests.",
		)
	}

	if f.MaxRetries != nil {
		flag.IntVar(f.MaxRetries, FlagMaxRetries, *f.MaxRetries, "Maximum number of retries.")
	}

	if f.RetryInterval != nil {
		flags.DurationVar(
			f.RetryInterval,
			FlagRetryInterval,
			*f.RetryInterval,
			"The interval time between each attempt.",
		)
	}
}

var optionsExample = `
		# Print flags inherited by all commands
		iamctl options`

// NewCmdOptions implements the options command.
// 其执行的就是将全局的flag打印出来.
func NewCmdOptions(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "options",
		Short:   "Print the list of flags inherited by all commands",
		Long:    "Print the list of flags inherited by all commands",
		Example: optionsExample,
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Usage()
		},
	}

	// The `options` command needs write its output to the `out` stream
	// (typically stdout). Without calling SetOutput here, the Usage()
	// function call will fall back to stderr.
	cmd.SetOutput(out)
	// 指定options命令的输出模板
	templates.UseOptionsTemplates(cmd)

	return cmd
}
