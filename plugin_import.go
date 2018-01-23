package main

import (
	"EventFlow/common/interface/pluginbase"
	"EventFlow/common/tool/throttlingtool"
	"EventFlow/plugin/action/email"
	"EventFlow/plugin/policy/actionlimitation"
	"EventFlow/plugin/policy/alwaysfired"
	"EventFlow/plugin/policy/triggerthreshold"
	"EventFlow/plugin/trigger/http"
	"EventFlow/plugin/trigger/swaggerstats"
)

type PluginImportLoader struct {
}

func (loader PluginImportLoader) Load() (triggerFactories map[string]pluginbase.ITriggerFactory, actionFactories map[string]pluginbase.IActionFactory) {

	//create trigger factories
	var triggerFactoryMap = make(map[string]pluginbase.ITriggerFactory)
	triggers := []pluginbase.ITriggerFactory{}

	triggers = append(triggers, swaggerstats.SwaggerStatsFactory{}, http.HttpFactory{})

	for _, trigger := range triggers {
		triggerFactoryMap[trigger.GetIdentifyName()] = trigger
	}

	//create action factories
	var actionFactoryMap = make(map[string]pluginbase.IActionFactory)
	actions := []pluginbase.IActionFactory{}

	actions = append(actions, email.EmailFactory)

	for _, action := range actions {
		actionFactoryMap[action.GetIdentifyName()] = action
	}

	//create policy factories
	throttlingtool.AddPolicyFactory(alwaysfired.AlwaysFiredFactory{})
	throttlingtool.AddPolicyFactory(triggerthreshold.TriggerThresholdFactory{})
	throttlingtool.AddPolicyFactory(actionlimitation.ActionLimitationFactory{})

	return triggerFactoryMap, actionFactoryMap
}
