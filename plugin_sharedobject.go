package main

import (
	"EventFlow/common/interface/pluginbase"
	"log"
	"os"
	"path/filepath"
	"plugin"
)

type pluginSharedObjectLoader struct {
	triggerFactoryMap map[string]pluginbase.ITriggerFactory
	actionFactoryMap  map[string]pluginbase.IActionFactory
}

func (loader pluginSharedObjectLoader) Load() (triggerFactories map[string]pluginbase.ITriggerFactory, actionFactories map[string]pluginbase.IActionFactory) {

	triggerFactoryMap = make(map[string]pluginbase.ITriggerFactory)
	actionFactoryMap = make(map[string]pluginbase.IActionFactory)

	loader.loadPlugin()

	return triggerFactoryMap, actionFactoryMap
}

func (loader pluginSharedObjectLoader) loadPlugin() {

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
