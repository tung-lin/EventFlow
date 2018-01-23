package email

import (
	"EventFlow/common/interface/pluginbase"
)

type EmailPlugin struct {
	Setting SettingConfig
}

func (trigger *EmailPlugin) FireAction(triggerPlugin pluginbase.ITriggerPlugin, messageFromTrigger string, throttlingIdFromTrigger string) {

}
