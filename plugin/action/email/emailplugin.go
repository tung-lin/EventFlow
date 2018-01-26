package email

type EmailPlugin struct {
	Setting SettingConfig
}

func (trigger *EmailPlugin) FireAction(messageFromTrigger *string, parameters *map[string]interface{}) {

}
