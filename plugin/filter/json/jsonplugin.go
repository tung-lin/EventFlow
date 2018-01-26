package json

import (
	"EventFlow/common/tool/arraytool"
	"encoding/json"
	"log"
	"strings"
)

type JSONPlugin struct {
	Setting SettingConfig
}

func (filter *JSONPlugin) DoFilter(messageFromTrigger *string, parameters *map[string]interface{}) (doNextPipeline bool) {

	err, jsonContent := isJSONFormat(*messageFromTrigger)

	if err != nil {
		log.Printf("[Filter][JSON] Content is not a valid json format: %v", err)
		return false
	}

	for key, fieldPath := range filter.Setting.AddMetadata {
		fieldValue := getField(fieldPath, jsonContent)

		if fieldValue != "" {
			(*parameters)[key] = fieldValue
		}
	}

	return true
}

func isJSONFormat(content string) (err error, jsonContent map[string]*json.RawMessage) {

	data := []byte(content)
	err = json.Unmarshal(data, &jsonContent)

	return err, jsonContent
}

func getField(fieldPath string, jsonContent map[string]*json.RawMessage) (fieldValue string) {

	paths := strings.Split(fieldPath, "]")
	paths = arraytool.RemoveItem(paths, "")

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
			var stringValue string
			err := json.Unmarshal(*value, &stringValue)

			if err != nil {
				log.Printf("%v", err)
				break
			}

			return stringValue
		} else {
			err := json.Unmarshal(*value, &objectMap)

			if err != nil {
				log.Printf("%v", err)
				break
			}
		}
	}

	return ""
}
