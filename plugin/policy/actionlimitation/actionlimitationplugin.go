package actionlimitation

import (
	"sync/atomic"
	"time"
)

type ActionLimitationPlugin struct {
	tick                <-chan time.Time
	currentExecuteCount int32
	Setting             SettingConfig
}

func (plugin *ActionLimitationPlugin) Throttling(throttlingIdFromTrigger string) bool {

	if plugin.tick == nil {

		plugin.tick = time.Tick(time.Second * time.Duration(plugin.Setting.Second))

		go func() {
			for {
				select {
				case <-plugin.tick:
					atomic.StoreInt32(&plugin.currentExecuteCount, 0)
				}
			}
		}()
	}

	atomic.AddInt32(&plugin.currentExecuteCount, 1)
	canExecuteAction := atomic.LoadInt32(&plugin.currentExecuteCount) <= int32(plugin.Setting.Limitation)

	return canExecuteAction
}
