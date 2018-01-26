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
}

func (filter *ThrottlePlugin) DoFilter(messageFromTrigger *string, parameters *map[string]interface{}) (doNextPipeline bool) {

	throttleValue, existed := (*parameters)[parametertool.GetParameterKey(filter.Setting.Key)]

	if !existed {
		throttleValue = filter.throttleKey
	}

	throttleKey := throttleValue.(string)

	filter.mutex.Lock()

	throttlePolicy, existed := filter.throttlePolicyMap[throttleKey]

	if !existed {
		if filter.Setting.TriggerCount <= 0 {
			throttlePolicy = policy.NewActionLimitation(filter.Setting)
		} else {
			throttlePolicy = policy.NewTriggerThreshold(filter.Setting)
		}
		filter.throttlePolicyMap[throttleKey] = throttlePolicy
	}

	filter.mutex.Unlock()

	return throttlePolicy.Throttling()
}
