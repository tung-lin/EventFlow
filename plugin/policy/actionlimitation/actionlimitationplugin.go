package actionlimitation

import (
	"sync"
	"sync/atomic"
	"time"
)

type ActionLimitationPlugin struct {
	timer               *time.Timer
	mutex               *sync.Mutex
	currentExecuteCount int32
	Setting             SettingConfig
}

func (plugin *ActionLimitationPlugin) Throttling(throttlingIdFromTrigger string) bool {

	plugin.mutex.Lock()

	if plugin.timer == nil {
		plugin.timer = time.AfterFunc(time.Second*time.Duration(plugin.Setting.Second), func() {
			atomic.StoreInt32(&plugin.currentExecuteCount, 0)
		})
	} else {
		if atomic.LoadInt32(&plugin.currentExecuteCount) == 0 {
			plugin.timer.Reset(time.Second * time.Duration(plugin.Setting.Second))
		}
	}

	plugin.mutex.Unlock()

	atomic.AddInt32(&plugin.currentExecuteCount, 1)
	canExecuteAction := atomic.LoadInt32(&plugin.currentExecuteCount) <= int32(plugin.Setting.Limitation)

	return canExecuteAction
}
