package stringtool

import (
	"math/rand"
	"time"
)

var random *rand.Rand

const defaultStringLen = 8
const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func init() {
	random = rand.New(rand.NewSource(time.Now().UnixNano()))
}

//CreateRandomString generates a eight bytes random string
func CreateRandomString() string {

	result := make([]byte, defaultStringLen)

	for i := range result {
		result[i] = chars[random.Intn(len(chars))]
	}

	return string(result)
}

//CreateRandomStringWithLen generates a strlen bytes random string
func CreateRandomStringWithLen(strlen int) string {

	result := make([]byte, strlen)

	for i := range result {
		result[i] = chars[random.Intn(len(chars))]
	}

	return string(result)
}

//Equal reports whether s and t are equal
func Equal(s, t string, caseSensitive bool) bool {

	return s == t
}

//ToStringArray convert interface(string or string array) to string array
func ToStringArray(data interface{}) []string {

	if value, ok := data.(string); ok {
		return []string{value}
	} else if value, ok := data.([]interface{}); ok {

		var results []string

		for _, v := range value {
			if stringValue, ok := v.(string); ok {
				results = append(results, stringValue)
			}
		}

		return results
	}

	return nil
}
