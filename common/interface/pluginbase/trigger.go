package pluginbase

type ITriggerFactory interface {
	GetIdentifyName() string
	CreateTrigger(config interface{}) ITriggerPlugin
}

type ITriggerPlugin interface {
	IActionHandler
	Start()
	Stop()
}

type IActionHandler interface {
	ActionHandleFunc(f func(triggerPlugin *ITriggerPlugin, messageFromTrigger *string))
	FireAction(triggerPlugin *ITriggerPlugin, messageFromTrigger *string)
}

type ActionHandler struct {
	actionHandleFunc func(triggerPlugin *ITriggerPlugin, messageFromTrigger *string)
}

func (h *ActionHandler) ActionHandleFunc(f func(triggerPlugin *ITriggerPlugin, messageFromTrigger *string)) {
	h.actionHandleFunc = f
}

func (h *ActionHandler) FireAction(triggerPlugin *ITriggerPlugin, messageFromTrigger *string) {
	if h.actionHandleFunc != nil {
		h.actionHandleFunc(triggerPlugin, messageFromTrigger)
	}
}
