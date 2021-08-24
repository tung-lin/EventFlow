package basic

import (
	"EventFlow/common/tool/yamltool"
	"EventFlow/plugin/trigger/httppoll/authbase"
	"log"
)

type BasicFactory struct {
}

type basicSettingConfig struct {
	Username   string            `yaml:"username"`
	Password   string            `yaml:"password"`
	AuthHeader string            `yaml:"authheader"`
	Headers    []authbase.Header `yaml:"headers"`
}

func (factory BasicFactory) GetIdentifyName() string {
	return "basichttp"
}

func (factory BasicFactory) CreateAuth(config interface{}) authbase.IAuthPlugin {
	var settingConfig basicSettingConfig
	err := yamltool.UnmarshalToType(config, &settingConfig)

	if err != nil {
		log.Print(err)
	}

	return &BasicPlugin{Setting: settingConfig}
}
