package actionlimitation

type ActionLimitationPlugin struct {
	Setting SettingConfig
}

func (plugin ActionLimitationPlugin) FireAction(throttlingIdFromTrigger string) bool {
	return true
}
