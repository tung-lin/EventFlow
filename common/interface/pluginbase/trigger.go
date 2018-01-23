package pluginbase

import "EventFlow/common/tool/stringtool"

type ITriggerFactory interface {
	GetIdentifyName() string
	CreateTrigger(config interface{}) ITriggerPlugin
}

type ITriggerPlugin interface {
	IPolicyHandler
	Start()
	Stop()
}

type IPolicyHandler interface {
	PolicyHandleFunc(f func(triggerPlugin *ITriggerPlugin, throttlingId string, messageFromTrigger *string))
	FireAction(triggerPlugin *ITriggerPlugin, messageFromTrigger *string)
	FireActionWithCustomId(triggerPlugin *ITriggerPlugin, throttlingId string, messageFromTrigger *string)
}

type PolicyHandler struct {
}

var policyHandleFunc func(triggerPlugin *ITriggerPlugin, throttlingId string, messageFromTrigger *string)
var throttlingId string

func (h *PolicyHandler) PolicyHandleFunc(f func(triggerPlugin *ITriggerPlugin, throttlingId string, messageFromTrigger *string)) {
	policyHandleFunc = f
}

func (h *PolicyHandler) FireAction(triggerPlugin *ITriggerPlugin, messageFromTrigger *string) {
	if policyHandleFunc != nil {
		if throttlingId == "" {
			throttlingId = stringtool.CreateRandomString()
		}
		policyHandleFunc(triggerPlugin, throttlingId, messageFromTrigger)
	}
}

func (h *PolicyHandler) FireActionWithCustomId(triggerPlugin *ITriggerPlugin, customTriggerId string, messageFromTrigger *string) {
	if policyHandleFunc != nil {
		policyHandleFunc(triggerPlugin, customTriggerId, messageFromTrigger)
	}
}
