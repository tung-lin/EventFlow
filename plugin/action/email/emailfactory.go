package email

import (
	"EventFlow/common/interface/pluginbase"
	"EventFlow/common/tool/yamltool"
	"log"
)

type SettingConfig struct {
	Address   string      `yaml:"address"`
	Port      int         `yaml:"port"`
	Username  string      `yaml:"username"`
	Password  string      `yaml:"password"`
	From      string      `yaml:"from"`
	StaticTo  interface{} `yaml:"staticto"`
	DynamicTo string      `yaml:"dynamicto"`
	Subject   string      `yaml:"subject"`
	Body      string      `yaml:"body"`
}

type EmailFactory struct {
}

func (factory EmailFactory) GetIdentifyName() string {
	return "email"
}

func (factory EmailFactory) CreateAction(config interface{}) pluginbase.IActionPlugin {

	var settingConfig SettingConfig
	err := yamltool.UnmarshalToType(config, &settingConfig)

	if err != nil {
		log.Print(err)
	}

	return &EmailPlugin{Setting: settingConfig}
}
