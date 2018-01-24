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

func CreateRandomString() string {

	result := make([]byte, defaultStringLen)

	for i := range result {
		result[i] = chars[random.Intn(len(chars))]
	}

	return string(result)
}

func CreateRandomStringWithLen(strlen int) string {

	result := make([]byte, strlen)

	for i := range result {
		result[i] = chars[random.Intn(len(chars))]
	}

	return string(result)
}