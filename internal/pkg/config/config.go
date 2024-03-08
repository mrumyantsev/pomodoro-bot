package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/mrumyantsev/pomodoro-bot/pkg/lib/e"
)

// A Config is the application configuration structure.
type Config struct {
	IsEnableDebugLogs bool `envconfig:"ENABLE_DEBUG_LOGS" default:"false"`

	BotToken               string `envconfig:"POM_BOT_TOKEN" default:""`
	UpdatesProcessLimit    int    `envconfig:"UPDATES_PROCESS_LIMIT" default:"15"`
	UpdatesCheckPeriodSecs int    `envconfig:"UPDATES_CHECK_PERIOD_SECS" default:"3"`
	RequestRetryAttempts   int    `envconfig:"REQUEST_RETRY_ATTEMPTS" default:"3"`
	ResponseRetryAttempts  int    `envconfig:"RESPONSE_RETRY_ATTEMPTS" default:"3"`

	RemoveStoppedTimersPeriodSecs int `envconfig:"REMOVE_STOPPED_TIMERS_PERIOD_SECS" default:"60"`

	DefaultTimeMins int    `envconfig:"DEFAULT_TIME_MINS" default:"25"`
	DefaultNotice   string `envconfig:"DEFAULT_NOTICE" default:"Pomodoro!"`
}

// New creates application configuration.
func New() *Config {
	return &Config{}
}

// Init initializes application configuration.
func (c *Config) Init() error {
	if err := envconfig.Process("", c); err != nil {
		return e.Wrap("could not populate struct with environment variables", err)
	}

	return nil
}
