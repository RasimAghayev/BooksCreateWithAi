package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

type Client struct {
	BotToken string
	ChatID   string
	Timeout  time.Duration
}

func NewClient(botToken, chatID string, timeout time.Duration) *Client {
	return &Client{
		BotToken: botToken,
		ChatID:   chatID,
		Timeout:  timeout,
	}
}

type SendDocumentResponse struct {
	OK          bool   `json:"ok"`
	MessageID   int64  `json:"message_id"`
	FileID      string `json:"file_id"`
	Error       string `json:"description,omitempty"`
}

func (c *Client) SendDocument(ctx context.Context, filename string, jsonData []byte, caption string) (*SendDocumentResponse, error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	part, err := writer.CreateFormFile("document", filename)
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}

	if _, err := io.Copy(part, bytes.NewReader(jsonData)); err != nil {
		return nil, fmt.Errorf("failed to copy data: %w", err)
	}

	if err := writer.WriteField("chat_id", c.ChatID); err != nil {
		return nil, fmt.Errorf("failed to write chat_id: %w", err)
	}

	if err := writer.WriteField("caption", caption); err != nil {
		return nil, fmt.Errorf("failed to write caption: %w", err)
	}

	writer.Close()

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendDocument", c.BotToken)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, &body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{Timeout: c.Timeout}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result SendDocumentResponse
	if err := parseTelegramResponse(respBody, &result); err != nil {
		return nil, err
	}

	if !result.OK {
		return nil, fmt.Errorf("telegram error: %s", result.Error)
	}

	return &result, nil
}

func parseTelegramResponse(data []byte, result *SendDocumentResponse) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if ok, okCast := raw["ok"].(bool); !okCast || !ok {
		if desc, okDesc := raw["description"].(string); okDesc {
			result.Error = desc
		}
		return fmt.Errorf("telegram API returned error")
	}

	if res, okRes := raw["result"].(map[string]interface{}); okRes {
		if msg, okMsg := res["message_id"].(float64); okMsg {
			result.MessageID = int64(msg)
		}
		if doc, okDoc := res["document"].(map[string]interface{}); okDoc {
			if fid, okFid := doc["file_id"].(string); okFid {
				result.FileID = fid
			}
		}
	}

	result.OK = true
	return nil
}
