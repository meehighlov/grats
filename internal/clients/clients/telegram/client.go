package telegram

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/meehighlov/grats/internal/config"
)

// --------------------------------------------------------------- Types ---------------------------------------------------------------
type UpdatesChannel chan Update

type Client struct {
	host               string
	token              string
	basePath           string
	httpClient         *http.Client
	fastHttpClient     *http.Client // Для быстрых операций типа answerCallbackQuery
	longPollHttpClient *http.Client // Для long polling операций типа getUpdates
	logger             *slog.Logger
}

type SendMessageOption func(q url.Values) error

func WithParseMode(parseMode string) SendMessageOption {
	return func(q url.Values) error {
		q.Add("parse_mode", parseMode)
		return nil
	}
}

func WithMarkDown() SendMessageOption {
	return WithParseMode("MarkDown")
}

func WithDisableNotification() SendMessageOption {
	return func(q url.Values) error {
		q.Add("disable_notification", "true")
		return nil
	}
}

func WithReplyMurkup(replyMarkup []*[]map[string]interface{}) SendMessageOption {
	return func(q url.Values) error {
		mrakup_ := map[string][]*[]map[string]interface{}{}
		mrakup_["inline_keyboard"] = replyMarkup
		markup, err := json.Marshal(mrakup_)
		if err != nil {
			return nil
		}
		q.Add("reply_markup", string(markup))
		return nil
	}
}

// AutoDeleteOption - опция для автоудаления сообщения
type AutoDeleteOption struct {
	Duration time.Duration
}

func WithAutoDelete(duration time.Duration) SendMessageOption {
	return func(q url.Values) error {
		// Сохраняем данные об автоудалении в специальном поле
		q.Add("__auto_delete_duration", duration.String())
		return nil
	}
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

func New(cfg *config.Config, logger *slog.Logger) *Client {
	httpClient := &http.Client{Timeout: 10 * time.Second}
	fastHttpClient := &http.Client{Timeout: 3 * time.Second}
	longPollHttpClient := &http.Client{Timeout: 60 * time.Second}
	host := "api.telegram.org"

	return &Client{
		token:              cfg.TelegramToken,
		host:               host,
		basePath:           "bot" + cfg.TelegramToken,
		httpClient:         httpClient,
		fastHttpClient:     fastHttpClient,
		longPollHttpClient: longPollHttpClient,

		// do we need turn off logger from outside?
		logger: setupLogger(logger),
	}
}

func (tc *Client) sendRequest(ctx context.Context, method string, query url.Values) (data []byte, err error) {
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
		tc.logger.Error("error making request: " + err.Error())
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusBadRequest {
		tc.logger.Error(
			"Telegram client",
			"sendRequest", "Bad request",
			"method", method,
			"Query params was", query,
			"Response body", string(body),
		)
	}

	return body, nil
}

func (tc *Client) sendFastRequest(ctx context.Context, method string, query url.Values) (data []byte, err error) {
	defer func() { err = wrapIfErr("can't do fast request", err) }()

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

	resp, err := tc.fastHttpClient.Do(req)

	if err != nil {
		tc.logger.Error("error making fast request: " + err.Error())
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusBadRequest {
		tc.logger.Error(
			"Telegram client",
			"sendFastRequest", "Bad request",
			"method", method,
			"Query params was", query,
			"Response body", string(body),
		)
	}

	return body, nil
}

func (tc *Client) sendLongPollRequest(ctx context.Context, method string, query url.Values) (data []byte, err error) {
	defer func() { err = wrapIfErr("can't do long poll request", err) }()

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

	resp, err := tc.longPollHttpClient.Do(req)

	if err != nil {
		tc.logger.Error("error making long poll request: " + err.Error())
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusBadRequest {
		tc.logger.Error(
			"Telegram client",
			"sendLongPollRequest", "Bad request",
			"method", method,
			"Query params was", query,
			"Response body", string(body),
		)
	}

	return body, nil
}

// --------------------------------------------------------------- API methods implementation ---------------------------------------------------------------

func (tc *Client) SendMessage(ctx context.Context, chatId, text string, opts ...SendMessageOption) (*Message, error) {
	q := url.Values{}
	q.Add("chat_id", chatId)
	q.Add("text", text)

	var (
		autoDeleteDuration time.Duration
		response           SendMessageResponse
	)

	for _, optSetter := range opts {
		err := optSetter(q)
		if err != nil {
			tc.logger.Error(
				"telegram client sendMessage error preparing query params",
				"error contins", err.Error(),
			)
		}
	}

	if durationStr := q.Get("__auto_delete_duration"); durationStr != "" {
		if duration, err := time.ParseDuration(durationStr); err == nil {
			autoDeleteDuration = duration
		}
		q.Del("__auto_delete_duration")
	}

	data, err := tc.sendRequest(ctx, "sendMessage", q)
	if err != nil {
		return &response.Result, err
	}

	if err := json.Unmarshal(data, &response); err != nil {
		return &response.Result, err
	}

	if autoDeleteDuration > 0 {
		go func() {
			time.Sleep(autoDeleteDuration)
			err := tc.DeleteMessage(context.Background(), chatId, strconv.Itoa(response.Result.MessageId))
			if err != nil {
				tc.logger.Error("error auto-deleting message: " + err.Error())
			}
		}()
	}

	return &response.Result, err
}

func (tc *Client) EditMessageReplyMarkup(
	ctx context.Context,
	chatId string,
	messageId string,
	opts ...SendMessageOption,
) (*Message, error) {
	q := url.Values{}
	q.Add("chat_id", chatId)
	q.Add("message_id", messageId)

	for _, optSetter := range opts {
		err := optSetter(q)
		if err != nil {
			tc.logger.Error(
				"telegram client editMessageReplyMarkup error preparing query params",
				"error contins", err.Error(),
			)
		}
	}

	data, err := tc.sendRequest(ctx, "editMessageReplyMarkup", q)
	if err != nil {
		return nil, err
	}

	model := Message{}
	if err := json.Unmarshal(data, &model); err != nil {
		return nil, err
	}

	return &model, err
}

func (tc *Client) EditMessageText(ctx context.Context, chatId, messageId, text string, opts ...SendMessageOption) (*Message, error) {
	q := url.Values{}
	q.Add("chat_id", chatId)
	q.Add("message_id", messageId)
	q.Add("text", text)

	for _, optSetter := range opts {
		err := optSetter(q)
		if err != nil {
			tc.logger.Error(
				"telegram client editMessageText error preparing query params",
				"error contins", err.Error(),
			)
		}
	}

	data, err := tc.sendRequest(ctx, "editMessageText", q)
	if err != nil {
		tc.logger.Error("error editing message: " + err.Error())
		return nil, err
	}

	model := Message{}
	if err := json.Unmarshal(data, &model); err != nil {
		return nil, err
	}

	return &model, err
}

func (tc *Client) AnswerCallbackQuery(ctx context.Context, update *Update) error {
	queryId := update.CallbackQuery.Id
	if queryId == "" {
		return nil
	}

	ctx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()

	q := url.Values{}
	q.Add("callback_query_id", queryId)

	_, err := tc.sendFastRequest(ctx, "answerCallbackQuery", q)
	if err != nil {
		tc.logger.Error("answer callback query request error: " + err.Error())
		return err
	}

	return nil
}

func (tc *Client) GetChat(ctx context.Context, chatId string) (*GetChatResponse, error) {
	q := url.Values{}
	q.Add("chat_id", chatId)

	tc.logger.Debug("Telegram client", "chat id", chatId)

	model := GetChatResponse{}

	data, err := tc.sendRequest(ctx, "getChat", q)
	if err != nil {
		tc.logger.Error("Telegram client error", "getChat", err.Error())
		return &model, err
	}

	if err := json.Unmarshal(data, &model); err != nil {
		return &model, err
	}

	return &model, err
}

func (tc *Client) GetUpdates(ctx context.Context, offset int) (*UpdateResponse, error) {
	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(100))
	q.Add("timeout", "60")

	data, err := tc.sendLongPollRequest(ctx, "getUpdates", q)
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

func (tc *Client) GetUpdatesChannel(ctx context.Context) UpdatesChannel {
	updatesChannelSize := 100
	offset := -1

	ch := make(UpdatesChannel, updatesChannelSize)

	go func() {
		for {
			updates, err := tc.GetUpdates(ctx, offset)
			if err != nil {
				tc.logger.Error(err.Error())
				tc.logger.Error("Failed to get updates, retrying in 100ms...")
				time.Sleep(time.Millisecond * time.Duration(100))

				continue
			}

			for _, update := range updates.Result {
				if update.UpdateId >= offset {
					offset = update.UpdateId + 1
					tc.logger.Debug("GetUpdatesChannel", "update", update)
					ch <- update
				}
			}
		}
	}()

	return ch
}

func (tc *Client) DeleteMessage(ctx context.Context, chatId string, messageId string) error {
	q := url.Values{}
	q.Add("chat_id", chatId)
	q.Add("message_id", messageId)

	_, err := tc.sendRequest(ctx, "deleteMessage", q)
	if err != nil {
		tc.logger.Error("error deleting message: " + err.Error())
		return err
	}

	return nil
}

func (tc *Client) GetChatMember(ctx context.Context, userId string) (*SingleChatMemberResponse, error) {
	q := url.Values{}
	q.Add("chat_id", userId)
	q.Add("user_id", userId)

	model := SingleChatMemberResponse{}

	data, err := tc.sendRequest(ctx, "getChatMember", q)
	if err != nil {
		tc.logger.Error("error getting chat member: " + err.Error())
		return &model, err
	}

	if err := json.Unmarshal(data, &model); err != nil {
		return &model, err
	}

	return &model, nil
}
