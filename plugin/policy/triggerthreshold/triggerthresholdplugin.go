package triggerthreshold

import (
	"sync"
	"sync/atomic"
	"time"
)

type TriggerThresholdPlugin struct {
	timer               *time.Timer
	mutex               *sync.Mutex
	currentTriggerCount int32
	currentExecuteCount int32
	Setting             SettingConfig
}

func (plugin *TriggerThresholdPlugin) Throttling(throttlingIdFromTrigger string) bool {

	plugin.mutex.Lock()

	if plugin.timer == nil {
		plugin.timer = time.AfterFunc(time.Second*time.Duration(plugin.Setting.Second), func() {
			atomic.StoreInt32(&plugin.currentExecuteCount, 0)
			atomic.StoreInt32(&plugin.currentTriggerCount, 0)
		})
	} else {
		if atomic.LoadInt32(&plugin.currentExecuteCount) == 0 && atomic.LoadInt32(&plugin.currentTriggerCount) == 0 {
			plugin.timer.Reset(time.Second * time.Duration(plugin.Setting.Second))
		}
	}

	plugin.mutex.Unlock()

	atomic.AddInt32(&plugin.currentTriggerCount, 1)
	canExecuteAction := atomic.LoadInt32(&plugin.currentTriggerCount) >= int32(plugin.Setting.Threshold) && atomic.LoadInt32(&plugin.currentExecuteCount) <= int32(plugin.Setting.Limitation)

	if canExecuteAction {
		atomic.AddInt32(&plugin.currentExecuteCount, 1)
	}

	return canExecuteAction
}
