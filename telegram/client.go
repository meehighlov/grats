package telegram

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"time"

	"encoding/json"
)

// --------------------------------------------------------------- Types ---------------------------------------------------------------
type UpdatesChannel chan Update

type telegramClient struct {
	host       string
	token      string
	basePath   string
	httpClient *http.Client
	logger     *slog.Logger
}

type ApiCaller interface {
	SendMessage(context.Context, string, string) (*Message, error)
}

type apiCaller interface {
	SendMessage(context.Context, string, string) (*Message, error)
	GetUpdates(context.Context, int) (*UpdateResponse, error)
	GetUpdatesChannel(context.Context) UpdatesChannel
}

// --------------------------------------------------------------- telegram client  ---------------------------------------------------------------

func setupLogger(logger *slog.Logger) *slog.Logger {
	if logger != nil {
		return logger
	} else {
		return slog.New(
			slog.NewJSONHandler(
				os.Stdout,
				&slog.HandlerOptions{Level: slog.LevelDebug},
			),
		)
	}
}

func NewClient(token string, logger *slog.Logger) ApiCaller {
	// http client timeout > telegram getUpdates timeout
	httpClient := &http.Client{Timeout: 25 * time.Second}
	host := "api.telegram.org"

	return &telegramClient{
		token:      token,
		host:       host,
		basePath:   "bot" + token,
		httpClient: httpClient,

		// do we need turn off logger from outside?
		logger: setupLogger(logger),
	}
}

func newClient(token string, logger *slog.Logger) apiCaller {
	// http client timeout > telegram getUpdates timeout
	httpClient := &http.Client{Timeout: 20 * time.Second}
	host := "api.telegram.org"
	return &telegramClient{
		token:      token,
		host:       host,
		basePath:   "bot" + token,
		httpClient: httpClient,

		// do we need turn off logger from outside?
		logger: setupLogger(logger),
	}
}

func (tc *telegramClient) sendRequest(ctx context.Context, method string, query url.Values) (data []byte, err error) {
	defer func() { err = wrapIfErr("can't do request", err) }()

	u := url.URL{
		Scheme: "https",
		Host:   tc.host,
		Path:   path.Join(tc.basePath, method),
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = query.Encode()

	resp, err := tc.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// --------------------------------------------------------------- API methods implementation ---------------------------------------------------------------

func (tc *telegramClient) SendMessage(ctx context.Context, chatId, text string) (*Message, error) {
	q := url.Values{}
	q.Add("chat_id", chatId)
	q.Add("text", text)

	data, err := tc.sendRequest(ctx, "sendMessage", q)
	if err != nil {
		return nil, err
	}

	model := Message{}
	if err := json.Unmarshal(data, &model); err != nil {
		return nil, err
	}

	return &model, err
}

func (tc *telegramClient) GetUpdates(ctx context.Context, offset int) (*UpdateResponse, error) {
	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(100))
	q.Add("timeout", "20")

	data, err := tc.sendRequest(ctx, "getUpdates", q)
	if err != nil {
		return nil, err
	}

	model := UpdateResponse{}
	if err := json.Unmarshal(data, &model); err != nil {
		return nil, err
	}

	return &model, nil
}

// --------------------------------------------------------------- polling ---------------------------------------------------------------

func (tc *telegramClient) GetUpdatesChannel(ctx context.Context) UpdatesChannel {
	updatesChannelSize := 100
	offset := -1

	ch := make(chan Update, updatesChannelSize)

	go func() {
		for {
			updates, err := tc.GetUpdates(ctx, offset)
			if err != nil {
				tc.logger.Error(err.Error())
				tc.logger.Error("Failed to get updates, retrying in 3 seconds...")
				time.Sleep(time.Second * time.Duration(3))

				continue
			}

			for _, update := range updates.Result {
				if update.UpdateId >= offset {
					offset = update.UpdateId + 1
					ch <- update
				}
			}
		}
	}()

	return ch
}
