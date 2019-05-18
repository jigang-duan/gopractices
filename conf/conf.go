package conf

import (
	"github.com/fsnotify/fsnotify"
	"github.com/jigang-duan/gopractices/bootstrap"
	"github.com/spf13/viper"
)

func Configure(b *bootstrap.Bootstrapper) {
	if b.ConfigFile != "" {
		viper.SetConfigFile(b.ConfigFile)
	} else {
		viper.AddConfigPath("conf")
		viper.SetConfigName("config")
	}

	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		b.Logger().Fatalf("加载配置文件时出错 v%", err)
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		b.Logger().Printf("Config file changed: %s\n", e.Name)
	})

	b.AppName = viper.GetString("app.name")
	b.AppOwner = viper.GetString("app.owner")
}
