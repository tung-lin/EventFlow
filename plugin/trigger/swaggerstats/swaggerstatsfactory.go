package swaggerstats

import (
	"EventFlow/common/interface/pluginbase"
	"EventFlow/common/tool/yamltool"
	"log"
	"strings"
)

//SettingConfig represents a configurations for api status monitoring
type SettingConfig struct {
	SwaggerFile        string      `yaml:"swagger_file"`
	SwaggerURL         string      `yaml:"swagger_url"`
	APIIP              string      `yaml:"api_ip"`
	APIPath            string      `yaml:"api_path"`
	MonitorIntervalSec int         `yaml:"monitor_interval_sec"`
	MonitorTimeoutSec  int         `yaml:"monitor_timeout_sec"`
	NotifyHTTPOK       bool        `yaml:"notify_http_ok"`
	ODataTop           int         `yaml:"odata_top"`
	ODataFormat        string      `yaml:"odata_format"`
	SkipOperations     []string    `yaml:"skip_operations"`
	Operations         []operation `yaml:"operations"`
}

type operation struct {
	OperationID        string   `yaml:"operationid"`
	Condition          string   `yaml:"condition"`
	MonitorIntervalSec int      `yaml:"monitor_interval_sec"`
	MonitorTimeoutSec  int      `yaml:"monitor_timeout_sec"`
	ODataTop           int      `yaml:"odata_top"`
	ODataFormat        string   `yaml:"odata_format"`
	Fields             []string `yaml:"fields"`
	ThresholdType      string   `yaml:"thresholdtype"`
	Threshold          string   `yaml:"threshold"`
	Parameters         []struct {
		Values []struct {
			Name  string      `yaml:"name"`
			value interface{} `yaml:"value"`
		} `yaml:"values"`
	} `yaml:"parameters"`
}

type SwaggerStatsFactory struct {
}

func (factory SwaggerStatsFactory) GetIdentifyName() string {
	return "swaggerstats"
}

func (factory SwaggerStatsFactory) CreateTrigger(config interface{}) pluginbase.ITriggerPlugin {

	var settingConfig SettingConfig
	err := yamltool.UnmarshalToType(config, &settingConfig)

	if err != nil {
		log.Print(err)
	} else {

		if settingConfig.ODataFormat != "" && !strings.EqualFold(settingConfig.ODataFormat, "json") && !strings.EqualFold(settingConfig.ODataFormat, "xml") {
			settingConfig.ODataFormat = "json"
		}

		if settingConfig.MonitorIntervalSec < 20 {
			settingConfig.MonitorIntervalSec = 20
		}

		if settingConfig.MonitorTimeoutSec <= 0 {
			settingConfig.MonitorTimeoutSec = 10
		}

		if settingConfig.MonitorTimeoutSec > settingConfig.MonitorIntervalSec {
			settingConfig.MonitorTimeoutSec = settingConfig.MonitorIntervalSec - 1
		}

		if settingConfig.ODataTop > 20 {
			settingConfig.ODataTop = 1
		}
	}

	return &SwaggerStatusPlugin{Setting: settingConfig}
}
