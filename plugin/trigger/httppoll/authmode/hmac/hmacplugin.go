package hmac

import (
	"EventFlow/common/tool/arraytool"
	"EventFlow/plugin/trigger/httppoll/authbase"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"hash"
	"net/http"
	"time"
)

type HMacPlugin struct {
	Setting      hmacSettingConfig
	hashFunction func() hash.Hash
}

var allowHMacAlgorithm = []string{"hmac-sha1", "hmac-sha256", "hmac-sha512"}

func (auth *HMacPlugin) CheckParameter() error {

	if auth.Setting.Algorithm == "" {
		return fmt.Errorf("http polling is diasbled due to empty hmac algorithm setting")
	}

	if existed, _ := arraytool.InArray(auth.Setting.Algorithm, allowHMacAlgorithm); !existed {
		return fmt.Errorf(fmt.Sprintf("hmac algorithm '%s' is not supported", auth.Setting.Algorithm))
	}

	switch auth.Setting.Algorithm {
	case "hmac-sha1":
		auth.hashFunction = sha1.New
	case "hmac-sha256":
		auth.hashFunction = sha256.New
	case "hmac-sha512":
		auth.hashFunction = sha512.New
	}

	return nil
}

func (auth *HMacPlugin) CreateHttpHeaders() []authbase.Header {

	var headers []authbase.Header

	switch auth.Setting.Type {
	case "ptx":
		xDate := time.Now().UTC().Format(http.TimeFormat)
		encryptXDate := "x-date: " + xDate

		key := []byte(auth.Setting.APPKey)
		mac := hmac.New(auth.hashFunction, key)
		mac.Write([]byte(encryptXDate))

		encryptSign := base64.StdEncoding.EncodeToString(mac.Sum(nil))
		encryptSign = "hmac username=\"" + auth.Setting.APPID + "\", algorithm=\"" + auth.Setting.Algorithm + "\", headers=\"x-date\", signature=\"" + encryptSign + "\""

		headers = append(headers, authbase.Header{Key: "Authorization", Value: encryptSign}, authbase.Header{Key: "x-date", Value: xDate})
	}

	return headers
}
