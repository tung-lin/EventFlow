package jsontool

import (
	"encoding/json"
	"log"
)

func UnmarshalToType(in interface{}, out interface{}) (err error) {

	bytes := structToByteArray(in)
	err = json.Unmarshal(bytes, out)

	return err
}

func structToByteArray(setting interface{}) (bytes []byte) {

	if setting == nil {
		return nil
	}

	bytes, err := json.Marshal(setting)

	if err != nil {
		log.Printf("Marshal json config to byte array failed: %v", err)
		return nil
	}

	return bytes
}
