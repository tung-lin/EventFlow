package policy

import (
	"EventFlow/plugin/filter/throttle/common"
	"sync/atomic"
	"time"
)

type ActionLimitation struct {
	timer               *time.Timer
	currentExecuteCount int32
	Setting             common.SettingConfig
}

func NewActionLimitation(setting common.SettingConfig) *ActionLimitation {
	return &ActionLimitation{Setting: setting}
}

func (plugin *ActionLimitation) Throttling() bool {

	if plugin.timer == nil {
		plugin.timer = time.AfterFunc(time.Second*time.Duration(plugin.Setting.PeriodSecond), func() {
			atomic.StoreInt32(&plugin.currentExecuteCount, 0)
		})
	} else {
		if atomic.LoadInt32(&plugin.currentExecuteCount) == 0 {
			plugin.timer.Reset(time.Second * time.Duration(plugin.Setting.PeriodSecond))
		}
	}

	atomic.AddInt32(&plugin.currentExecuteCount, 1)
	canExecuteAction := atomic.LoadInt32(&plugin.currentExecuteCount) <= int32(plugin.Setting.ActionCount)

	return canExecuteAction
}
