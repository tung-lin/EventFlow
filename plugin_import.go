package main

import (
	"EventFlow/common/interface/pluginbase"
	"EventFlow/plugin/action/email"
	"EventFlow/plugin/action/line"
	"EventFlow/plugin/filter/dbmysql"
	"EventFlow/plugin/filter/json"
	"EventFlow/plugin/filter/throttle"
	"EventFlow/plugin/trigger/http"
	"EventFlow/plugin/trigger/swaggerstats"
)

type PluginImportLoader struct {
}

func (loader PluginImportLoader) Load() (triggerFactories map[string]pluginbase.ITriggerFactory, filterFactories map[string]pluginbase.IFilterFactory, actionFactories map[string]pluginbase.IActionFactory) {

	//create trigger factories
	var triggerFactoryMap = make(map[string]pluginbase.ITriggerFactory)
	triggers := []pluginbase.ITriggerFactory{}

	triggers = append(triggers, swaggerstats.SwaggerStatsFactory{}, http.HttpFactory{})

	for _, trigger := range triggers {
		triggerFactoryMap[trigger.GetIdentifyName()] = trigger
	}

	//create filter factories
	var filterFactoryMap = make(map[string]pluginbase.IFilterFactory)
	filters := []pluginbase.IFilterFactory{}

	filters = append(filters, json.JSONFactory{}, throttle.ThrottleFactory{}, dbmysql.MySQLFactory{})

	for _, filter := range filters {
		filterFactoryMap[filter.GetIdentifyName()] = filter
	}

	//create action factories
	var actionFactoryMap = make(map[string]pluginbase.IActionFactory)
	actions := []pluginbase.IActionFactory{}

	actions = append(actions, email.EmailFactory{}, line.LineFactory{})

	for _, action := range actions {
		actionFactoryMap[action.GetIdentifyName()] = action
	}

	return triggerFactoryMap, filterFactoryMap, actionFactoryMap
}
