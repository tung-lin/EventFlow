package policy

import (
	"EventFlow/plugin/filter/throttle/common"
	"sync"
	"sync/atomic"
	"time"
)

type TriggerThreshold struct {
	timer               *time.Timer
	mutex               *sync.Mutex
	currentTriggerCount int32
	currentExecuteCount int32
	Setting             common.SettingConfig
}

func NewTriggerThreshold(setting common.SettingConfig) *TriggerThreshold {
	return &TriggerThreshold{Setting: setting, mutex: &sync.Mutex{}}
}

func (plugin *TriggerThreshold) Throttling() bool {

	plugin.mutex.Lock()

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

	plugin.mutex.Unlock()

	atomic.AddInt32(&plugin.currentTriggerCount, 1)
	canExecuteAction := atomic.LoadInt32(&plugin.currentTriggerCount) >= int32(plugin.Setting.TriggerCount) && atomic.LoadInt32(&plugin.currentExecuteCount) < int32(plugin.Setting.ActionCount)

	if canExecuteAction {
		atomic.AddInt32(&plugin.currentExecuteCount, 1)
	}

	return canExecuteAction
}
