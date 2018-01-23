package swaggerstats

import "EventFlow/common/interface/pluginbase"

type SwaggerStatusPlugin struct {
	pluginbase.PolicyHandler
	Setting SettingConfig
}

func (trigger *SwaggerStatusPlugin) Start() {

}

func (trigger *SwaggerStatusPlugin) Stop() {

}
