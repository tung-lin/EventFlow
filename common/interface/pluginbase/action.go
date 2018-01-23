package pluginbase

type IActionFactory interface {
	GetIdentifyName() string
	CreateAction(config interface{}) IActionPlugin
}

type IActionPlugin interface {
	FireAction(triggerPlugin ITriggerPlugin, messageFromTrigger string, throttlingIdFromTrigger string)
}
