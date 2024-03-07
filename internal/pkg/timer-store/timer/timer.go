package timer

import (
	"fmt"
	"time"

	"github.com/mrumyantsev/logx/log"
	"github.com/mrumyantsev/pomodoro-bot/internal/pkg/config"
	botclients "github.com/mrumyantsev/pomodoro-bot/pkg/bot-clients"
)

type Timer struct {
	config    *config.Config
	sender    botclients.MessageSender
	chatId    int
	timeMins  int
	notice    string
	stopCh    chan byte
	isStopped bool
}

func New(
	cfg *config.Config,
	snd botclients.MessageSender,
	chatId int,
	timeMins int,
	notice string,
) *Timer {
	t := &Timer{
		config:   cfg,
		sender:   snd,
		chatId:   chatId,
		timeMins: timeMins,
		notice:   notice,
		stopCh:   make(chan byte),
	}

	go t.Run()

	return t
}

func (t *Timer) IsStopped() bool {
	return t.isStopped
}

func (t *Timer) TimeMins() int {
	return t.timeMins
}

func (t *Timer) Run() {
	defer t.markAsStopped()

	ticker := time.NewTicker(1 * time.Second)
	secondsLeft := uint64(60 * t.timeMins)

	for {
		select {
		case <-t.stopCh:
			return
		case <-ticker.C:
			if t.config.IsEnableDebugLogs {
				log.Debug(fmt.Sprint("tick ", t.timeMins, secondsLeft))
			}

			if secondsLeft == 0 {
				if err := t.sender.SendMessage(t.chatId, t.notice); err != nil {
					log.Fatal("could not respond", err)
				}

				return
			}

			secondsLeft--
		}
	}
}

func (t *Timer) Stop() {
	t.stopCh <- 0
}

func (t *Timer) markAsStopped() {
	t.isStopped = true

	if t.config.IsEnableDebugLogs {
		log.Debug(fmt.Sprint("stopped ", t.timeMins))
	}
}
