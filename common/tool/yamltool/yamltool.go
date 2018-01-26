package yamltool

import (
	"log"

	"gopkg.in/yaml.v2"
)

func UnmarshalToType(in interface{}, out interface{}) (err error) {

	bytes := structToByteArray(in)
	err = yaml.Unmarshal(bytes, out)

	return err
}

func structToByteArray(setting interface{}) (bytes []byte) {

	if setting == nil {
		return nil
	}

	bytes, err := yaml.Marshal(setting)

	if err != nil {
		log.Printf("Marshal yaml config to byte array failed: %v", err)
		return nil
	}

	return bytes
}
