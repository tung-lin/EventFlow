package basic

import "EventFlow/plugin/trigger/httppoll/authbase"

type BasicPlugin struct {
	Setting basicSettingConfig
}

func (auth *BasicPlugin) CheckParameter() error {
	return nil
}

func (auth *BasicPlugin) CreateHttpHeaders() []authbase.Header {
	return nil
}
