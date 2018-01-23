package triggerthreshold

import (
	"sync/atomic"
	"time"
)

type TriggerThresholdPlugin struct {
	tick                <-chan time.Time
	currentTriggerCount int32
	currentExecuteCount int32
	Setting             SettingConfig
}

func (plugin *TriggerThresholdPlugin) Throttling(throttlingIdFromTrigger string) bool {

	if plugin.tick == nil {

		plugin.tick = time.Tick(time.Second * time.Duration(plugin.Setting.Second))

		go func() {
			for {
				select {
				case <-plugin.tick:
					atomic.StoreInt32(&plugin.currentExecuteCount, 0)
					atomic.StoreInt32(&plugin.currentTriggerCount, 0)
				}
			}
		}()
	}

	atomic.AddInt32(&plugin.currentTriggerCount, 1)
	canExecuteAction := atomic.LoadInt32(&plugin.currentTriggerCount) >= int32(plugin.Setting.Threshold) && atomic.LoadInt32(&plugin.currentExecuteCount) <= int32(plugin.Setting.Limitation)

	if canExecuteAction {
		atomic.AddInt32(&plugin.currentExecuteCount, 1)
	}

	return canExecuteAction
}
