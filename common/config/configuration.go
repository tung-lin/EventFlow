package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

type configuration struct {
	Log struct {
		Level         string `yaml:"level"`
		Path          string `yaml:"path"`
		Verbose       bool   `yaml:"verbose"`
		MaxSizeKB     int    `yaml:"maxsizekb"`
		IncludeCaller bool   `yaml:"includecaller"`
	} `yaml:"log"`
}

//Config represents a global configuration
var Config configuration

func init() {

	configFilePath, _ := os.Getwd()
	file, err := ioutil.ReadFile(configFilePath + "/config/config.yaml")

	if err != nil {
		log.Printf("[tool][config] %s", fmt.Sprintf("read config file failed: %v", err))
		os.Exit(1)
	}

	err = yaml.Unmarshal(file, &Config)

	if err != nil {
		log.Printf("[tool][config] %s", fmt.Sprintf("unmarshal config file failed: %v", err))
		os.Exit(1)
	}

	getEnv("log.level", &Config.Log.Level)
	getEnv("log.path", &Config.Log.Path)

	if Config.Log.MaxSizeKB == 0 {
		Config.Log.MaxSizeKB = 1024
	}

	if Config.Log.Level == "" {
		Config.Log.Level = "info"
	}
}

func getEnv(envName string, variable *string) {

	envName = strings.ToLower(envName)
	value, existed := os.LookupEnv(envName)

	if existed {
		*variable = value
	}
}
