package pluginbase

type ITriggerFactory interface {
	IFactoryBase
	GetIdentifyName() string
	CreateTrigger(config interface{}) ITriggerPlugin
}

type ITriggerPlugin interface {
	AddAction(actionPlugin *IActionPlugin)
}
