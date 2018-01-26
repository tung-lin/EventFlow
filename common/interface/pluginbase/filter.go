package pluginbase

type IFilterFactory interface {
	GetIdentifyName() string
	CreateFilter(config interface{}) IFilterPlugin
}

type IFilterPlugin interface {
	DoFilter(messageFromTrigger *string, parameters *map[string]interface{}) (doNextPipeline bool)
}
