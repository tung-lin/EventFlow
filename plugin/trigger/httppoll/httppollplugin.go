package httppoll

import (
	"EventFlow/common/interface/pluginbase"
	"EventFlow/common/tool/arraytool"
	"EventFlow/common/tool/jsontool"
	"EventFlow/common/tool/logtool"
	"EventFlow/plugin/trigger/httppoll/authbase"
	"EventFlow/plugin/trigger/httppoll/authmode/basic"
	"EventFlow/plugin/trigger/httppoll/authmode/hmac"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/robfig/cron"
)

type pluginOutput struct {
	Url            string
	TimeElapsedMS  int64
	ResponseLength int64
	ResponseCode   int
	ResponseBody   interface{}
}

var authFactoryMap = make(map[string]authbase.IAuthFactory)
var allowHMacMode = []string{"PTX"}
var allowCodec = []string{"json"}

type HttpPollPlugin struct {
	pluginbase.ActionHandler
	authPlugin authbase.IAuthPlugin
	Setting    SettingConfig
	cron       *cron.Cron
	httpClient *http.Client
}

func init() {

	authList := []authbase.IAuthFactory{}
	authList = append(authList, basic.BasicFactory{}, hmac.HMacFactory{})

	for _, auth := range authList {
		authFactoryMap[auth.GetIdentifyName()] = auth
	}
}

func (trigger *HttpPollPlugin) Start() {

	if len(trigger.Setting.Urls) == 0 {
		logtool.Warn("trigger", "http_poll", "http polling is disabled due to empty url setting")
		return
	}

	if trigger.Setting.Cron == "" {
		logtool.Warn("trigger", "http_poll", "http polling is disabled due to empty cron setting")
		return
	}

	if existed, _ := arraytool.InArray(trigger.Setting.Codec, allowCodec); !existed {
		logtool.Error("trigger", "http_poll", fmt.Sprintf("codec '%s' is not supported", trigger.Setting.Codec))
		return
	}

	for index := range trigger.Setting.Urls {

		urlItem := &trigger.Setting.Urls[index]

		if urlItem.Codec == "" {
			urlItem.Codec = trigger.Setting.Codec
		} else {
			if existed, _ := arraytool.InArray(urlItem.Codec, allowCodec); !existed {
				logtool.Error("trigger", "http_poll", fmt.Sprintf("codec '%s' is not supported", urlItem.Codec))
				return
			}
		}
	}

	if authFactory, existed := authFactoryMap[trigger.Setting.Auth.Mode]; existed {
		trigger.authPlugin = authFactory.CreateAuth(trigger.Setting.Auth.Setting)

		err := trigger.authPlugin.CheckParameter()

		if err != nil {
			logtool.Error("trigger", "http_poll", err.Error())
			return
		}

		trigger.httpClient = &http.Client{
			Timeout: time.Duration(trigger.Setting.TimeoutMS) * time.Millisecond,
		}

		trigger.cron = cron.New()
		err = trigger.cron.AddFunc(trigger.Setting.Cron, trigger.cronjobInvokeHandler)

		if err != nil {
			logtool.Error("trigger", "http_poll", fmt.Sprintf("create cron(%v) failed: %v", trigger.Setting.Cron, err.Error()))
		} else {
			trigger.cron.Start()
		}

	} else {
		logtool.Error("trigger", "http_poll", fmt.Sprintf("auth mode '%s' not found", trigger.Setting.Auth.Mode))
	}
}

func (trigger *HttpPollPlugin) cronjobInvokeHandler() {
	headers := trigger.authPlugin.CreateHttpHeaders()

	var httpHeaders = map[string][]string{}

	for _, header := range headers {
		if values, existed := httpHeaders[header.Key]; existed {
			values = append(values, header.Value)
		} else {
			httpHeaders[header.Key] = []string{header.Value}
		}
	}

	for _, urlItem := range trigger.Setting.Urls {

		request, _ := http.NewRequest("GET", urlItem.Url, nil)
		request.Header = httpHeaders

		start := time.Now()
		response, err := trigger.httpClient.Do(request)
		elapsedMS := time.Since(start).Nanoseconds() / 1000000

		pluginOutput := pluginOutput{
			Url:            urlItem.Url,
			TimeElapsedMS:  elapsedMS,
			ResponseCode:   response.StatusCode,
			ResponseLength: response.ContentLength,
		}

		if err != nil {
			logtool.Error("trigger", "http_poll", fmt.Sprintf("an error occurred while calling http request: %v", err.Error()))
		} else {
			defer response.Body.Close()

			//if header, existed := response.Header["Content-Type"]; existed {
			//	if len(header) > 0 {

			//	}
			//}

			body, err := ioutil.ReadAll(response.Body)

			if err != nil {
				logtool.Error("trigger", "http_poll", fmt.Sprintf("an error occurred while reading http response body: %v", err.Error()))
			} else {
				pluginOutput.ResponseBody = string(body)

			}

			var content string
			jsontool.MarshalToString(pluginOutput, &content)

			var triggerPlugin pluginbase.ITriggerPlugin = trigger
			trigger.FireAction(&triggerPlugin, &content)
		}
	}
}

func (trigger *HttpPollPlugin) Stop() {
	if trigger.cron != nil {
		trigger.cron.Stop()
	}
}
