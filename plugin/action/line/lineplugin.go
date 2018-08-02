package line

import (
	"EventFlow/common/tool/logtool"
	"EventFlow/common/tool/parametertool"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	lineNotifyMethod      = "POST"
	lineNotifyURL         = "https://notify-api.line.me/api/notify"
	lineNotifyContentType = "application/x-www-form-urlencoded"
	limit                 = "X-RateLimit-Limit"
	remaining             = "X-RateLimit-Remaining"
	imagelimit            = "X-RateLimit-ImageLimit"
	imageremaining        = "X-RateLimit-ImageRemaining"
	reset                 = "X-RateLimit-Reset"
)

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
	accessTokens := parametertool.ReplaceWithParameter(&filter.Setting.AccessToken, parameters)

	values := url.Values{}
	values.Add("message", message)

	body := strings.NewReader(values.Encode())

	for _, accessToken := range strings.Split(accessTokens, ",") {

		go func(token string) {
			request, err := http.NewRequest(lineNotifyMethod, lineNotifyURL, body)
			request.Header.Set("Content-Type", lineNotifyContentType)
			request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

			client := http.DefaultClient
			response, err := client.Do(request)

			if err != nil {
				logtool.Error("action", "line", fmt.Sprintf("send http request failed: %v", err))
				return
			}

			defer response.Body.Close()

			notifyResponse := notifyResponse{}
			err = json.NewDecoder(response.Body).Decode(&notifyResponse)

			if err != nil {
				logtool.Error("action", "line", fmt.Sprintf("decode http response body failed: %v", err))
				return
			}

			logtool.Debug("action", "line", fmt.Sprintf("http response message: %s, status: %d", notifyResponse.Message, notifyResponse.Status))

			if response.StatusCode == http.StatusOK {
				rateLimit := parseRateLimit(response.Header)
				logtool.Debug("action", "line", fmt.Sprintf("limit: %d, remaining: %d, imagelimit: %d, imageremaining: %d, reset: %s", rateLimit.Limit, rateLimit.Remaining, rateLimit.ImageLimit, rateLimit.ImageRemaining, rateLimit.Reset.Format(time.RFC1123)))
			}
		}(accessToken)
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
