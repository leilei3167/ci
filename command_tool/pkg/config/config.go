package config

import (
	"log"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

func LoadConfig(cfg string, defaultName string) {
	// 是否指定了配置文件路径
	if cfg != "" {
		viper.SetConfigFile(cfg)
	} else {
		// 未指定,添加几个默认路径去按照默认配置文件名查找
		viper.AddConfigPath(".")                                          // 1.当前的文件夹
		viper.AddConfigPath(filepath.Join(HomeDir(), RecommendedHomeDir)) // $HOME/.xxx 文件夹
		viper.AddConfigPath("/etc/xxx")

		viper.SetConfigName(defaultName)
	}

	viper.SetConfigType("yaml")
	viper.AutomaticEnv() // 处理环境变量
	// 设置环境变量前缀,建议是大写程序名,Viper会自动捕获带有此前缀的环境变量
	viper.SetEnvPrefix(RecommendedEnvPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_")) // 替换为_

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("WARNING: viper failed to discover and load the configuration file: %s", err.Error())
	}
}

const (
	// RecommendedHomeDir defines the default directory used to place all iam service configurations.
	RecommendedHomeDir = ".iam"

	// RecommendedEnvPrefix defines the ENV prefix used by all iam service.
	RecommendedEnvPrefix = "IAM"
)
