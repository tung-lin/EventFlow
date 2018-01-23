package alwaysfired

import "EventFlow/common/interface/pluginbase"

type AlwaysFiredFactory struct {
}

func (factory AlwaysFiredFactory) IsDefaultFactory() bool {
	return true
}

func (factory AlwaysFiredFactory) GetIdentifyName() string {
	return "alwaysfired"
}

func (factory AlwaysFiredFactory) CreatePolicy(config interface{}) pluginbase.IThrottlingPolicy {
	return AlwaysFiredPlugin{}
}
