package telegram

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"
)

func GetUpdates() []Update {
	response, err := request("/getUpdates")

	defer response.Body.Close()

	var updatesResponse struct {
		OK     bool     `json:"ok"`
		Result []Update `json:"result"`
	}

	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&updatesResponse)

	if err != nil {
		panic("Parsing the response body failed!")
	}

	if !updatesResponse.OK {
		panic("Fetching updates from telegram failed!")
	}

	return updatesResponse.Result
}

// func SendMessage(chatID int, message string) map[string]interface{} {
// 	params := make(map[string]string)
// 	params["chat_id"] = strconv.Itoa(chatID)
// 	params["text"] = message

// 	return request("/sendMessage", params)
// }

func request(path string, params ...map[string]string) (*http.Response, error) {
	requestUrl := url.URL{
		Scheme: "https",
		Host:   "api.telegram.org",
		Path:   "bot" + os.Getenv("TELEGRAM_BOT_TOKEN") + path,
	}

	if len(params) > 0 {
		requestParams := url.Values{}

		for key, value := range params[0] {
			requestParams.Add(key, value)
		}

		requestUrl.RawQuery = requestParams.Encode()
	}

	response, err := http.Get(requestUrl.String())

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	return response, nil
}
