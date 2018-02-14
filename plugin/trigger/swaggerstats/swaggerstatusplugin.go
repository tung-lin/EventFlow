package swaggerstats

import (
	"EventFlow/common/interface/pluginbase"
	"EventFlow/common/tool/arraytool"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
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

func init() {
	resultChannel = make(chan monitorResult, 20)
	createMonitorResultWriterRoutine()
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

	for _, operation := range trigger.Setting.Operations {

		if operation.OperationID == "" {
			continue
		}

		if operation.MonitorIntervalSec == 0 || operation.MonitorIntervalSec < 20 {
			operation.MonitorIntervalSec = trigger.Setting.MonitorIntervalSec
		}

		if operation.MonitorTimeoutSec == 0 {
			operation.MonitorTimeoutSec = trigger.Setting.MonitorTimeoutSec
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
			if strings.EqualFold(method, "GET") {
				continue
			}

			//skip when OperationID existed in SkipOperations array or api marked as deprecated
			if existed, _ := arraytool.InArray(methodDetail.OperationID, trigger.Setting.SkipOperations); existed || methodDetail.Deprecated {
				continue
			}

			apiURL := url.URL{
				Host:   trigger.swaggerConfig.Host,
				Path:   trigger.swaggerConfig.BasePath + path,
				Scheme: scheme,
			}

			trigger.createAPIMonitor(method, apiURL.String(), methodDetail.OperationID)
		}
	}
}

func (trigger *SwaggerStatusPlugin) createAPIMonitor(method, apiURL, operationID string) {

	var odataParameters []string
	var odatatop int
	var odataformat string
	var specialOperation operation
	var acceptFormat string

	if operation, existed := trigger.operations[operationID]; existed {
		specialOperation = operation
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

	if &specialOperation != nil {
		trigger.createSpecialAPIMonitorRoutine(method, acceptFormat, apiURL, specialOperation)
	} else {
		trigger.createNormalAPIMonitorRoutine(method, acceptFormat, apiURL)
	}
}

func (trigger *SwaggerStatusPlugin) createNormalAPIMonitorRoutine(method, acceptFormat, apiURL string) {
	go func() {
		time.AfterFunc(time.Second*time.Duration(trigger.Setting.MonitorIntervalSec), func() {

			response, _ := getHTTPResponse(trigger.Setting.MonitorTimeoutSec, acceptFormat, apiURL)

			result := monitorResult{
				method:          method,
				url:             apiURL,
				responseCode:    response.StatusCode,
				responseMessage: getErrorResponseContent(response.Body),
			}

			resultChannel <- result

			if response.StatusCode != 200 {

				msg := fmt.Sprintf("url: %s\r\nresponse code: %d\r\nmessage: %s\r\n", apiURL, response.StatusCode, result.responseMessage)

				var triggerPlugin pluginbase.ITriggerPlugin = trigger
				trigger.FireAction(&triggerPlugin, &msg)
			} else {

			}
		})
	}()
}

func (trigger *SwaggerStatusPlugin) createSpecialAPIMonitorRoutine(method, acceptFormat, apiURL string, operation operation) {
	go func() {

		time.AfterFunc(time.Second*time.Duration(operation.MonitorIntervalSec), func() {

			response, _ := getHTTPResponse(operation.MonitorTimeoutSec, acceptFormat, apiURL)

			result := monitorResult{
				method:          method,
				url:             apiURL,
				responseCode:    response.StatusCode,
				responseMessage: getErrorResponseContent(response.Body),
			}

			resultChannel <- result

		})
	}()
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

func getHTTPResponse(monitorTimeoutSec int, method, acceptFormat, apiURL string) (response *http.Response, err error) {
	var client = &http.Client{
		Timeout: time.Second * time.Duration(monitorTimeoutSec),
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	if acceptFormat != "" {
		request, _ := http.NewRequest("GET", apiURL, nil)
		request.Header.Set("Accept", acceptFormat)
		response, err = client.Do(request)
	} else {
		response, err = client.Get(apiURL)
	}

	return
}

func getErrorResponseContent(body io.ReadCloser) string {
	bodyBytes, _ := ioutil.ReadAll(body)
	bodyString := string(bodyBytes)

	return bodyString
}
