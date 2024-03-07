package timerstore

import (
	"container/list"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/mrumyantsev/logx/log"
	"github.com/mrumyantsev/pomodoro-bot/internal/pkg/config"
	"github.com/mrumyantsev/pomodoro-bot/internal/pkg/timer-store/timer"
	botclients "github.com/mrumyantsev/pomodoro-bot/pkg/bot-clients"
)

var (
	errNoTimersToUnset = errors.New("no timers to Unset")
	errNoSuchTimer     = errors.New("no such timer")
)

type TimerStore struct {
	config *config.Config
	sender botclients.MessageSender
	list   *list.List
	mu     sync.Mutex
}

func New(
	cfg *config.Config,
	snd botclients.MessageSender,
) *TimerStore {
	ts := &TimerStore{
		config: cfg,
		sender: snd,
		list:   list.New(),
	}

	if cfg.RemoveStoppedTimersPeriodSecs > 0 {
		go ts.removeStoppedBySchedule()
	}

	return ts
}

func (t *TimerStore) Set(chatId int, timeMins int, notice string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	timer := timer.New(
		t.config,
		t.sender,
		chatId,
		timeMins,
		notice,
	)

	t.list.PushBack(timer)
}

func (t *TimerStore) Unset() (int, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.list.Len() == 0 {
		return 0, errNoTimersToUnset
	}

	el := t.list.Back()

	timer, _ := el.Value.(*timer.Timer)

	timer.Stop()

	t.list.Remove(el)

	return timer.TimeMins(), nil
}

func (t *TimerStore) UnsetTime(TimeMins int) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.list.Len() == 0 {
		return errNoTimersToUnset
	}

	for el := t.list.Front(); el != nil; el = el.Next() {
		timer, _ := el.Value.(*timer.Timer)

		if timer.TimeMins() == TimeMins {
			timer.Stop()
			t.list.Remove(el)

			return nil
		}
	}

	return errNoSuchTimer
}

func (t *TimerStore) UnsetAll() {
	t.mu.Lock()
	defer t.mu.Unlock()

	for el := t.list.Front(); el != nil; el = el.Next() {
		timer, _ := el.Value.(*timer.Timer)

		timer.Stop()
	}

	t.list.Init()
}

func (t *TimerStore) Remove(tim *timer.Timer) {
	t.mu.Lock()
	defer t.mu.Unlock()

	for el := t.list.Front(); el != nil; el = el.Next() {
		timer, _ := el.Value.(*timer.Timer)

		if timer == tim {
			t.list.Remove(el)
		}
	}
}

func (t *TimerStore) RemoveStopped() {
	t.mu.Lock()
	defer t.mu.Unlock()

	for el := t.list.Front(); el != nil; el = el.Next() {
		timer, _ := el.Value.(*timer.Timer)

		if timer.IsStopped() {
			if t.config.IsEnableDebugLogs {
				log.Debug(fmt.Sprintf("remove %d\n", timer.TimeMins()))
			}

			t.list.Remove(el)
		}
	}
}

func (t *TimerStore) removeStoppedBySchedule() {
	log.Info("timer autoremover initiated")

	for {
		log.Debug("automatic removing of stopped timers started")

		time.Sleep(
			time.Duration(t.config.RemoveStoppedTimersPeriodSecs) *
				time.Minute,
		)

		t.RemoveStopped()
	}
}
