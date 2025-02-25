package telegram

import (
	"encoding/json"
	"io"
	"link-saver-bot/lib/e"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

// client for interacting with telegram api
type Client struct {
	host     string
	basePath string // example:tg-bot.com/bot{token}
	client   http.Client
}

func New(host, token string) *Client {

	return &Client{
		host:     host,
		basePath: generateBasePath(token),
		client:   http.Client{},
	}
}

func generateBasePath(token string) string {
	return "bot" + token
}

// SendMessage sends message to telegram  https://core.telegram.org/bots/api#sendmessage
func (client *Client) SendMessage(chatId int, text string) error {
	query := url.Values{}
	query.Add("chatId", strconv.Itoa(chatId))
	query.Add("text", text)

	_, err := client.doRequest("sendMessage", query)

	if err != nil {
		return e.Wrap("can't send message", err)
	}

	return nil
}

// Updates fetches updates from telegram
func (client *Client) Updates(offset, limit int) ([]Updates, error) {

	query := url.Values{}
	query.Add("offset", strconv.Itoa(offset))
	query.Add("limit", strconv.Itoa(limit))

	//@todo make "getUpdates" http request
	data, err := client.doRequest("getUpdates", query)
	if err != nil {
		return nil, err
	}

	var result UpdateResponse

	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return result.Result, nil
}

func (client *Client) doRequest(method string, query url.Values) (data []byte, err error) {

	defer func() { err = e.Wrap("can't do request", err) }()

	u := url.URL{
		Scheme: "https",
		Host:   client.host,
		Path:   path.Join(client.basePath, method),
	}

	request, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, e.Wrap("can't create new request", err)
	}

	request.URL.RawQuery = query.Encode()

	response, err := client.client.Do(request)
	if err != nil {
		return nil, e.Wrap("can't do request", err)
	}

	defer func() { _ = response.Body.Close() }()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, e.Wrap("can't read response", err)
	}

	return body, nil
}
