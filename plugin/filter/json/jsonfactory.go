package json

import (
	"EventFlow/common/interface/pluginbase"
	"EventFlow/common/tool/yamltool"
	"log"
)

type SettingConfig struct {
	AddMetadata map[string]string `yaml:"addmetadata"`
}

type JSONFactory struct {
}

func (factory JSONFactory) GetIdentifyName() string {
	return "json"
}

func (factory JSONFactory) CreateFilter(config interface{}) pluginbase.IFilterPlugin {

	var settingConfig SettingConfig
	err := yamltool.UnmarshalToType(config, &settingConfig)

	if err != nil {
		log.Print(err)
	}

	return &JSONPlugin{Setting: settingConfig}
}
