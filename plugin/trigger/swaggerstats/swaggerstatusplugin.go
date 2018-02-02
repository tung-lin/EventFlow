package swaggerstats

import (
	"EventFlow/common/interface/pluginbase"
	"EventFlow/common/tool/arraytool"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type SwaggerConfig struct {
	Host     string   `json:"host"`
	BasePath string   `json:"basePath"`
	Schemes  []string `json:"schemes"`
	Paths    map[string]map[string]struct {
		Tags        []string `json:"tags"`
		Summary     string   `json:"summary"`
		Description string   `json:"description"`
		OperationID string   `json:"operationId"`
		Deprecated  bool     `json:"deprecated"`
	} `json:"paths"`
}

type SwaggerStatusPlugin struct {
	pluginbase.ActionHandler
	Setting       SettingConfig
	swaggerConfig SwaggerConfig
}

func (trigger *SwaggerStatusPlugin) Start() {

	if trigger.Setting.SwaggerURL == "" {
		return
	}

	go func() {
		trigger.swaggerParsing()

	}()
}

func (trigger *SwaggerStatusPlugin) Stop() {

}

func (trigger *SwaggerStatusPlugin) swaggerParsing() {

	response, err := http.Get(trigger.Setting.SwaggerURL)

	if err != nil {
		log.Printf("[trigger][swagger] Download swagger json failed: %v", err)
		return
	}

	defer response.Body.Close()

	swaggerJSONContents, err := ioutil.ReadAll(response.Body)

	if err != nil {
		fmt.Printf("[trigger][swagger] Read http response body failed: %v", err)
		return
	}

	err = json.Unmarshal(swaggerJSONContents, &trigger.swaggerConfig)

	if err != nil {
		log.Printf("[trigger][swagger] Unmarshal swagger json failed: %v", err)
		return
	}
}

func (trigger *SwaggerStatusPlugin) createAPIMonitor() {

	var scheme string

	if len(trigger.swaggerConfig.Schemes) == 0 {
		scheme = "http"
	} else {
		scheme = trigger.swaggerConfig.Schemes[0]
	}

	for path, pathDetail := range trigger.swaggerConfig.Paths {
		for method, methodDetail := range pathDetail {

			//skip http method other than http GET
			if strings.EqualFold(method, "GET") {
				continue
			}

			if existed, _ := arraytool.InArray(methodDetail.OperationID, trigger.Setting.SkipOperations); existed || methodDetail.Deprecated {
				continue
			}

			apiURL := url.URL{
				Host:   trigger.swaggerConfig.Host,
				Path:   trigger.swaggerConfig.BasePath + path,
				Scheme: scheme,
			}

			go trigger.startAPIMonitor(method, apiURL.String())
		}
	}
}

func (trigger *SwaggerStatusPlugin) startAPIMonitor(method string, apiURL string) {
	//request, err := http.NewRequest(method, apiURL, nil)

	//if err != nil {
	//	log.Printf("[trigger][swagger]%v", err)
	//	return
	//}

}
