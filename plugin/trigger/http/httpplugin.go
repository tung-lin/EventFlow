package http

import (
	"EventFlow/common/interface/pluginbase"
	"EventFlow/common/tool/arraytool"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	nethttp "net/http"

	"github.com/gorilla/mux"
)

type HttpPlugin struct {
	currentServer *nethttp.Server
	pluginbase.PolicyHandler
	Setting SettingConfig
}

var allowHttpMethods []string

func (trigger *HttpPlugin) Start() {

	allowHttpMethods = []string{"GET", "POST"}
	handler := trigger.handleRequestFunc

	router := mux.NewRouter()
	router.HandleFunc("/", handler).Methods("POST")
	router.HandleFunc("/", handler).Methods("GET")
	router.HandleFunc("/", handler).Methods("PUT")
	router.HandleFunc("/", handler).Methods("DELETE")
	router.HandleFunc("/", handler).Methods("HEAD")

	addr := fmt.Sprintf(":%d", trigger.Setting.Port)

	trigger.currentServer = &nethttp.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	trigger.currentServer.RegisterOnShutdown(trigger.onServerShutdown)

	log.Printf("[trigger][http] Listening at %s", addr)

	if err := trigger.currentServer.ListenAndServe(); err != nil {
		log.Printf("[trigger][http] Listening at %s failed\r\n%s", addr, err)
	}
}

func (trigger *HttpPlugin) Stop() {

	if trigger.currentServer != nil {

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := trigger.currentServer.Shutdown(ctx); err != nil {
			log.Printf("Server shutdown failed:%s", err)
		}
	}
}

func (trigger *HttpPlugin) handleRequestFunc(w nethttp.ResponseWriter, r *nethttp.Request) {

	log.Printf("http %s %s%s", r.Method, r.Host, r.URL)

	if existed, _ := arraytool.InArray(r.Method, allowHttpMethods); existed {

		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			log.Printf("[trigger][http] Read http response body failed : %v", err)
			nethttp.Error(w, "can't read body", nethttp.StatusBadRequest)
			return
		}

		bodyContent := string(body)

		var triggerPlugin pluginbase.ITriggerPlugin = trigger
		trigger.FireAction(&triggerPlugin, &bodyContent)

	} else {
		nethttp.Error(w, fmt.Sprintf("[trigger][http] Method '%s' not allowed", r.Method), nethttp.StatusMethodNotAllowed)
	}

}

func (trigger *HttpPlugin) onServerShutdown() {
	log.Printf("[trigger][http] Server '%s' shutdown...", trigger.currentServer.Addr)
}
