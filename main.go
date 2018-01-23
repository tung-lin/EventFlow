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
		Mode             string      `yaml:mode`
		Setting          interface{} `yaml:setting`
		ThrottlingPolicy struct {
			Mode    string      `yaml:mode`
			Setting interface{} `yaml:setting`
		} `yaml:throttlingpolicy`
	} `yaml:action`
}

var loader IPluginLoader
var triggerFactoryMap map[string]pluginbase.ITriggerFactory
var actionFactoryMap map[string]pluginbase.IActionFactory
var pipelineMap map[pluginbase.ITriggerPlugin][]pluginbase.IActionPlugin
var actionPolicyMap map[pluginbase.IActionPlugin]pluginbase.IThrottlingPolicyPlugin
var ch chan bool

func main() {

	loader = PluginImportLoader{}
	//loader = PluginSharedObjectLoader{}

	triggerFactoryMap, actionFactoryMap = loader.Load()
	pipelineMap = make(map[pluginbase.ITriggerPlugin][]pluginbase.IActionPlugin)
	actionPolicyMap = make(map[pluginbase.IActionPlugin]pluginbase.IThrottlingPolicyPlugin)

	LoadConfig()

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

		if len(config.Trigger) == 0 {
			continue
		}

		var triggerPluginList []*pluginbase.ITriggerPlugin

		for _, trigger := range config.Trigger {
			triggerFactory, existed := triggerFactoryMap[trigger.Mode]

			if existed {
				triggerPlugin := triggerFactory.CreateTrigger(trigger.Setting)

				if triggerPlugin != nil {
					triggerPlugin.PolicyHandleFunc(policyHandleFunc)
					triggerPluginList = append(triggerPluginList, &triggerPlugin)
				}
			}
		}

		if len(triggerPluginList) == 0 {
			continue
		}

		for _, action := range config.Action {
			actionFactory, existed := actionFactoryMap[action.Mode]

			if existed {

				for _, triggerPlugin := range triggerPluginList {

					actionPlugin := actionFactory.CreateAction(action.Setting)

					if actionPlugin != nil {
						pipelineMap[*triggerPlugin] = append(pipelineMap[*triggerPlugin], actionPlugin)

						if action.ThrottlingPolicy.Mode != "" {
							policyPlugin := throttlingtool.CreatePolicy(action.ThrottlingPolicy.Mode, action.ThrottlingPolicy.Setting)
							actionPolicyMap[actionPlugin] = policyPlugin
							//log.Printf("%s-%s", &actionPlugin, &policyPlugin)
						}
					}
				}
			}
		}
	}

	for triggerPlugin := range pipelineMap {
		go triggerPlugin.Start()
	}
}

func policyHandleFunc(triggerPlugin *pluginbase.ITriggerPlugin, throttlingId string, messageFromTrigger *string) {
	for _, actionPlugin := range pipelineMap[*triggerPlugin] {

		canFireAction := true

		if policyPlugin, existed := actionPolicyMap[actionPlugin]; existed {
			//log.Printf("%s-%s", &actionPlugin, &policyPlugin)
			canFireAction = policyPlugin.Throttling(throttlingId)
		}

		if canFireAction {
			log.Print("action fired")
			actionPlugin.FireAction(throttlingId, messageFromTrigger)
		}
	}
}
