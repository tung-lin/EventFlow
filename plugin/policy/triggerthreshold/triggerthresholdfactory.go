package triggerthreshold

import (
	"EventFlow/common/interface/pluginbase"
	"EventFlow/common/tool/yamltool"
	"log"
)

type SettingConfig struct {
	Second     int `yaml:second`
	Threshold  int `yaml:threshold`
	Limitation int `yaml:limitation`
}

type TriggerThresholdFactory struct {
}

func (factory TriggerThresholdFactory) IsDefaultFactory() bool {
	return false
}

func (factory TriggerThresholdFactory) GetIdentifyName() string {
	return "triggerthreshold"
}

func (factory TriggerThresholdFactory) CreatePolicy(config interface{}) pluginbase.IThrottlingPolicyPlugin {

	var settingConfig SettingConfig
	err := yamltool.UnmarshalToType(config, &settingConfig)

	if err != nil {
		log.Print(err)
	}

	return &TriggerThresholdPlugin{Setting: settingConfig}
}
