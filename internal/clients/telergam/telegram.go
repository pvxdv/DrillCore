package tgClient

import (
	"bytes"
	"drillCore/internal/config"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

type Client struct {
	host     string
	basePath string
	client   http.Client
}

const (
	getUpdatesMethod  = "getUpdates"
	sendMessageMethod = "sendMessage"
)

func New(cfg *config.TelegramEnvs) *Client {
	return &Client{
		host:     cfg.BaseUrl,
		basePath: newBasePath(cfg.Token),
		client:   http.Client{},
	}
}

func newBasePath(token string) string {
	return "bot" + token
}

func (c *Client) Updates(offset int, limit int) ([]Update, error) {
	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))
	q.Add("timeout", "60")

	data, err := c.doRequest(getUpdatesMethod, q)
	if err != nil {
		return nil, err
	}

	var res UpdatesResponse
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, fmt.Errorf("failed to unmarshal updates response: %w", err)
	}

	if !res.Ok {
		return nil, errors.New("failed to get updates: telegram api returned not ok")
	}

	return res.Result, nil
}

func (c *Client) SendMessage(chatID int, text string) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("text", text)

	_, err := c.doRequest(sendMessageMethod, q)
	if err != nil {
		return fmt.Errorf("failed to send message:%w", err)
	}

	return nil
}

func (c *Client) SendMessageWithKeyboard(chatID int, text string, keyboard ReplyMarkup) error {
	req := struct {
		ChatID      int         `json:"chat_id"`
		Text        string      `json:"text"`
		ReplyMarkup ReplyMarkup `json:"reply_markup"`
	}{
		ChatID:      chatID,
		Text:        text,
		ReplyMarkup: keyboard,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("can't marshal request: %w", err)
	}

	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, sendMessageMethod),
	}

	resp, err := c.client.Post(u.String(), "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("can't do request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("telegram error: %s", string(body))
	}

	return nil
}

func (c *Client) doRequest(method string, query url.Values) (data []byte, err error) {
	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request:%w", err)
	}

	req.URL.RawQuery = query.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do request:%w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	data, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response:%w", err)
	}

	return data, nil
}
