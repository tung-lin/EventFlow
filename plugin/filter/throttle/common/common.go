package common

type SettingConfig struct {
	TriggerCount int    `yaml:triggercount`
	ActionCount  int    `yaml:actioncount`
	PeriodSecond int    `yaml:periodsecond`
	Key          string `yaml:key`
}

type IThrottlingPolicy interface {
	Throttling() (canExecute bool)
}
