package swaggerstats

import "EventFlow/common/interface/pluginbase"

type SwaggerStatusPlugin struct {
	Setting SettingConfig
}

var actionList []*pluginbase.IActionPlugin

func (trigger *SwaggerStatusPlugin) AddAction(actionPlugin *pluginbase.IActionPlugin) {
	actionList = append(actionList, actionPlugin)
}

func (trigger *SwaggerStatusPlugin) Start() {

}

func (trigger *SwaggerStatusPlugin) Stop() {

}
