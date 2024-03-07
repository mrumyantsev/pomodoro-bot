package eventprocessor

type FetchProcessor interface {
	Fetcher
	Processor
}

type Fetcher interface {
	Fetch(limit int) ([]Event, error)
}

type Processor interface {
	Process(events []Event)
}

type Type int

const (
	Undefined Type = iota
	Start
	Help
	SetDefault
	SetTime
	SetTimeNotice
	Unset
	UnsetTime
	UnsetAll
)

type Event struct {
	Type   Type
	ChatId int
	Extras *Extras
}

type Extras struct {
	TimeMins int
	Notice   string
}
