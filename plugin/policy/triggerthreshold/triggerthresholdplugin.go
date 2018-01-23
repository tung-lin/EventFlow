package triggerthreshold

type TriggerThresholdPlugin struct {
	Setting SettingConfig
}

func (plugin TriggerThresholdPlugin) FireAction(throttlingIdFromTrigger string) bool {
	return true
}
