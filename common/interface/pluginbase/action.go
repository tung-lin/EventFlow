package pluginbase

type IActionFactory interface {
	GetIdentifyName() string
	CreateAction(config interface{}) IActionPlugin
}

type IActionPlugin interface {
	FireAction(messageFromTrigger *string, parameters *map[string]interface{})
}
