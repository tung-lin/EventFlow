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

	key := parametertool.ReplaceWithParameter(&filter.Setting.Key, parameters)

	filter.mutex.Lock()

	throttlePolicy, existed := filter.throttlePolicyMap[key]

	if !existed {
		if filter.Setting.TriggerCount <= 0 {
			throttlePolicy = policy.NewActionLimitation(filter.Setting)
		} else {
			throttlePolicy = policy.NewTriggerThreshold(filter.Setting)
		}
		filter.throttlePolicyMap[key] = throttlePolicy
	}

	filter.mutex.Unlock()

	return throttlePolicy.Throttling()
}
