package pluginbase

//IFilterFactory interface for filter factory
type IFilterFactory interface {
	GetIdentifyName() string
	CreateFilter(config interface{}) IFilterPlugin
}

//IFilterPlugin interface for filter plugin
type IFilterPlugin interface {
	DoFilter(messageFromTrigger *string, parameters *map[string]interface{}) (doNextPipeline bool)
}
