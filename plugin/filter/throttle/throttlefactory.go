package throttle

import (
	"EventFlow/common/interface/pluginbase"
	"EventFlow/common/tool/stringtool"
	"EventFlow/common/tool/yamltool"
	"EventFlow/plugin/filter/throttle/common"
	"log"
	"sync"
)

type ThrottleFactory struct {
}

func (factory ThrottleFactory) GetIdentifyName() string {
	return "throttle"
}

func (factory ThrottleFactory) CreateFilter(config interface{}) pluginbase.IFilterPlugin {

	var settingConfig common.SettingConfig
	err := yamltool.UnmarshalToType(config, &settingConfig)

	if err != nil {
		log.Print(err)
	}

	return &ThrottlePlugin{
		Setting:           settingConfig,
		throttleKey:       stringtool.CreateRandomString(),
		throttlePolicyMap: make(map[string]common.IThrottlingPolicy),
		mutex:             &sync.Mutex{}}
}
