package pipeline

import (
	"EventFlow/common/interface/pluginbase"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"plugin"

	"gopkg.in/yaml.v2"
)

var Config EventFlow

type EventFlow struct {
	Trigger []struct {
		Mode    string      `yaml:mode`
		Setting interface{} `yaml:setting`
	} `yaml:trigger`

	Action []struct {
		Mode    string      `yaml:mode`
		Setting interface{} `yaml:setting`
	} `yaml:action`
}

type Action struct {
}

type Setting struct {
	Swagger_url string      `yaml:swagger_url`
	API_IP      string      `yaml:api_ip`
	API_Path    string      `yaml:api_path`
	Operations  []Operation `yaml:operations`
}

type Operation struct {
	OperationId    string   `yaml:operationid`
	Condition      string   `yaml:condition`
	Fields         []string `yaml:fields`
	Threshold_Type string   `yaml:threshold_type`
	Threshold      string   `yaml:threshold`
}

func init() {

	triggerFactoryMap = make(map[string]pluginbase.ITriggerFactory)
	actionFactoryMap = make(map[string]pluginbase.IActionFactory)

	loadPLugin()
}

var triggerFactoryMap map[string]pluginbase.ITriggerFactory
var actionFactoryMap map[string]pluginbase.IActionFactory

func loadPLugin() {

	currentePath, _ := os.Getwd()
	pluginPath := currentePath + "/plugin/"

	triggerPath := pluginPath + "trigger/"
	actionPath := pluginPath + "action/"

	triggerFiles := getSymbolFiles(triggerPath)
	actionFiles := getSymbolFiles(actionPath)

	for _, triggerFile := range triggerFiles {

		plugin, err := plugin.Open(triggerFile)

		if err != nil {
			log.Printf("load symbol file '%s' failed: %s", triggerFile, err)
			continue
		}

		symbol, err := plugin.Lookup("Factory")

		if err != nil {
			log.Printf("look up a symbol 'Factory' failed: %s", err)
			continue
		}

		triggerFactory, ok := symbol.(pluginbase.ITriggerFactory)

		if ok {
			triggerFactoryMap[triggerFactory.GetIdentifyName()] = triggerFactory
		}
	}

	for _, actionFile := range actionFiles {

		plugin, err := plugin.Open(actionFile)

		if err != nil {
			log.Printf("load symbol file '%s' failed: %s", actionFile, err)
			continue
		}

		symbol, err := plugin.Lookup("Factory")

		if err != nil {
			log.Printf("look up a symbol 'Factory' failed: %s", err)
			continue
		}

		actionFactory, ok := symbol.(pluginbase.IActionFactory)

		if ok {
			actionFactoryMap[actionFactory.GetIdentifyName()] = actionFactory
		}
	}
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

		for _, trigger := range config.Trigger {
			triggerFactory, existed := triggerFactoryMap[trigger.Mode]

			if existed {
				triggerFactory.CreateTrigger(trigger.Setting)
			}
		}

		for _, action := range config.Action {
			actionFactory, existed := actionFactoryMap[action.Mode]

			if existed {
				actionFactory.CreateAction(action.Setting)
			}
		}
	}
}

func structToByteArray(setting interface{}) (bytes []byte) {

	if setting == nil {
		return nil
	}

	str, err := yaml.Marshal(setting)

	if err != nil {
		log.Printf("Marshal setting config to string failed: %v", err)
		return nil
	}

	return []byte(str)
}

func getSymbolFiles(path string) (files []string) {

	triggerFiles := []string{}

	filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if filepath.Ext(path) == ".so" {
			triggerFiles = append(triggerFiles, path)
		}
		return nil
	})

	return triggerFiles
}

func Test() {
	configFilePath, _ := os.Getwd()
	file, err := ioutil.ReadFile(configFilePath + "/config/test.yaml")

	if err != nil {
		log.Fatalf("Read config file failed: %v", err)
	}

	err = yaml.Unmarshal(file, &Config)

	if err != nil {
		log.Fatalf("Unmarshal config file failed: %v", err)
	}

	fmt.Printf("%+v\n", Config.Trigger[0].Setting)

	var setting Setting

	str, err := yaml.Marshal(Config.Trigger[0].Setting)
	fmt.Print(str)
	err = yaml.Unmarshal([]byte(str), &setting)

	fmt.Printf("%+v\n", setting)
}
