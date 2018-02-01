package jsontool

import (
	"EventFlow/common/tool/logtool"
	"encoding/json"
	"fmt"
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
		logtool.Error("tool", "json", fmt.Sprintf("marshal json config to byte array failed: %v", err))
		return nil
	}

	return bytes
}
