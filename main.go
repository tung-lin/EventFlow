package main

import (
	"EventFlow/common/interface/pluginbase"
	"EventFlow/common/tool/throttlingtool"
	"io/ioutil"
	"log"
	"os"

	yaml "gopkg.in/yaml.v2"
)

type IPluginLoader interface {
	Load() (triggerFactories map[string]pluginbase.ITriggerFactory, actionFactories map[string]pluginbase.IActionFactory)
}

type EventFlow struct {
	Trigger []struct {
		Mode    string      `yaml:mode`
		Setting interface{} `yaml:setting`
	} `yaml:trigger`

	Action []struct {
		Mode   string `yaml:mode`
		Policy struct {
			Mode    string      `yaml:mode`
			Setting interface{} `yaml:setting`
		} `yaml:policy`
		Setting interface{} `yaml:setting`
	} `yaml:action`
}

var loader IPluginLoader
var triggerFactoryMap map[string]pluginbase.ITriggerFactory
var actionFactoryMap map[string]pluginbase.IActionFactory

func main() {

	loader = PluginImportLoader{}
	//loader = PluginSharedObjectLoader{}

	triggerFactoryMap, actionFactoryMap = loader.Load()
	LoadConfig()
}

func LoadConfig() {

	currentePath, _ := os.Getwd()
	configPath := currentePath + "/config/"
	files, err := ioutil.ReadDir(configPath)

	if err != nil {
		log.Print(err)
		return
	}

	var allTriggerPluginList []*pluginbase.ITriggerPlugin

	for _, file := range files {
		pipelineFile, err := ioutil.ReadFile(configPath + file.Name())

		if err != nil {
			log.Fatalf("Read config file failed: %v", err)
			continue
		}

		var config EventFlow
		err = yaml.Unmarshal(pipelineFile, &config)

		if err != nil {
			log.Fatalf("Unmarshal config file failed: %v", err)
			continue
		}

		var triggerPluginList []*pluginbase.ITriggerPlugin

		for _, trigger := range config.Trigger {
			triggerFactory, existed := triggerFactoryMap[trigger.Mode]

			if existed {
				triggerPlugin := triggerFactory.CreateTrigger(trigger.Setting)

				if triggerPlugin != nil {
					triggerPluginList = append(triggerPluginList, &triggerPlugin)
					allTriggerPluginList = append(allTriggerPluginList, &triggerPlugin)
				}
			}
		}

		for _, action := range config.Action {
			actionFactory, existed := actionFactoryMap[action.Mode]

			if existed {

				if action.Policy != nil {
					policy := throttlingtool.CreatePolicy(action.Policy.Mode, action.Policy.Setting)
				}

				actionPlugin := actionFactory.CreateAction(action.Setting)

				if actionPlugin != nil {
					for _, triggerPlugin := range triggerPluginList {
						(*triggerPlugin).AddAction(&actionPlugin)
					}
				}
			}
		}

		for _, triggerPlugin := range triggerPluginList {
			(*triggerPlugin).Start()
		}
	}
}
