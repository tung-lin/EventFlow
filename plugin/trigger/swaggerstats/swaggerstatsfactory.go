package swaggerstats

import (
	"EventFlow/common/interface/pluginbase"
	"EventFlow/common/tool/yamltool"
	"log"
)

//yaml config
type SettingConfig struct {
	Swagger_url string `yaml:swagger_url`
	API_IP      string `yaml:api_ip`
	API_Path    string `yaml:api_path`
	Operations  []struct {
		OperationId   string   `yaml:operaionid`
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
