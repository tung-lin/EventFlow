package config

import (
	"io/ioutil"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

type Configuration struct {
	Port string `yaml:port`
}

var Config Configuration

func init() {

	configFilePath, _ := os.Getwd()
	file, err := ioutil.ReadFile(configFilePath + "/datacreator.yaml")

	if err != nil {
		log.Fatalf("Read config file failed: %v", err)
	}

	err = yaml.Unmarshal(file, &Config)

	if err != nil {
		log.Fatalf("Unmarshal config file failed: %v", err)
	}

	GetEnv("port", &Config.Port)
}

func GetEnv(envName string, variable *string) {

	envName = strings.ToLower(envName)
	value, existed := os.LookupEnv(envName)

	if existed {
		*variable = value
	}
}
