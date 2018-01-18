package main

import (
	"EventFlow/common/interface/pluginbase"
)

func main() {

}

var Factory SwaggerStatsFactory

type SwaggerStatsFactory struct {
	pluginbase.FactoryBase
}

func (factory *SwaggerStatsFactory) GetIdentifyName() string {
	return "swagger-stats"
}

func (factory *SwaggerStatsFactory) CreateTrigger(config interface{}) pluginbase.ITriggerPlugin {

	var settingConfig SettingConfig
	factory.UnmarshalToType(config, &settingConfig)

	return &SwaggerStatus{Setting: settingConfig}
}

type SettingConfig struct {
	Swagger_url string `yaml:swagger_url`
	API_IP      string `yaml:api_ip`
	API_Path    string `yaml:api_path`
	Operations  []struct {
		OperationId    string   `yaml:operationid`
		Condition      string   `yaml:condition`
		Fields         []string `yaml:fields`
		Threshold_Type string   `yaml:threshold_type`
		Threshold      string   `yaml:threshold`
	} `yaml:operations`
}

type SwaggerStatus struct {
	Setting SettingConfig
}

func (trigger *SwaggerStatus) AddAction(actionPlugin *pluginbase.IActionPlugin) {

}
