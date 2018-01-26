package parametertool

import (
	"regexp"
	"strings"
)

var regex = regexp.MustCompile(`{[^}]*}`)
var replacer = strings.NewReplacer("{", "", "}", "")

func ReplaceWithParameter(target *string, parameters *map[string]interface{}) {

	matches := regex.FindAllString(*target, -1)

	for _, matchItem := range matches {

		parmKey := replacer.Replace(matchItem)
		parmValue := (*parameters)[parmKey]

		if newValue, ok := parmValue.(string); ok {
			*target = strings.Replace(*target, "%"+matchItem, newValue, -1)
		}
	}
}
