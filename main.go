package main

import (
	"EventFlow/common/interface/pluginbase"
	"EventFlow/common/tool/logtool"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	yaml "gopkg.in/yaml.v2"
)

type eventFlow struct {
	Trigger struct {
		Mode    string      `yaml:"mode"`
		Setting interface{} `yaml:"setting"`
	} `yaml:"trigger"`

	Filter []struct {
		Mode    string      `yaml:"mode"`
		Setting interface{} `yaml:"setting"`
	} `yaml:"filter"`

	Action []struct {
		Mode    string      `yaml:"mode"`
		Setting interface{} `yaml:"setting"`
	} `yaml:"action"`
}

var ch chan bool

var triggerFactoryMap map[string]pluginbase.ITriggerFactory
var actionFactoryMap map[string]pluginbase.IActionFactory
var filterFactoryMap map[string]pluginbase.IFilterFactory

var pipelineFilterMap map[pluginbase.ITriggerPlugin][]pluginbase.IFilterPlugin
var pipelineActionMap map[pluginbase.ITriggerPlugin][]pluginbase.IActionPlugin

func init() {
	loader := pluginImportLoader{}
	//loader := pluginSharedObjectLoader{}

	triggerFactoryMap, filterFactoryMap, actionFactoryMap = loader.Load()

	pipelineFilterMap = make(map[pluginbase.ITriggerPlugin][]pluginbase.IFilterPlugin)
	pipelineActionMap = make(map[pluginbase.ITriggerPlugin][]pluginbase.IActionPlugin)

	loadConfig()
}

func main() {

	<-ch
	log.Print("exist")
}

func loadConfig() {

	currentePath, _ := os.Getwd()
	pipelineConfigPath := currentePath + "/config/pipeline/"
	files, err := ioutil.ReadDir(pipelineConfigPath)

	if err != nil {
		log.Print(err)
		return
	}

	for _, file := range files {
		pipelineFile, err := ioutil.ReadFile(pipelineConfigPath + file.Name())

		if err != nil {
			logtool.Fatal("main", "main", fmt.Sprintf("read pipeline config file failed: %v", err))
			continue
		}

		var config eventFlow
		err = yaml.Unmarshal(pipelineFile, &config)

		if err != nil {
			logtool.Fatal("main", "main", fmt.Sprintf("unmarshal pipeline config file failed: %v", err))
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

		triggerPlugin.ActionHandleFunc(actionHandleFunc)

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

		go triggerPlugin.Start()
	}
}

func actionHandleFunc(triggerPlugin *pluginbase.ITriggerPlugin, messageFromTrigger *string) {

	go func() {
		parameters := make(map[string]interface{})

		//start := time.Now()

		for _, filterPlugin := range pipelineFilterMap[*triggerPlugin] {
			if doNextPipeline := filterPlugin.DoFilter(messageFromTrigger, &parameters); !doNextPipeline {
				return
			}
		}

		//log.Printf("took %s", time.Since(start))

		log.Printf("fire action...")

		for _, actionPlugin := range pipelineActionMap[*triggerPlugin] {
			go actionPlugin.FireAction(messageFromTrigger, &parameters)
		}
	}()
}
