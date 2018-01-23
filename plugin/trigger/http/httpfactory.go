package http

import (
	"EventFlow/common/interface/pluginbase"
	"EventFlow/common/tool/yamltool"
	"log"
)

type SettingConfig struct {
	Port int `yaml:port`
}

type HttpFactory struct {
}

func (factory HttpFactory) GetIdentifyName() string {
	return "http"
}

func (factory HttpFactory) CreateTrigger(config interface{}) pluginbase.ITriggerPlugin {

	var settingConfig SettingConfig
	err := yamltool.UnmarshalToType(config, &settingConfig)

	if err != nil {
		log.Print(err)
	}

	return &HttpPlugin{Setting: settingConfig}
}
