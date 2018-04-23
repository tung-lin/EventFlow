package http

import (
	"EventFlow/common/interface/pluginbase"
	"EventFlow/common/tool/arraytool"
	"EventFlow/common/tool/logtool"
	"context"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	nethttp "net/http"

	"github.com/gorilla/mux"
)

type HttpPlugin struct {
	currentServer *nethttp.Server
	pluginbase.ActionHandler
	Setting SettingConfig
}

var allowHttpMethods []string

func (trigger *HttpPlugin) Start() {

	allowHttpMethods = []string{"GET", "POST"}
	handler := trigger.handleRequestFunc

	router := mux.NewRouter()
	router.HandleFunc("/", handler)
	//router.HandleFunc("/trigger/{triggerid:.*}", handler)

	addr := fmt.Sprintf(":%d", trigger.Setting.Port)

	trigger.currentServer = &nethttp.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	trigger.currentServer.RegisterOnShutdown(trigger.onServerShutdown)

	logtool.Debug("trigger", "http", fmt.Sprintf("create http listener at %s...", addr))

	if err := trigger.currentServer.ListenAndServe(); err != nil {
		logtool.Error("trigger", "http", fmt.Sprintf("create http listener at %s failed: %v", addr, err))
	}
}

func (trigger *HttpPlugin) Stop() {

	if trigger.currentServer != nil {

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := trigger.currentServer.Shutdown(ctx); err != nil {
			logtool.Error("trigger", "http", fmt.Sprintf("http listener shutdown '%s' failed: %v", trigger.currentServer.Addr, err))
		}
	}
}

func (trigger *HttpPlugin) handleRequestFunc(w nethttp.ResponseWriter, r *nethttp.Request) {

	logtool.Debug("trigger", "http", fmt.Sprintf("http %s %s%s", r.Method, r.Host, r.URL))

	if existed, _ := arraytool.InArray(r.Method, allowHttpMethods); !existed {
		nethttp.Error(w, fmt.Sprintf("http method '%s' not allowed", r.Method), nethttp.StatusMethodNotAllowed)
		logtool.Error("trigger", "http", fmt.Sprintf("http method '%s' not allowed", r.Method))
		return
	}

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		nethttp.Error(w, "can't read body", nethttp.StatusBadRequest)
		logtool.Error("trigger", "http", fmt.Sprintf("read http response body failed : %v", err))
		return
	}

	var bodyContent, contentType string

	if strings.EqualFold(r.Method, "POST") {
		bodyContent = string(body)
		contentType = r.Header.Get("Content-Type")

		switch contentType {
		case "application/json":
		default:
			nethttp.Error(w, fmt.Sprintf("Content-Type '%s' not supported", contentType), nethttp.StatusUnsupportedMediaType)
			logtool.Error("trigger", "http", fmt.Sprintf("Content-Type '%s' not supported", contentType))
			return
		}
	}

	var triggerPlugin pluginbase.ITriggerPlugin = trigger
	trigger.FireAction(&triggerPlugin, &bodyContent)
}

func (trigger *HttpPlugin) onServerShutdown() {
	logtool.Debug("trigger", "http", fmt.Sprintf("http listener '%s' shutdown...", trigger.currentServer.Addr))
}
