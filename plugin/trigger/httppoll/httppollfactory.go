package httppoll

import (
	"EventFlow/common/interface/pluginbase"
	"EventFlow/common/tool/yamltool"
	"log"
)

type HttpPollFactory struct {
}

type SettingConfig struct {
	Urls      []Url  `yaml:"urls"`
	Cron      string `yaml:"cron"`
	TimeoutMS int    `yaml:"timeout_ms"`
	Codec     string `yaml:"codec"`
	Auth      Auth   `yaml:"auth"`
}

type Url struct {
	Url   string `yaml:"url"`
	Codec string `yaml:"codec"`
}

type Auth struct {
	Mode    string      `yaml:"mode"`
	Setting interface{} `yaml:"setting"`
}

func (factory HttpPollFactory) GetIdentifyName() string {
	return "http_poll"
}

func (factory HttpPollFactory) CreateTrigger(config interface{}) pluginbase.ITriggerPlugin {

	var settingConfig SettingConfig
	err := yamltool.UnmarshalToType(config, &settingConfig)

	if err != nil {
		log.Print(err)
	}

	timeoutMS := &settingConfig.TimeoutMS
	codec := &settingConfig.Codec

	if timeoutMS == nil {
		settingConfig.TimeoutMS = 1000
	}

	if codec == nil {
		settingConfig.Codec = "json"
	}

	return &HttpPollPlugin{Setting: settingConfig}
}
