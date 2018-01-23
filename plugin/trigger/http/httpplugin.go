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
	Setting SettingConfig
}

var actionList []*pluginbase.IActionPlugin
var currentServer *nethttp.Server
var allowHttpMethods []string

func (trigger *HttpPlugin) AddAction(actionPlugin *pluginbase.IActionPlugin) {
	actionList = append(actionList, actionPlugin)
}

func (trigger *HttpPlugin) Start() {

	allowHttpMethods = []string{"GET", "POST"}
	handler := trigger.handleRequestFunc

	router := mux.NewRouter()
	router.HandleFunc("/", handler).Methods("POST")
	router.HandleFunc("/", handler).Methods("GET")
	router.HandleFunc("/", handler).Methods("PUT")
	router.HandleFunc("/", handler).Methods("DELETE")
	router.HandleFunc("/", handler).Methods("HEAD")

	currentServer = &nethttp.Server{
		Addr:    fmt.Sprintf(":%d", trigger.Setting.Port),
		Handler: router,
	}

	currentServer.RegisterOnShutdown(onServerShutdown)
	currentServer.ListenAndServe()
}

func (trigger *HttpPlugin) Stop() {

	if currentServer != nil {

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := currentServer.Shutdown(ctx); err != nil {
			log.Printf("Server shutdown failed:%s", err)
		}
	}
}

func (trigger *HttpPlugin) handleRequestFunc(w nethttp.ResponseWriter, r *nethttp.Request) {

	log.Printf("http %s %s%s", r.Method, r.Host, r.URL)

	if existed, _ := arraytool.InArray(r.Method, allowHttpMethods); existed {

		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			log.Printf("Error reading body: %v", err)
			nethttp.Error(w, "can't read body", nethttp.StatusBadRequest)
			return
		}

		for _, action := range actionList {
			(*action).FireAction(trigger, string(body))
		}

	} else {
		nethttp.Error(w, fmt.Sprintf("Method '%s' not allowed", r.Method), nethttp.StatusMethodNotAllowed)
	}

}

func onServerShutdown() {
	log.Printf("Server '%s' shutdown...", currentServer.Addr)
}
