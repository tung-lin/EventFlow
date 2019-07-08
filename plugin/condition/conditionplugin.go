package condition

import (
	"EventFlow/common/tool/pipelinetool"
)

type ConditionPlugin struct {
	Setting pipelinetool.Condition
}

func (condition *ConditionPlugin) IsMatch(parameters *map[string]interface{}) bool {

	if ok, err := pipelinetool.IsMatchCondition(&condition.Setting, parameters); !ok && len(condition.Setting) > 0 {
		return false
	}

	return true
}
