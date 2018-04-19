package swaggerstats

import (
	"EventFlow/common/interface/pluginbase"
	"EventFlow/common/tool/arraytool"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type swaggerConfig struct {
	Host     string   `json:"host"`
	BasePath string   `json:"basePath"`
	Schemes  []string `json:"schemes"`
	Paths    map[string]map[string]struct {
		Tags        []string `json:"tags"`
		Summary     string   `json:"summary"`
		Description string   `json:"description"`
		OperationID string   `json:"operationId"`
		Deprecated  bool     `json:"deprecated"`
		Parameters  []struct {
			Name        string `json:"name"`
			In          string `json:"in"`
			Description string `json:"description"`
			Required    bool   `json:"required"`
			Type        string `json:"type"`
			Enum        []struct {
				Text  string
				Value interface{}
			} `json:"enum"`
		} `json:"parameters"`
	} `json:"paths"`
}

type SwaggerStatusPlugin struct {
	pluginbase.ActionHandler
	Setting       SettingConfig
	swaggerConfig swaggerConfig
	operations    map[string]operation
}

type monitorResult struct {
	method          string
	url             string
	responseCode    int
	responseMessage string
}

var resultChannel chan monitorResult

func (trigger *SwaggerStatusPlugin) Start() {

	if trigger.Setting.SwaggerURL == "" {
		return
	}

	resultChannel = make(chan monitorResult, 20)
	createMonitorResultWriterRoutine()

	trigger.operations = make(map[string]operation)

	go func() {
		trigger.swaggerParsing()
	}()
}

func (trigger *SwaggerStatusPlugin) Stop() {

}

func (trigger *SwaggerStatusPlugin) swaggerParsing() {

	currentePath, _ := os.Getwd()
	swaggerFilePath := fmt.Sprintf("%s/%s", currentePath, trigger.Setting.SwaggerFile)

	var swaggerJSONContents []byte

	if _, err := os.Stat(swaggerFilePath); err == nil {
		jsonContent, err := ioutil.ReadFile(swaggerFilePath)

		if err != nil {
			fmt.Printf("[trigger][swagger] Read swagger file failed: %v", err)
			return
		}

		swaggerJSONContents = jsonContent

	} else {
		response, err := http.Get(trigger.Setting.SwaggerURL)

		if err != nil {
			log.Printf("[trigger][swagger] Download swagger json failed: %v", err)
			return
		}

		defer response.Body.Close()

		jsonContent, err := ioutil.ReadAll(response.Body)

		if err != nil {
			fmt.Printf("[trigger][swagger] Read http response body failed: %v", err)
			return
		}

		swaggerJSONContents = jsonContent
	}

	err := json.Unmarshal(swaggerJSONContents, &trigger.swaggerConfig)

	if err != nil {
		log.Printf("[trigger][swagger] Unmarshal swagger json failed: %v", err)
		return
	}

	for _, operation := range trigger.Setting.Operations {

		if operation.OperationID == "" {
			continue
		}

		if operation.MonitorIntervalSec < 20 {
			operation.MonitorIntervalSec = 20
		}

		if operation.MonitorTimeoutSec <= 0 {
			operation.MonitorTimeoutSec = 10
		}

		if operation.MonitorTimeoutSec > operation.MonitorIntervalSec {
			operation.MonitorTimeoutSec = operation.MonitorIntervalSec - 1
		}

		if operation.ODataTop > 20 {
			operation.ODataTop = trigger.Setting.ODataTop
		}

		if operation.ODataFormat != "" && !strings.EqualFold(operation.ODataFormat, "json") && !strings.EqualFold(operation.ODataFormat, "xml") {
			operation.ODataFormat = "json"
		}

		trigger.operations[operation.OperationID] = operation
	}

	trigger.apiParsing()
}

func (trigger *SwaggerStatusPlugin) apiParsing() {

	var scheme string

	if len(trigger.swaggerConfig.Schemes) == 0 {
		scheme = "http"
	} else {
		scheme = trigger.swaggerConfig.Schemes[0]
	}

	for path, pathDetail := range trigger.swaggerConfig.Paths {
		for method, methodDetail := range pathDetail {

			//skip http method other than http GET
			if method == "GET" {
				continue
			}

			//skip when OperationID existed in SkipOperations array or api marked as deprecated
			if existed, _ := arraytool.InArray(methodDetail.OperationID, trigger.Setting.SkipOperations); existed || methodDetail.Deprecated {
				continue
			}

			hasParameterInPath := false
			supportOData_Top := false
			supportOData_Format := false

			for _, parameter := range methodDetail.Parameters {
				if parameter.In == "path" {
					hasParameterInPath = true
				} else if parameter.In == "query" {
					if parameter.Name == "$top" {
						supportOData_Top = true
					} else if parameter.Name == "$format" {
						supportOData_Format = true
					}
				}
			}

			operation, existed := trigger.operations[methodDetail.OperationID]

			if hasParameterInPath {
				if !existed {
					continue
				}

				for _, parameter := range operation.Parameters {

					pathReplacedWithParameter := path

					for _, pair := range parameter.Values {
						pathReplacedWithParameter = strings.Replace(pathReplacedWithParameter, fmt.Sprintf("{%s}", pair.Name), pair.value.(string), -1)
					}

					apiURL := url.URL{
						Host:   trigger.swaggerConfig.Host,
						Path:   trigger.swaggerConfig.BasePath + pathReplacedWithParameter,
						Scheme: scheme,
					}

					trigger.createAPIMonitor(apiURL.String(), operation)
				}
			} else {
				apiURL := url.URL{
					Host:   trigger.swaggerConfig.Host,
					Path:   trigger.swaggerConfig.BasePath + path,
					Scheme: scheme,
				}

				trigger.createAPIMonitor(apiURL.String(), operation)
			}
		}
	}
}

func (trigger *SwaggerStatusPlugin) createAPIMonitor(apiURL string, operation operation) {

	var odataParameters []string
	var odatatop int
	var odataformat string
	var acceptFormat string

	if &operation != nil {
		odatatop = operation.ODataTop
		odataformat = operation.ODataFormat
	} else {
		odatatop = trigger.Setting.ODataTop
		odataformat = trigger.Setting.ODataFormat
	}

	if odatatop != 0 {
		odataParameters = append(odataParameters, fmt.Sprintf("$top=%d", odatatop))
	}

	if odataformat != "" {
		odataParameters = append(odataParameters, fmt.Sprintf("$format=%s", odataformat))
	} else {
		acceptFormat = "application/json"
	}

	if len(odataParameters) != 0 {
		apiURL = fmt.Sprintf("%s?%s", apiURL, strings.Join(odataParameters, "&"))
	}

	trigger.createAPIMonitorTimer(acceptFormat, apiURL, operation)
}

func (trigger *SwaggerStatusPlugin) createAPIMonitorTimer(acceptFormat, apiURL string, operation operation) {

	var interval int
	var timeout int

	if &operation != nil {
		interval = operation.MonitorIntervalSec
		timeout = operation.MonitorTimeoutSec
	} else {
		interval = trigger.Setting.MonitorIntervalSec
		timeout = trigger.Setting.MonitorTimeoutSec
	}

	time.AfterFunc(time.Second*time.Duration(interval), func() {

		result := getHTTPResponse(timeout, acceptFormat, apiURL)

		msg := fmt.Sprintf("url: %s\r\nresponse code: %d\r\nmessage: %s\r\n", apiURL, result.responseCode, result.responseMessage)

		if result.responseCode != 200 || (result.responseCode == 200 && trigger.Setting.NotifyHTTPOK) {
			var triggerPlugin pluginbase.ITriggerPlugin = trigger
			trigger.FireAction(&triggerPlugin, &msg)
		}
	})
}

func getHTTPResponse(monitorTimeoutSec int, acceptFormat, apiURL string) (result monitorResult) {
	var client = &http.Client{
		Timeout: time.Second * time.Duration(monitorTimeoutSec),
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	var response *http.Response

	if acceptFormat != "" {
		request, _ := http.NewRequest("GET", apiURL, nil)
		request.Header.Set("Accept", acceptFormat)
		response, _ = client.Do(request)
	} else {
		response, _ = client.Get(apiURL)
	}

	bodyBytes, _ := ioutil.ReadAll(response.Body)
	bodyMessage := string(bodyBytes)

	result = monitorResult{
		method:          "GET",
		url:             apiURL,
		responseCode:    response.StatusCode,
		responseMessage: bodyMessage,
	}

	resultChannel <- result

	return
}

func createMonitorResultWriterRoutine() {
	go func() {
		for {
			result := <-resultChannel

			if result.method != "" {

			}
		}
	}()
}
