package main

import (
	"EventFlow/common/interface/pluginbase"
	"EventFlow/common/tool/logtool"
	"EventFlow/common/tool/pipelinetool"
	"EventFlow/common/tool/stringtool"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"

	yaml "gopkg.in/yaml.v2"
)

type eventFlow struct {
	Trigger   []triggerConfig        `yaml:"trigger"`
	Condition pipelinetool.Condition `yaml:"condition"`
	Filter    []filterConfig         `yaml:"filter"`
	Action    []actionConfig         `yaml:"action"`
}

type triggerConfig struct {
	Mode    string      `yaml:"mode"`
	Disable bool        `yaml:"disable"`
	Setting interface{} `yaml:"setting"`
}

type filterConfig struct {
	Mode      string                 `yaml:"mode"`
	Disable   bool                   `yaml:"disable"`
	Setting   interface{}            `yaml:"setting"`
	Condition pipelinetool.Condition `yaml:"condition"`
}

type actionConfig struct {
	Mode      string                 `yaml:"mode"`
	Disable   bool                   `yaml:"disable"`
	Setting   interface{}            `yaml:"setting"`
	Condition pipelinetool.Condition `yaml:"condition"`
}

type pipeline struct {
	PipelineID       string
	PipelineFileName string
	Trigger          []pluginbase.ITriggerPlugin
	Condition        pluginbase.IConditionPlugin
	Filter           []pluginbase.IFilterPlugin
	FilterConfig     []filterConfig
	Action           []pluginbase.IActionPlugin
	ActionConfig     []actionConfig
}

var ch chan os.Signal
var runForever bool

var triggerFactoryMap map[string]pluginbase.ITriggerFactory
var conditionFactory pluginbase.IConditionFactory
var actionFactoryMap map[string]pluginbase.IActionFactory
var filterFactoryMap map[string]pluginbase.IFilterFactory

var pipelineMap map[string]*pipeline

func init() {
	loader := pluginImportLoader{}
	//loader := pluginSharedObjectLoader{}

	triggerFactoryMap, conditionFactory, filterFactoryMap, actionFactoryMap = loader.Load()

	pipelineMap = make(map[string]*pipeline)

	runForever = false
	loadConfig()
}

func main() {

	exitFunc := func() {
		logtool.Debug("main", "main", "exit program")
		os.Exit(0)
	}

	if runForever {
		ch = make(chan os.Signal)
		signal.Notify(ch, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

		go func() {
			signal := <-ch
			logtool.Debug("main", "main", fmt.Sprintf("caught signal: %+v", signal))

			for _, pipeline := range pipelineMap {
				for _, triggerPlugin := range pipeline.Trigger {
					triggerPlugin.Stop()
				}
			}

			exitFunc()
		}()

		select {}

	} else {
		exitFunc()
	}
}

func loadConfig() {

	//load pipeline file folder
	currentePath, _ := os.Getwd()
	pipelineConfigPath := currentePath + "/config/pipeline/"
	files, err := ioutil.ReadDir(pipelineConfigPath)

	if err != nil {
		logtool.Error("main", "main", fmt.Sprintf("read pipeline config directory failed: %v", err))
		return
	}

	for _, file := range files {

		if file.IsDir() {
			continue
		}

		filename := file.Name()

		//load pieline file
		pipelineFile, err := ioutil.ReadFile(pipelineConfigPath + file.Name())

		if err != nil {
			logtool.Error("main", "main", fmt.Sprintf("read pipeline config file '%s 'failed: %v", filename, err))
			continue
		}

		//unmarshal yaml file to struct
		var config eventFlow
		err = yaml.Unmarshal(pipelineFile, &config)

		if err != nil {
			logtool.Fatal("main", "main", fmt.Sprintf("unmarshal pipeline config file '%s' failed: %v", filename, err))
			continue
		}

		//create pipeline struct
		pipelineID := stringtool.CreateRandomString()

		for _, existed := pipelineMap[pipelineID]; existed; {
			pipelineID = stringtool.CreateRandomString()
		}

		pipeline := pipeline{
			PipelineID:       pipelineID,
			PipelineFileName: file.Name(),
			Trigger:          []pluginbase.ITriggerPlugin{},
			Filter:           []pluginbase.IFilterPlugin{},
			FilterConfig:     []filterConfig{},
			Action:           []pluginbase.IActionPlugin{},
			ActionConfig:     []actionConfig{},
		}

		pipelineMap[pipelineID] = &pipeline

		logtool.Debug("main", "main", fmt.Sprintf("read pipeline config file '%s' (%d trigger(s), %d filter(s), %d action(s))", filename, len(config.Trigger), len(config.Filter), len(config.Action)))

		//create condition
		if conditionFactory != nil {
			pipeline.Condition = conditionFactory.CreateCondition(config.Condition)
		}

		//create filters
		for _, filter := range config.Filter {

			if filter.Disable {
				logtool.Info("main", "main", fmt.Sprintf("filter mode '%s' in file '%s' is disabled", filter.Mode, filename))
				continue
			}

			filterFactory, existed := filterFactoryMap[filter.Mode]

			if !existed {
				continue
			}

			filterPlugin := filterFactory.CreateFilter(filter.Setting)

			if filterPlugin == nil {
				continue
			}

			pipeline.Filter = append(pipeline.Filter, filterPlugin)
			pipeline.FilterConfig = append(pipeline.FilterConfig, filter)
		}

		//create actions
		for _, action := range config.Action {

			if action.Disable {
				logtool.Info("main", "main", fmt.Sprintf("action mode '%s' in file '%s' is disabled", action.Mode, filename))
				continue
			}

			actionFactory, existed := actionFactoryMap[action.Mode]

			if !existed {
				continue
			}

			actionPlugin := actionFactory.CreateAction(action.Setting)

			if actionPlugin == nil {
				continue
			}

			pipeline.Action = append(pipeline.Action, actionPlugin)
			pipeline.ActionConfig = append(pipeline.ActionConfig, action)
		}

		//create and execute triggers
		for _, trigger := range config.Trigger {

			if trigger.Mode == "" {
				logtool.Info("main", "main", fmt.Sprintf("trigger mode is undefined in file '%s'", filename))
				continue
			}

			if trigger.Disable {
				logtool.Info("main", "main", fmt.Sprintf("trigger mode '%s' in file '%s' is disabled", trigger.Mode, filename))
				continue
			}

			triggerFactory, existed := triggerFactoryMap[trigger.Mode]

			if !existed {
				logtool.Info("main", "main", fmt.Sprintf("trigger mode '%s' not found", trigger.Mode))
				continue
			}

			triggerPlugin := triggerFactory.CreateTrigger(trigger.Setting)

			if triggerPlugin == nil {
				continue
			}

			pipeline.Trigger = append(pipeline.Trigger, triggerPlugin)

			triggerPlugin.ActionHandleFunc(pipeline.PipelineID, actionHandleFunc)

			go triggerPlugin.Start()
		}
	}

	runForever = true
}

func actionHandleFunc(pipelineID string, triggerPlugin *pluginbase.ITriggerPlugin, messageFromTrigger *string) {

	go func() {

		//find pipeline by ID
		if pipeline, existed := pipelineMap[pipelineID]; existed {

			doAction := true
			parameters := make(map[string]interface{})

			//execute condition
			if pipeline.Condition != nil && !pipeline.Condition.IsMatch(&parameters) {
				return
			}

			//execute filters
			for filterConfigIndex, filterPlugin := range pipeline.Filter {

				filterConfig := pipeline.FilterConfig[filterConfigIndex]

				if len(filterConfig.Condition) > 0 && !pipelinetool.IsMatchCondition(&filterConfig.Condition, &parameters) {
					continue
				}

				if doNextPipeline := filterPlugin.DoFilter(messageFromTrigger, &parameters); !doNextPipeline {
					doAction = false
					break
				}
			}

			//execute action
			if doAction {
				for actionConfigIndex, actionPlugin := range pipeline.Action {

					actionConfig := pipeline.ActionConfig[actionConfigIndex]

					if len(actionConfig.Condition) > 0 && !pipelinetool.IsMatchCondition(&actionConfig.Condition, &parameters) {
						continue
					}

					go actionPlugin.FireAction(messageFromTrigger, &parameters)
				}
			}
		}
	}()
}
