package botclients

type UpdateSender interface {
	Updater
	MessageSender
}

type Updater interface {
	Updates(offset int, limit int) ([]Update, error)
}

type MessageSender interface {
	SendMessage(chatId int, text string) error
}

type UpdatesResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type Update struct {
	UpdateId int     `json:"update_id"`
	Message  Message `json:"message"`
}

type Message struct {
	Chat Chat   `json:"chat"`
	Text string `json:"text"`
}

type Chat struct {
	Id int `json:"id"`
}
