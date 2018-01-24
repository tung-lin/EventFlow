package actionlimitation

import (
	"EventFlow/common/interface/pluginbase"
	"EventFlow/common/tool/yamltool"
	"log"
	"sync"
)

type SettingConfig struct {
	Second     int `yaml:second`
	Limitation int `yaml:limitation`
}

type ActionLimitationFactory struct {
}

func (factory ActionLimitationFactory) IsDefaultFactory() bool {
	return false
}

func (factory ActionLimitationFactory) GetIdentifyName() string {
	return "actionlimitation"
}

func (factory ActionLimitationFactory) CreatePolicy(config interface{}) pluginbase.IThrottlingPolicyPlugin {

	var settingConfig SettingConfig
	err := yamltool.UnmarshalToType(config, &settingConfig)

	if err != nil {
		log.Print(err)
	}

	return &ActionLimitationPlugin{Setting: settingConfig, mutex: &sync.Mutex{}}
}
