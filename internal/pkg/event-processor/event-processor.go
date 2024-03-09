package eventprocessor

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/mrumyantsev/logx/log"
	"github.com/mrumyantsev/pomodoro-bot/internal/pkg/config"
	timerstore "github.com/mrumyantsev/pomodoro-bot/internal/pkg/timer-store"
	botclients "github.com/mrumyantsev/pomodoro-bot/pkg/bot-clients"
	"github.com/mrumyantsev/pomodoro-bot/pkg/lib/e"
)

type EventProcessor struct {
	config    *config.Config
	store     *timerstore.TimerStore
	updSender botclients.UpdateSender
	offset    int
	wg        *sync.WaitGroup
}

func New(
	cfg *config.Config,
	store *timerstore.TimerStore,
	updSender botclients.UpdateSender,
) *EventProcessor {
	return &EventProcessor{
		config:    cfg,
		store:     store,
		updSender: updSender,
		wg:        &sync.WaitGroup{},
	}
}

func (p *EventProcessor) Fetch(limit int) ([]Event, error) {
	updates, err := p.updSender.Updates(p.offset, limit)
	if err != nil {
		return nil, e.Wrap("could not get events", err)
	}

	if p.config.IsEnableDebugLogs {
		log.Debug(fmt.Sprint("updates: ", updates))
	}

	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]Event, 0, len(updates))
	var ev Event

	for _, u := range updates {
		ev = event(u)
		ev.ChatId = u.Message.Chat.Id

		res = append(res, ev)
	}

	p.offset = updates[len(updates)-1].UpdateId + 1

	return res, nil
}

func (p *EventProcessor) Process(events []Event) {
	eventsCount := len(events)

	p.wg.Add(eventsCount)

	for i := 0; i < eventsCount; i++ {
		go func(i int) {
			defer p.wg.Done()

			p.doEvent(events[i])
		}(i)
	}

	p.wg.Wait()
}

func (p *EventProcessor) doEvent(event Event) {
	chatId := event.ChatId

	switch event.Type {
	case Start:
		p.mustSendMsg(chatId, "Start information here.")
	case Help:
		p.mustSendMsg(chatId, "Help information here.")
	case SetDefault:
		p.store.Set(chatId, p.config.DefaultTimeMins, p.config.DefaultNotice)
		p.mustSendMsg(chatId, p.sayTimerSet(p.config.DefaultTimeMins, p.config.DefaultNotice))
	case SetTime:
		p.store.Set(chatId, event.Extras.TimeMins, p.config.DefaultNotice)
		p.mustSendMsg(chatId, p.sayTimerSet(event.Extras.TimeMins, p.config.DefaultNotice))
	case SetTimeNotice:
		p.store.Set(chatId, event.Extras.TimeMins, event.Extras.Notice)
		p.mustSendMsg(chatId, p.sayTimerSet(event.Extras.TimeMins, event.Extras.Notice))
	case Unset:
		mins, err := p.store.Unset()
		if err != nil {
			p.mustSendMsg(chatId, e.ToPrettyString(err))
		} else {
			p.mustSendMsg(chatId, sayTimerUnset(mins))
		}
	case UnsetTime:
		if err := p.store.UnsetTime(event.Extras.TimeMins); err != nil {
			p.mustSendMsg(chatId, e.ToPrettyString(err))
		} else {
			p.mustSendMsg(chatId, sayTimerUnset(event.Extras.TimeMins))
		}
	case UnsetAll:
		p.store.UnsetAll()
		p.mustSendMsg(chatId, "All timers are unset.")
	case Undefined:
		p.mustSendMsg(chatId, "Unknown command.")
	}
}

func event(upd botclients.Update) Event {
	switch upd.Message.Text {
	case "/start":
		return newCmd(Start)
	case "/help":
		return newCmd(Help)
	default:
		return mustCmd(upd.Message.Text)
	}
}

func (p *EventProcessor) mustSendMsg(chatId int, text string) {
	if err := p.updSender.SendMessage(chatId, text); err != nil {
		log.Fatal("could not send message", err)
	}
}

func (p *EventProcessor) sayTimerSet(timeMins int, notice string) string {
	sleepDuration := time.Duration(timeMins) * time.Minute
	ringTime := time.Now().Add(sleepDuration)

	if notice != p.config.DefaultNotice {
		notice = fmt.Sprintf(`, notice: "%s"`, notice)
	} else {
		notice = ""
	}

	return fmt.Sprintf(
		"Timer is set on %d mins (ring on %v%s).",
		timeMins,
		ringTime.Format(time.TimeOnly),
		notice,
	)
}

func sayTimerUnset(timeMins int) string {
	return fmt.Sprintf("%d mins timer is unset.", timeMins)
}

func newCmd(t Type) Event {
	return Event{Type: t}
}

func mustCmd(text string) Event {
	event := mustSetCmd(text)
	if event.Type != Undefined {
		return event
	}

	event = mustUnsetCmd(text)
	if event.Type != Undefined {
		return event
	}

	return newCmd(Undefined)
}

func mustSetCmd(text string) Event {
	if len(text) < 4 {
		return newCmd(Undefined)
	}

	words := strings.Split(text, " ")

	if words[0] != "/set" {
		return newCmd(Undefined)
	}

	if len(words) == 1 {
		return newCmd(SetDefault)
	}

	timeMins, err := strconv.Atoi(words[1])
	if err != nil {
		return newCmd(Undefined)
	}

	ext := &Extras{TimeMins: timeMins}

	if len(words) == 2 {
		return newCmdWithExtras(SetTime, ext)
	}

	ext.Notice = strings.Join(words[2:], " ")

	return newCmdWithExtras(SetTimeNotice, ext)
}

func mustUnsetCmd(text string) Event {
	if len(text) < 6 {
		return newCmd(Undefined)
	}

	words := strings.Split(text, " ")

	if words[0] != "/unset" {
		return newCmd(Undefined)
	}

	if len(words) == 1 {
		return newCmd(Unset)
	}

	if len(words) > 2 {
		return newCmd(Undefined)
	}

	if strings.EqualFold(words[1], "all") {
		return newCmd(UnsetAll)
	}

	timeMins, err := strconv.Atoi(words[1])
	if err != nil {
		return newCmd(Undefined)
	}

	ext := &Extras{TimeMins: timeMins}

	return newCmdWithExtras(UnsetTime, ext)
}

func newCmdWithExtras(t Type, ext *Extras) Event {
	return Event{Type: t, Extras: ext}
}
