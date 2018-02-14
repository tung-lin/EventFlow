package policy

import (
	"EventFlow/plugin/filter/throttle/common"
	"sync/atomic"
	"time"
)

//An TriggerThreshold represents throttling parameters
type TriggerThreshold struct {
	timer               *time.Timer
	currentTriggerCount int32
	currentExecuteCount int32
	Setting             common.SettingConfig
}

//NewTriggerThreshold create a new event throttling policy
func NewTriggerThreshold(setting common.SettingConfig) *TriggerThreshold {
	return &TriggerThreshold{Setting: setting}
}

//Throttling decides whether to fire action or not
func (plugin *TriggerThreshold) Throttling() bool {

	if plugin.timer == nil {
		plugin.timer = time.AfterFunc(time.Second*time.Duration(plugin.Setting.PeriodSecond), func() {
			atomic.StoreInt32(&plugin.currentExecuteCount, 0)
			atomic.StoreInt32(&plugin.currentTriggerCount, 0)
		})
	} else {
		if atomic.LoadInt32(&plugin.currentExecuteCount) == 0 && atomic.LoadInt32(&plugin.currentTriggerCount) == 0 {
			plugin.timer.Reset(time.Second * time.Duration(plugin.Setting.PeriodSecond))
		}
	}

	atomic.AddInt32(&plugin.currentTriggerCount, 1)
	canExecuteAction := atomic.LoadInt32(&plugin.currentTriggerCount) >= int32(plugin.Setting.TriggerCount) && atomic.LoadInt32(&plugin.currentExecuteCount) < int32(plugin.Setting.ActionCount)

	if canExecuteAction {
		atomic.AddInt32(&plugin.currentExecuteCount, 1)
	}

	return canExecuteAction
}
