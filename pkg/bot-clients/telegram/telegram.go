package telegram

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"time"

	"github.com/mrumyantsev/logx/log"
	botclients "github.com/mrumyantsev/pomodoro-bot/pkg/bot-clients"
	"github.com/mrumyantsev/pomodoro-bot/pkg/lib/e"
)

const (
	telegramApiHost = "api.telegram.org"

	methodGetUpdates  = "getUpdates"
	methodSendMessage = "sendMessage"
)

type Config struct {
	BotToken               string
	InitialRetryPeriodSecs int
	RequestRetryAttempts   int
	ResponseRetryAttempts  int
}

type BotClient struct {
	config *Config
	client *http.Client
}

func New(cfg *Config) *BotClient {
	return &BotClient{
		client: &http.Client{},
		config: cfg,
	}
}

func (c *BotClient) Updates(offset int, limit int) ([]botclients.Update, error) {
	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))

	var (
		data       []byte
		attepmts   int
		periodSecs = c.config.InitialRetryPeriodSecs
		err        error
	)

	for attepmts = 1; attepmts <= c.config.RequestRetryAttempts; attepmts++ {
		if data, err = c.doRequest(methodGetUpdates, q); err == nil {
			break
		}

		log.Error("request failed, retry (attempt "+
			strconv.Itoa(attepmts)+")", nil)

		time.Sleep(time.Duration(periodSecs) * time.Second)

		periodSecs *= 2
	}
	if err != nil {
		return nil, e.Wrap("could not do request after "+
			strconv.Itoa(c.config.RequestRetryAttempts)+" retry attepmts", err)
	}

	var res botclients.UpdatesResponse

	if err = json.Unmarshal(data, &res); err != nil {
		return nil, e.Wrap("could not unmarshal request to slice of updates", err)
	}

	return res.Result, nil
}

func (c *BotClient) SendMessage(chatId int, text string) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatId))
	q.Add("text", text)

	var (
		attepmts   int
		periodSecs = c.config.InitialRetryPeriodSecs
		err        error
	)

	for attepmts = 1; attepmts <= c.config.ResponseRetryAttempts; attepmts++ {
		if _, err = c.doRequest(methodSendMessage, q); err == nil {
			break
		}

		log.Error("response failed, retry (attempt "+
			strconv.Itoa(attepmts)+")", nil)

		time.Sleep(time.Duration(periodSecs) * time.Second)

		periodSecs *= 2
	}
	if err != nil {
		return e.Wrap("could not do request after "+
			strconv.Itoa(c.config.ResponseRetryAttempts)+" retry attepmts", err)
	}

	return nil
}

func (c *BotClient) doRequest(method string, query url.Values) ([]byte, error) {
	url := url.URL{
		Scheme: "https",
		Host:   telegramApiHost,
		Path:   path.Join("bot"+c.config.BotToken, method),
	}

	req, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, e.Wrap("could not create request", err)
	}

	req.URL.RawQuery = query.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, e.Wrap("could not do request", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, e.Wrap("could not read from response body", err)
	}

	return body, nil
}
