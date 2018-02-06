package pluginbase

//ITriggerFactory interfacr for factory
type ITriggerFactory interface {
	GetIdentifyName() string
	CreateTrigger(config interface{}) ITriggerPlugin
}

//ITriggerPlugin interface for trigger plugin
type ITriggerPlugin interface {
	IActionHandler
	Start()
	Stop()
}

//IActionHandler interface for action handler
type IActionHandler interface {
	ActionHandleFunc(f func(triggerPlugin *ITriggerPlugin, messageFromTrigger *string))
	FireAction(triggerPlugin *ITriggerPlugin, messageFromTrigger *string)
}

//ActionHandler struct for action handler
type ActionHandler struct {
	actionHandleFunc func(triggerPlugin *ITriggerPlugin, messageFromTrigger *string)
}

//ActionHandleFunc add handler function
func (h *ActionHandler) ActionHandleFunc(f func(triggerPlugin *ITriggerPlugin, messageFromTrigger *string)) {
	h.actionHandleFunc = f
}

//FireAction fire action
func (h *ActionHandler) FireAction(triggerPlugin *ITriggerPlugin, messageFromTrigger *string) {
	if h.actionHandleFunc != nil {
		h.actionHandleFunc(triggerPlugin, messageFromTrigger)
	}
}
