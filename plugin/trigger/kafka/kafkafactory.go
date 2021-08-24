package kafka

import (
	"EventFlow/common/interface/pluginbase"
	"EventFlow/common/tool/yamltool"
	"log"
)

type SettingConfig struct {
	Brokers    interface{} `yaml:"brokers"`
	GroupID    string      `yaml:"groupid"`
	Topic      string      `yaml:"topic"`
	FromOldest bool        `yaml:"from_oldest"`
}

type KafkaFactory struct {
}

func (factory KafkaFactory) GetIdentifyName() string {
	return "kafka"
}

func (factory KafkaFactory) CreateTrigger(config interface{}) pluginbase.ITriggerPlugin {

	var settingConfig SettingConfig
	err := yamltool.UnmarshalToType(config, &settingConfig)

	if err != nil {
		log.Print(err)
	}

	return &KafkaPlugin{Setting: settingConfig}
}
