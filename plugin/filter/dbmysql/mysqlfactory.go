package dbmysql

import (
	"EventFlow/common/interface/pluginbase"
	"EventFlow/common/tool/cachetool"
	"EventFlow/common/tool/yamltool"
	"log"
)

type SettingConfig struct {
	IP          string            `yaml:ip`
	User        string            `yaml:user`
	Password    string            `yaml:password`
	Database    string            `yaml:database`
	Command     string            `yaml:command`
	AddMetadata map[string]string `yaml:addmetadata`
	Cache       cachetool.Cache   `yaml:cache`
}

type MySQLFactory struct {
}

func (factory MySQLFactory) GetIdentifyName() string {
	return "mysql"
}

func (factory MySQLFactory) CreateFilter(config interface{}) pluginbase.IFilterPlugin {

	var settingConfig SettingConfig
	err := yamltool.UnmarshalToType(config, &settingConfig)

	if err != nil {
		log.Print(err)
	}

	return &MySQLPlugin{Setting: settingConfig}
}
