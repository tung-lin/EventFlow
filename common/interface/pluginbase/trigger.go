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
	PolicyHandleFunc(f func(triggerPlugin *ITriggerPlugin, throttlingID string, messageFromTrigger *string))
	FireAction(triggerPlugin *ITriggerPlugin, messageFromTrigger *string)
	FireActionWithThrottlingID(triggerPlugin *ITriggerPlugin, throttlingID string, messageFromTrigger *string)
}

type PolicyHandler struct {
	throttlingID     string
	policyHandleFunc func(triggerPlugin *ITriggerPlugin, throttlingID string, messageFromTrigger *string)
}

func (h *PolicyHandler) PolicyHandleFunc(f func(triggerPlugin *ITriggerPlugin, throttlingID string, messageFromTrigger *string)) {
	h.throttlingID = stringtool.CreateRandomString()
	h.policyHandleFunc = f
}

func (h *PolicyHandler) FireAction(triggerPlugin *ITriggerPlugin, messageFromTrigger *string) {
	if h.policyHandleFunc != nil {
		h.policyHandleFunc(triggerPlugin, h.throttlingID, messageFromTrigger)
	}
}

func (h *PolicyHandler) FireActionWithThrottlingID(triggerPlugin *ITriggerPlugin, throttlingID string, messageFromTrigger *string) {
	if h.policyHandleFunc != nil {
		h.policyHandleFunc(triggerPlugin, throttlingID, messageFromTrigger)
	}
}
