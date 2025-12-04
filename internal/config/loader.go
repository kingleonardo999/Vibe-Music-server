package config

import (
	"errors"
	"github.com/spf13/viper"
)

func load() (*viper.Viper, error) {
	v := viper.New()

	// 读取 templateConfig.yaml 文件
	v.AddConfigPath("./config")
	v.SetConfigName("templateConfig")
	v.SetConfigType("yaml")
	err1 := v.MergeInConfig()

	// 读取 config.yaml 文件
	v.AddConfigPath("./config")
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	err2 := v.MergeInConfig()

	if err1 != nil && err2 != nil {
		return nil, errors.New("failed to load config files")
	}

	return v, nil
}
