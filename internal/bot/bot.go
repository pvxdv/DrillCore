package bot

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"

	"drillCore/internal/config"

	"go.uber.org/zap"
)

type Client struct {
	host     string
	basePath string
	client   http.Client
	logger   *zap.SugaredLogger
}

const (
	getUpdatesMethod  = "getUpdates"
	sendMessageMethod = "sendMessage"
)

func New(cfg *config.TelegramEnvs, logger *zap.SugaredLogger) *Client {
	return &Client{
		host:     cfg.BaseUrl,
		basePath: newBasePath(cfg.Token),
		client:   http.Client{},
		logger:   logger,
	}
}

func newBasePath(token string) string {
	return "bot" + token
}

func (c *Client) Updates(ctx context.Context, offset int, limit int) ([]Update, error) {
	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))
	q.Add("timeout", "60")

	data, err := c.doRequest(ctx, getUpdatesMethod, q)
	if err != nil {
		return nil, err
	}

	var res UpdatesResponse
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, fmt.Errorf("failed to unmarshal updates response: %w", err)
	}

	if !res.Ok {
		c.logger.Debugf("failed to get updates: event-processor api returned not ok:%v", err)
	}

	return res.Result, nil
}

func (c *Client) SendMessage(ctx context.Context, chatID int, text string) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("text", text)

	_, err := c.doRequest(ctx, sendMessageMethod, q)
	if err != nil {
		return fmt.Errorf("failed to send message:%w", err)
	}

	return nil
}

func (c *Client) SendMessageWithKeyboard(ctx context.Context, chatID int, text string, keyboard ReplyMarkup) error {
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
		return fmt.Errorf("event-processor error: %s", string(body))
	}

	return nil
}

func (c *Client) SendPhotoWithKeyBoard(ctx context.Context, chatID int, photoPath, caption string, keyboard ReplyMarkup) error {
	file, err := os.Open(photoPath)
	if err != nil {
		return fmt.Errorf("failed to open photo: %w", err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("photo", filepath.Base(photoPath))
	if err != nil {
		return fmt.Errorf("failed to create form file: %w", err)
	}
	if _, err = io.Copy(part, file); err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	_ = writer.WriteField("chat_id", strconv.Itoa(chatID))
	if caption != "" {
		_ = writer.WriteField("caption", caption)
	}

	if keyboard.InlineKeyboard != nil {
		kbData, err := json.Marshal(keyboard)
		if err != nil {
			return fmt.Errorf("failed to marshal keyboard: %w", err)
		}
		_ = writer.WriteField("reply_markup", string(kbData))
	}

	if err = writer.Close(); err != nil {
		return fmt.Errorf("failed to close writer: %w", err)
	}

	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, "sendPhoto"),
	}

	req, err := http.NewRequestWithContext(ctx, "POST", u.String(), body)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send photo: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("telegram API error: %s", string(body))
	}

	return nil
}

func (c *Client) doRequest(ctx context.Context, method string, query url.Values) (data []byte, err error) {
	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
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
