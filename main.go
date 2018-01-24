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
	Trigger struct {
		Mode    string      `yaml:mode`
		Setting interface{} `yaml:setting`
	} `yaml:trigger`

	Action []struct {
		Mode             string           `yaml:mode`
		Setting          interface{}      `yaml:setting`
		ThrottlingPolicy ThrottlingPolicy `yaml:throttlingpolicy`
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
var pipelineMap map[pluginbase.ITriggerPlugin][]pluginbase.IActionPlugin
var actionPolicyMap map[pluginbase.IActionPlugin]ThrottlingPolicy
var policyInstanceMap map[string]pluginbase.IThrottlingPolicyPlugin

func init() {
	loader = PluginImportLoader{}
	//loader = PluginSharedObjectLoader{}

	triggerFactoryMap, actionFactoryMap = loader.Load()
	pipelineMap = make(map[pluginbase.ITriggerPlugin][]pluginbase.IActionPlugin)
	actionPolicyMap = make(map[pluginbase.IActionPlugin]ThrottlingPolicy)
	policyInstanceMap = make(map[string]pluginbase.IThrottlingPolicyPlugin)

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

		if triggerPlugin != nil {
			triggerPlugin.PolicyHandleFunc(policyHandleFunc)
		}

		for _, action := range config.Action {
			actionFactory, existed := actionFactoryMap[action.Mode]

			if !existed {
				continue
			}

			actionPlugin := actionFactory.CreateAction(action.Setting)

			if actionPlugin != nil {
				pipelineMap[triggerPlugin] = append(pipelineMap[triggerPlugin], actionPlugin)

				if action.ThrottlingPolicy.Mode != "" {
					actionPolicyMap[actionPlugin] = action.ThrottlingPolicy
				}
			}
		}
	}

	for triggerPlugin := range pipelineMap {
		go triggerPlugin.Start()
	}
}

func policyHandleFunc(triggerPlugin *pluginbase.ITriggerPlugin, throttlingID string, messageFromTrigger *string) {
	for _, actionPlugin := range pipelineMap[*triggerPlugin] {

		canFireAction := true
		policyInstance, existed := policyInstanceMap[throttlingID]

		if !existed {
			if policyConfig, existed := actionPolicyMap[actionPlugin]; existed {
				policyInstance = throttlingtool.CreatePolicy(policyConfig.Mode, policyConfig.Setting)

				if policyInstance != nil {
					policyInstanceMap[throttlingID] = policyInstance
				}
			}
		}

		if policyInstance != nil {
			canFireAction = policyInstance.Throttling(throttlingID)
		}

		if canFireAction {
			log.Print("action fired")
			actionPlugin.FireAction(throttlingID, messageFromTrigger)
		}
	}
}
