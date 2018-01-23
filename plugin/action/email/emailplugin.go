package email

type EmailPlugin struct {
	Setting SettingConfig
}

func (trigger *EmailPlugin) FireAction(throttlingIdFromTrigger string, messageFromTrigger *string) {

}
