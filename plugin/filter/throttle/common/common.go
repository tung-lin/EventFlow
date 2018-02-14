package common

//SettingConfig represents an event throttling policy configuration
type SettingConfig struct {
	TriggerCount int    `yaml:"triggercount"`
	ActionCount  int    `yaml:"actioncount"`
	PeriodSecond int    `yaml:"periodsecond"`
	Key          string `yaml:"key"`
}

//IThrottlingPolicy represents an interface for throttling policy
type IThrottlingPolicy interface {
	Throttling() (canExecute bool)
}
