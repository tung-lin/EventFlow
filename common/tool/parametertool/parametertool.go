package parametertool

import (
	"EventFlow/common/tool/logtool"
	"fmt"
	"regexp"
	"strings"
)

var regex = regexp.MustCompile(`{[^}]*}`)
var replacer = strings.NewReplacer("{", "", "}", "")

func ReplaceWithParameter(target *string, parameters *map[string]interface{}) (replaceResult string) {

	result := *target
	matches := regex.FindAllString(result, -1)

	for _, matchItem := range matches {

		parmKey := replacer.Replace(matchItem)
		parmValue, existed := (*parameters)[parmKey]

		if !existed {
			parmValue = ""
		}

		matchItem = "%" + matchItem

		if value, ok := parmValue.(string); ok {
			result = strings.Replace(result, matchItem, value, -1)
		} else if values, ok := parmValue.([]interface{}); ok {
			results := make([]string, len(values))
			for index, value := range values {
				results[index] = value.(string)
			}
			result = strings.Replace(result, matchItem, strings.Join(results, ","), -1)
		} else {
			logtool.Error("tool", "parameter", fmt.Sprintf("cannot replace parameter %s", matchItem))
		}
	}

	return result
}
