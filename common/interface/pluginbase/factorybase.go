package pluginbase

import (
	"log"

	"gopkg.in/yaml.v2"
)

type IFactoryBase interface {
	UnmarshalToType(in interface{}, out interface{}) (err error)
}

type FactoryBase struct {
}

func (plugin *FactoryBase) UnmarshalToType(in interface{}, out interface{}) (err error) {
	bytes := structToByteArray(in)
	err = yaml.Unmarshal(bytes, &out)

	return err
}

func structToByteArray(setting interface{}) (bytes []byte) {

	if setting == nil {
		return nil
	}

	str, err := yaml.Marshal(setting)

	if err != nil {
		log.Printf("Marshal setting config to string failed: %v", err)
		return nil
	}

	return []byte(str)
}
