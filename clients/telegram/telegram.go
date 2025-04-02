package telegram

import (
	"encoding/json"
	"io"
	"link-saver-bot/lib/e"
	"log/slog"
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
	basePath := generateBasePath(token)
	slog.Info("Initializing Telegram client", "host", host, "basePath", basePath)

	return &Client{
		host:     host,
		basePath: basePath,
		client:   http.Client{},
	}
}

func generateBasePath(token string) string {
	return "bot" + token
}

// SendMessage sends message to telegram  https://core.telegram.org/bots/api#sendmessage
func (client *Client) SendMessage(chatId int, text string) error {
	query := url.Values{}
	query.Add("chat_id", strconv.Itoa(chatId))
	query.Add("text", text)

	slog.Info("Sending message to Telegram", "chatId", chatId, "textLength", len(text))

	_, err := client.doRequest("sendMessage", query)
	if err != nil {
		slog.Error("Failed to send message", "chatId", chatId, "error", err)
		return e.Wrap("can't send message", err)
	}

	slog.Info("Message sent successfully", "chatId", chatId)
	return nil
}

// Updates fetches updates from telegram
func (client *Client) Updates(limit int, offset int) ([]Updates, error) {
	query := url.Values{}
	query.Add("offset", strconv.Itoa(offset))
	query.Add("limit", strconv.Itoa(limit))

	slog.Info("Fetching updates from Telegram", "limit", limit, "offset", offset)

	data, err := client.doRequest("getUpdates", query)
	if err != nil {
		slog.Error("Failed to get updates", "limit", limit, "offset", offset, "error", err)
		return nil, err
	}

	var result UpdateResponse

	if err := json.Unmarshal(data, &result); err != nil {
		slog.Error("Failed to unmarshal updates response", "error", err)
		return nil, err
	}

	slog.Info("Updates fetched successfully", "count", len(result.Result))
	return result.Result, nil
}

func (client *Client) doRequest(method string, query url.Values) (data []byte, err error) {
	u := url.URL{
		Scheme: "https",
		Host:   client.host,
		Path:   path.Join(client.basePath, method),
	}

	request, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		slog.Error("Failed to create request", "method", method, "error", err)
		return nil, e.Wrap("can't create new request", err)
	}

	request.URL.RawQuery = query.Encode()
	requestURL := request.URL.String()

	slog.Info("Making request to Telegram API", "method", method, "url", requestURL)

	response, err := client.client.Do(request)
	if err != nil {
		slog.Error("Failed to execute request", "method", method, "url", requestURL, "error", err)
		return nil, e.Wrap("can't do request", err)
	}

	defer func() { _ = response.Body.Close() }()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		slog.Error("Failed to read response body", "method", method, "error", err)
		return nil, e.Wrap("can't read response", err)
	}

	slog.Debug("Received response from Telegram",
		"statusCode", response.StatusCode,
		"body", string(body))

	return body, nil
}
