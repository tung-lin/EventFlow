package alwaysfired

type AlwaysFiredPlugin struct {
}

func (plugin AlwaysFiredPlugin) FireAction(throttlingIdFromTrigger string) bool {
	return true
}
