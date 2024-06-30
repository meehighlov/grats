package telegram

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"time"
	"strconv"

	"encoding/json"
)

// --------------------------------------------------------------- Types ---------------------------------------------------------------

type requestBodyType map[string]interface{}
type requestQueryParamsType map[string]string
type UpdatesChannel chan Update

type telegramClient struct {
	urlHead string
	token   string
	baseUrl string

	// todo use interface instead of concrete type
	httpClient *http.Client
}

type APICaller interface {
	SendMessage(chatId, text string, needForceReply bool) *Message
	GetUpdates(updatesOffset int) (*UpdateResponse, error)
	GetMe() (*User, error)
	GetUpdatesChannel() UpdatesChannel
}

// --------------------------------------------------------------- telegram client  ---------------------------------------------------------------

func NewClient(token string) APICaller {
	// http client timeout > telegram getUpdates timeout
	httpClient := &http.Client{Timeout: 20 * time.Second}
	urlHead := "https://api.telegram.org/bot"
	return telegramClient{
		token: token,
		urlHead: urlHead,
		baseUrl: urlHead + token,
		httpClient: httpClient,
	}
}

func (tc *telegramClient) send(request *http.Request) (*http.Response, error) {
	response, err := tc.httpClient.Do(request)

	if err != nil {
		log.Println("HTTP request failed", err.Error())
		return nil, err
	}

	return response, nil
}

func (tc *telegramClient) prepareRequestBody(requestBody *requestBodyType) (io.Reader, error) {
	if requestBody == nil {
		return nil, nil
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		log.Println("Failed to marshal request body " + err.Error())
		return nil, err
	}

	return bytes.NewBuffer(jsonData), nil
}

func (tc *telegramClient) prepareQueryParams(queryParams *requestQueryParamsType, requset *http.Request) error {
	if queryParams == nil {
		return nil
	}

	query := requset.URL.Query()

	for key, value := range *queryParams {
		query.Add(key, value)
	}

	requset.URL.RawQuery = query.Encode()

	return nil
}

func (tc *telegramClient) prepareRequest(method, urlTail string, requestBody *requestBodyType, queryParams *requestQueryParamsType) (*http.Request, error) {
	body, err := tc.prepareRequestBody(requestBody)

	if err != nil {
		return nil, err
	}

	url := tc.baseUrl + "/" + urlTail
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		log.Fatalf("Failed to create request " + err.Error())
	}
	req.Header.Add("Content-Type", "application/json")

	tc.prepareQueryParams(queryParams, req)

	return req, nil
}

func (tc *telegramClient) getBodyBytes(response *http.Response) ([]byte, error) {
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println("Failed to parse response body " + err.Error())
		return nil, err
	}

	return body, nil
}

func (tc *telegramClient) sendRequest(method, urlTail string, body *requestBodyType, queryParams *requestQueryParamsType, responseModel interface{}) error {
	request, err := tc.prepareRequest(method, urlTail, body, queryParams)

	if err != nil {
		return err
	}

	response, err := tc.send(request)

	if err != nil {
		return err
	}

	body_bytes, read_err := tc.getBodyBytes(response)

	if read_err != nil {
		return read_err
	}

	if response.StatusCode == http.StatusOK {
		json.Unmarshal(body_bytes, responseModel)
		return nil
	}

	log.Println("Bad status code:", response.StatusCode, "body:", string(body_bytes))

	return nil
}

// --------------------------------------------------------------- API methods implementation ---------------------------------------------------------------

func (tc telegramClient) SendMessage(chatId, text string, needForceReply bool) *Message {
	res := Message{}

	body := requestBodyType{
		"chat_id": chatId,
		"text": text,
		"reply_markup": requestBodyType{
			"force_reply": needForceReply,
			"selective": needForceReply,
		},
	}

	tc.sendRequest("POST", "sendMessage", &body, nil, res)

	return &res
}

func (tc telegramClient) GetUpdates(updatesOffset int) (*UpdateResponse, error) {
	res := UpdateResponse{}

	queryParams := requestQueryParamsType{
		// timeout should be less than http client timeout
		"timeout": "10",
		"offset": strconv.Itoa(updatesOffset),
	}

	err := tc.sendRequest("GET", "getUpdates", nil, &queryParams, &res)

	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (tc telegramClient) GetMe() (*User, error) {
	res := GetMeReponse{}
	tc.sendRequest("GET", "getMe", nil, nil, &res)
	return &res.Result, nil
}

// --------------------------------------------------------------- polling ---------------------------------------------------------------

func (tc telegramClient) GetUpdatesChannel() UpdatesChannel {
	updatesChannelSize := 100
	updatesOffset := -1

	ch := make(chan Update, updatesChannelSize)

	go func() {
		for {
			updates, err := tc.GetUpdates(updatesOffset)
			if err != nil {
				log.Println(err)
				log.Println("Failed to get updates, retrying in 3 seconds...")
				time.Sleep(time.Second * time.Duration(3))

				continue
			}

			for _, update := range updates.Result {
				if update.UpdateId >= updatesOffset {
					updatesOffset = update.UpdateId + 1
					ch <- update
				}
			}
		}
	}()

	return ch
}
