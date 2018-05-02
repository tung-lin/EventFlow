package udp

import (
	"EventFlow/common/interface/pluginbase"
	"EventFlow/common/tool/yamltool"
	"log"
)

type SettingConfig struct {
	Port int `yaml:"port"`
}

type UdpFactory struct {
}

func (factory UdpFactory) GetIdentifyName() string {
	return "udp"
}

func (factory UdpFactory) CreateTrigger(config interface{}) pluginbase.ITriggerPlugin {

	var settingConfig SettingConfig
	err := yamltool.UnmarshalToType(config, &settingConfig)

	if err != nil {
		log.Print(err)
	}

	return &UDPPlugin{Setting: settingConfig}
}
