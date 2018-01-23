package pluginbase

type ITriggerFactory interface {
	GetIdentifyName() string
	CreateTrigger(config interface{}) ITriggerPlugin
}

type ITriggerPlugin interface {
	AddAction(actionPlugin *IActionPlugin)
	Start()
	Stop()
}
