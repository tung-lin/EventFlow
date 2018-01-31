package line

import (
	"EventFlow/common/tool/parametertool"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const lineNotifyMethod = "POST"
const lineNotifyURL = "https://notify-api.line.me/api/notify"
const lineNotifyContentType = "application/x-www-form-urlencoded"
const limit = "X-RateLimit-Limit"
const remaining = "X-RateLimit-Remaining"
const imagelimit = "X-RateLimit-ImageLimit"
const imageremaining = "X-RateLimit-ImageRemaining"
const reset = "X-RateLimit-Reset"

type LinePlugin struct {
	Setting SettingConfig
}

type notifyResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type rateLimit struct {
	Limit          int
	Remaining      int
	ImageLimit     int
	ImageRemaining int
	Reset          time.Time
}

func (filter *LinePlugin) FireAction(messageFromTrigger *string, parameters *map[string]interface{}) {

	message := parametertool.ReplaceWithParameter(&filter.Setting.Message, parameters)

	values := url.Values{}
	values.Add("message", message)

	body := strings.NewReader(values.Encode())
	request, err := http.NewRequest(lineNotifyMethod, lineNotifyURL, body)

	if err != nil {
		log.Printf("[action][line] create http request failed: %v", err)
		return
	}

	request.Header.Set("Content-Type", lineNotifyContentType)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", filter.Setting.AccessToken))

	client := http.DefaultClient
	response, err := client.Do(request)

	if err != nil {
		log.Printf("[action][line] send http request failed: %v", err)
		return
	}

	defer response.Body.Close()

	notifyResponse := notifyResponse{}
	err = json.NewDecoder(response.Body).Decode(&notifyResponse)

	if err != nil {
		log.Printf("[action][line] decode http response body failed: %v", err)
		return
	}

	log.Printf("[action][line] http response message: %s, status: %d", notifyResponse.Message, notifyResponse.Status)

	if response.StatusCode == http.StatusOK {
		rateLimit := parseRateLimit(response.Header)
		log.Printf("[action][line] limit: %d, remaining: %d, imagelimit: %d, imageremaining: %d, reset: %s", rateLimit.Limit, rateLimit.Remaining, rateLimit.ImageLimit, rateLimit.ImageRemaining, rateLimit.Reset.Format(time.RFC1123))
	}
}

func parseRateLimit(header http.Header) rateLimit {

	rateLimit := rateLimit{}

	if v, err := strconv.Atoi(header.Get(limit)); err == nil {
		rateLimit.Limit = v
	}

	if v, err := strconv.Atoi(header.Get(remaining)); err == nil {
		rateLimit.Remaining = v
	}

	if v, err := strconv.Atoi(header.Get(imagelimit)); err == nil {
		rateLimit.ImageLimit = v
	}

	if v, err := strconv.Atoi(header.Get(imageremaining)); err == nil {
		rateLimit.ImageRemaining = v
	}

	if v, err := strconv.ParseInt(header.Get(reset), 10, 64); err == nil {
		rateLimit.Reset = time.Unix(v, 0)
	}

	return rateLimit
}
