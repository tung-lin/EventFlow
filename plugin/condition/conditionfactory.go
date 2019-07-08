package condition

import (
	"EventFlow/common/interface/pluginbase"
	"EventFlow/common/tool/pipelinetool"
)

type ConditionFactory struct {
}

func (factory ConditionFactory) GetIdentifyName() string {
	return "condition"
}

func (factory ConditionFactory) CreateCondition(config pipelinetool.Condition) pluginbase.IConditionPlugin {

	return &ConditionPlugin{Setting: config}
}
