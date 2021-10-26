package config

import (
	"errors"
	"log"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Load load config file from given path
func Load(confPath string, out interface{}) error {
	v := viper.New()
	if confPath != "" {
		v.SetConfigFile(confPath) // 如果指定了配置文件，则解析指定的配置文件
	} else {
		v.AddConfigPath("config") // 如果没有指定配置文件，则解析默认的配置文件
		v.SetConfigName("config")
	}
	v.SetConfigType("yaml") // 设置配置文件格式为YAML

	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return errors.New("config file not found")
		}
		return err
	}
	err := v.Unmarshal(&out)
	if err != nil {
		return err
	}

	// watch file change
	go func() {
		v.WatchConfig()
		v.OnConfigChange(func(e fsnotify.Event) {
			_ = v.Unmarshal(&out)
			log.Printf("Config file changed: %s", e.Name)
		})
	}()

	return nil
}
