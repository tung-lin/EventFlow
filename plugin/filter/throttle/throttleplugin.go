package throttle

import (
	"EventFlow/common/tool/parametertool"
	"EventFlow/plugin/filter/throttle/common"
	"EventFlow/plugin/filter/throttle/policy"
	"sync"
)

type ThrottlePlugin struct {
	Setting           common.SettingConfig
	mutex             *sync.Mutex
	throttleKey       string
	throttlePolicyMap map[string]common.IThrottlingPolicy
	parameterReplaced bool
}

func (filter *ThrottlePlugin) DoFilter(messageFromTrigger *string, parameters *map[string]interface{}) (doNextPipeline bool) {

	filter.mutex.Lock()

	if !filter.parameterReplaced {
		filter.parameterReplaced = true
		parametertool.ReplaceWithParameter(&filter.Setting.Key, parameters)
	}

	throttlePolicy, existed := filter.throttlePolicyMap[filter.Setting.Key]

	if !existed {
		if filter.Setting.TriggerCount <= 0 {
			throttlePolicy = policy.NewActionLimitation(filter.Setting)
		} else {
			throttlePolicy = policy.NewTriggerThreshold(filter.Setting)
		}
		filter.throttlePolicyMap[filter.Setting.Key] = throttlePolicy
	}

	filter.mutex.Unlock()

	return throttlePolicy.Throttling()
}
