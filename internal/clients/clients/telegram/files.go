package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"path"
)

func (c *Client) SendDocument(ctx context.Context, chatId string, file []byte, filename string, opts ...SendMessageOption) (*SendDocumentResponse, error) {
	var (
		result SendDocumentResponse
		b      bytes.Buffer
		q      = make(url.Values)
	)

	for _, optSetter := range opts {
		err := optSetter(q)
		if err != nil {
			return &result, fmt.Errorf("telegram client sendDocument error preparing query params: %w", err)
		}
	}

	w := multipart.NewWriter(&b)

	if err := w.WriteField("chat_id", chatId); err != nil {
		return &result, fmt.Errorf("failed to write chat_id field: %w", err)
	}

	part, err := w.CreateFormFile("document", filename)
	if err != nil {
		return &result, fmt.Errorf("failed to create form file: %w", err)
	}

	if _, err := io.Copy(part, bytes.NewReader(file)); err != nil {
		return &result, fmt.Errorf("failed to copy file content: %w", err)
	}

	if err := w.Close(); err != nil {
		return &result, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	u := url.URL{
		Scheme:   "https",
		Host:     c.host,
		Path:     path.Join(c.basePath, "sendDocument"),
		RawQuery: q.Encode(),
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), &b)
	if err != nil {
		return &result, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return &result, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return &result, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &result, fmt.Errorf("failed to read response body: %w", err)
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return &result, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}
