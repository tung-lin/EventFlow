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
