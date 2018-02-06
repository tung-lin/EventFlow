package pluginbase

// IActionFactory interface for action factory
type IActionFactory interface {
	GetIdentifyName() string
	CreateAction(config interface{}) IActionPlugin
}

//IActionPlugin interface for aciont plugin
type IActionPlugin interface {
	FireAction(messageFromTrigger *string, parameters *map[string]interface{})
}
