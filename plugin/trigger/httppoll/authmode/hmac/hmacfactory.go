package hmac

import (
	"EventFlow/common/tool/yamltool"
	"EventFlow/plugin/trigger/httppoll/authbase"
	"log"
)

type HMacFactory struct {
}

type hmacSettingConfig struct {
	Type      string            `yaml:"type"`
	APPID     string            `yaml:"appid"`
	APPKey    string            `yaml:"appkey"`
	Algorithm string            `yaml:"algorithm"`
	Headers   []authbase.Header `yaml:"headers"`
}

func (factory HMacFactory) GetIdentifyName() string {
	return "hmac"
}

func (factory HMacFactory) CreateAuth(config interface{}) authbase.IAuthPlugin {
	var settingConfig hmacSettingConfig
	err := yamltool.UnmarshalToType(config, &settingConfig)

	if err != nil {
		log.Print(err)
	}

	return &HMacPlugin{Setting: settingConfig}
}
