package line

import (
	"EventFlow/common/interface/pluginbase"
	"EventFlow/common/tool/yamltool"
	"log"
)

type SettingConfig struct {
	AccessToken string `yaml:accesstoken`
	Message     string `yaml:message`
}

type LineFactory struct {
}

func (factory LineFactory) GetIdentifyName() string {
	return "line"
}

func (factory LineFactory) CreateAction(config interface{}) pluginbase.IActionPlugin {

	var settingConfig SettingConfig
	err := yamltool.UnmarshalToType(config, &settingConfig)

	if err != nil {
		log.Print(err)
	}

	return &LinePlugin{Setting: settingConfig}
}
