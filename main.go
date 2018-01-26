package main

import (
	"EventFlow/common/interface/pluginbase"
	"io/ioutil"
	"log"
	"os"

	yaml "gopkg.in/yaml.v2"
)

type IPluginLoader interface {
	Load() (triggerFactories map[string]pluginbase.ITriggerFactory, filterFactories map[string]pluginbase.IFilterFactory, actionFactories map[string]pluginbase.IActionFactory)
}

type EventFlow struct {
	Trigger struct {
		Mode    string      `yaml:mode`
		Setting interface{} `yaml:setting`
	} `yaml:trigger`

	Filter []struct {
		Mode    string      `yaml:mode`
		Setting interface{} `yaml:setting`
	} `yaml:filter`

	Action []struct {
		Mode    string      `yaml:mode`
		Setting interface{} `yaml:setting`
	} `yaml:action`
}

type ThrottlingPolicy struct {
	Mode    string      `yaml:mode`
	Setting interface{} `yaml:setting`
}

var loader IPluginLoader
var ch chan bool

var triggerFactoryMap map[string]pluginbase.ITriggerFactory
var actionFactoryMap map[string]pluginbase.IActionFactory
var filterFactoryMap map[string]pluginbase.IFilterFactory

var pipelineFilterMap map[pluginbase.ITriggerPlugin][]pluginbase.IFilterPlugin
var pipelineActionMap map[pluginbase.ITriggerPlugin][]pluginbase.IActionPlugin

func init() {
	loader = PluginImportLoader{}
	//loader = PluginSharedObjectLoader{}

	triggerFactoryMap, filterFactoryMap, actionFactoryMap = loader.Load()

	pipelineFilterMap = make(map[pluginbase.ITriggerPlugin][]pluginbase.IFilterPlugin)
	pipelineActionMap = make(map[pluginbase.ITriggerPlugin][]pluginbase.IActionPlugin)

	LoadConfig()
}

func main() {
	<-ch
	log.Print("exist")
}

func LoadConfig() {

	currentePath, _ := os.Getwd()
	configPath := currentePath + "/config/"
	files, err := ioutil.ReadDir(configPath)

	if err != nil {
		log.Print(err)
		return
	}

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

		if config.Trigger.Mode == "" {
			continue
		}

		triggerFactory, existed := triggerFactoryMap[config.Trigger.Mode]

		if !existed {
			continue
		}

		triggerPlugin := triggerFactory.CreateTrigger(config.Trigger.Setting)

		if triggerPlugin == nil {
			continue
		}

		triggerPlugin.PolicyHandleFunc(policyHandleFunc)

		for _, filter := range config.Filter {
			filterFactory, existed := filterFactoryMap[filter.Mode]

			if !existed {
				continue
			}

			filterPlugin := filterFactory.CreateFilter(filter.Setting)

			if filterPlugin == nil {
				continue
			}

			pipelineFilterMap[triggerPlugin] = append(pipelineFilterMap[triggerPlugin], filterPlugin)
		}

		for _, action := range config.Action {
			actionFactory, existed := actionFactoryMap[action.Mode]

			if !existed {
				continue
			}

			actionPlugin := actionFactory.CreateAction(action.Setting)

			if actionPlugin == nil {
				continue
			}

			pipelineActionMap[triggerPlugin] = append(pipelineActionMap[triggerPlugin], actionPlugin)
		}
	}

	for triggerPlugin := range pipelineActionMap {
		go triggerPlugin.Start()
	}
}

func policyHandleFunc(triggerPlugin *pluginbase.ITriggerPlugin, throttlingID string, messageFromTrigger *string) {

	parameters := make(map[string]interface{})

	go func() {
		for _, filterPlugin := range pipelineFilterMap[*triggerPlugin] {
			if doNextPipeline := filterPlugin.DoFilter(messageFromTrigger, &parameters); !doNextPipeline {
				return
			}
		}

		log.Printf("fire action...")

		for _, actionPlugin := range pipelineActionMap[*triggerPlugin] {
			actionPlugin.FireAction(messageFromTrigger, &parameters)
		}
	}()
}
