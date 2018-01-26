package swaggerstats

import (
	"EventFlow/common/interface/pluginbase"
	"EventFlow/common/tool/yamltool"
	"log"
)

//yaml config
type SettingConfig struct {
	SwaggerURL     string   `yaml:swagger_url`
	APIIP          string   `yaml:api_ip`
	APIPath        string   `yaml:api_path`
	SkipOperations []string `yaml:skip_operations`
	Operations     []struct {
		OperationID   string   `yaml:operaionid`
		Condition     string   `yaml:condition`
		Fields        []string `yaml:fields`
		ThresholdType string   `yaml:thresholdtype`
		Threshold     string   `yaml:threshold`
	} `yaml:operations`
}

type SwaggerStatsFactory struct {
}

func (factory SwaggerStatsFactory) GetIdentifyName() string {
	return "swaggerstats"
}

func (factory SwaggerStatsFactory) CreateTrigger(config interface{}) pluginbase.ITriggerPlugin {

	var settingConfig SettingConfig
	err := yamltool.UnmarshalToType(config, &settingConfig)

	if err != nil {
		log.Print(err)
	}

	return &SwaggerStatusPlugin{Setting: settingConfig}
}
