package json

import (
	"EventFlow/common/tool/arraytool"
	"EventFlow/common/tool/logtool"
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

type JSONPlugin struct {
	Setting SettingConfig
}

func (filter *JSONPlugin) DoFilter(messageFromTrigger *string, parameters *map[string]interface{}) (doNextPipeline bool) {

	jsonContent, err := isJSONFormat(*messageFromTrigger)

	if err != nil {
		logtool.Error("filter", "JSON", fmt.Sprintf("content is not a valid json format data: %v", err))
		for metadataKey := range filter.Setting.AddMetadata {
			(*parameters)[metadataKey] = ""
		}
	} else {
		for metadataKey, metadataParm := range filter.Setting.AddMetadata {
			fieldValue := getField(metadataParm, jsonContent)
			(*parameters)[metadataKey] = fieldValue
		}
	}

	return true
}

func isJSONFormat(content string) (jsonContent map[string]*json.RawMessage, err error) {

	data := []byte(content)
	err = json.Unmarshal(data, &jsonContent)

	return jsonContent, err
}

func getField(fieldPath string, jsonContent map[string]*json.RawMessage) (fieldValue interface{}) {

	paths := strings.Split(fieldPath, "]")
	paths = arraytool.RemoveItem("", paths)

	if len(paths) == 0 {
		return ""
	}

	objectMap := jsonContent

	for index, path := range paths {

		if path == "" {
			continue
		}

		path = strings.Replace(path, "[", "", -1)

		value, existed := objectMap[path]

		if !existed {
			return ""
		}

		//is the last item in array
		if index == len(paths)-1 {

			var objectValue interface{}
			err := json.Unmarshal(*value, &objectValue)

			if err != nil {
				log.Printf("%v", err)
				break
			}

			return objectValue
		}

		err := json.Unmarshal(*value, &objectMap)

		if err != nil {
			log.Printf("%v", err)
			break
		}
	}

	return ""
}
