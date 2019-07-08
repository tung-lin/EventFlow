package pluginbase

import "EventFlow/common/tool/pipelinetool"

// IConditionFactory interface for condition factory
type IConditionFactory interface {
	GetIdentifyName() string
	CreateCondition(config pipelinetool.Condition) IConditionPlugin
}

//IConditionPlugin interface for condition plugin
type IConditionPlugin interface {
	IsMatch(parameters *map[string]interface{}) bool
}
