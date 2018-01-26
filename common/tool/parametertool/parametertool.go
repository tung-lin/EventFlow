package parametertool

import (
	"strings"
)

var replacer *strings.Replacer

func init() {
	replacer = strings.NewReplacer("%", "", "{", "", "}", "")
}

func GetParameterKey(parameter string) (key string) {
	return replacer.Replace(parameter)
}
