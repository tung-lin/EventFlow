package jsontool

import (
	"EventFlow/common/tool/logtool"
	"encoding/json"
	"fmt"
	"unsafe"
)

//UnmarshalToType converts unkonwn type in and assigns values into the out value
func UnmarshalToType(in interface{}, out interface{}) (err error) {

	bytes := structToByteArray(in)
	err = json.Unmarshal(bytes, out)

	return err
}

func MarshalToString(in interface{}, out *string) (err error) {
	bytes, err := json.Marshal(in)

	if err != nil {
		return err
	}

	*out = *(*string)(unsafe.Pointer(&bytes))

	return nil
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
