package yamltool

import (
	"EventFlow/common/tool/logtool"
	"fmt"

	"gopkg.in/yaml.v2"
)

//UnmarshalToType converts unkonwn type in and assigns values into the out value
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
		logtool.Error("tool", "yaml", fmt.Sprintf("marshal yaml config to byte array failed: %v", err))
		return nil
	}

	return bytes
}
