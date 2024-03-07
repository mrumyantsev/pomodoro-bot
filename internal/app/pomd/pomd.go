package pomd

import (
	"time"

	"github.com/mrumyantsev/logx"
	"github.com/mrumyantsev/logx/log"
	"github.com/mrumyantsev/pomodoro-bot/internal/pkg/config"
	eventprocessor "github.com/mrumyantsev/pomodoro-bot/internal/pkg/event-processor"
	timerstore "github.com/mrumyantsev/pomodoro-bot/internal/pkg/timer-store"
	"github.com/mrumyantsev/pomodoro-bot/pkg/bot-clients/telegram"
)

type App struct {
	config    *config.Config
	fetchProc eventprocessor.FetchProcessor
}

func New() *App {
	cfg := config.New()

	if err := cfg.Init(); err != nil {
		log.Error("could not initialize configuration", err)
	}

	initLogger(cfg)

	botClient := initBotClient(cfg)

	store := timerstore.New(cfg, botClient)

	return &App{
		config:    cfg,
		fetchProc: eventprocessor.New(cfg, store, botClient),
	}
}

func (a *App) Run() {
	log.Info("service started")

	var (
		events []eventprocessor.Event
		err    error
	)

	for {
		time.Sleep(time.Duration(a.config.UpdatesCheckPeriodSecs) * time.Second)

		events, err = a.fetchProc.Fetch(a.config.UpdatesProcessLimit)
		if err != nil {
			log.Fatal("could not get updates", err)
		}

		if len(events) == 0 {
			continue
		}

		a.fetchProc.Process(events)
	}
}

func initLogger(cfg *config.Config) {
	logXCfg := &logx.Config{
		IsDisableDebugLogs: !cfg.IsEnableDebugLogs,
	}

	log.ApplyConfig(logXCfg)
}

func initBotClient(cfg *config.Config) *telegram.BotClient {
	tgCfg := &telegram.Config{
		BotToken:               cfg.BotToken,
		InitialRetryPeriodSecs: cfg.UpdatesCheckPeriodSecs,
		RequestRetryAttempts:   cfg.RequestRetryAttempts,
		ResponseRetryAttempts:  cfg.ResponseRetryAttempts,
	}

	return telegram.New(tgCfg)
}
