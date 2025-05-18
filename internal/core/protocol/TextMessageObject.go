package protocol

import "time"

// структура текстового сообщения, передаваемого по протоколу
// такую структуру можно отправлять по сети в виде json

type TextMessageObject struct {
	ID          int64     `json:"ID"`       // идентификатор сообщения, соответствующий БД
	SenderID    int64     `json:"SenderID"` // вот тут может быть нужна именно ссылка на контакт
	ReplyChat   int64     `json:"ReplyTo"`
	Text        string    `json:"Text"`
	IsRead      bool      `json:"IsRead"`
	Updated     time.Time `json:"Updated"`
	Timestamp   time.Time `json:"Timestamp"`
	SendingTime time.Time `json:"SendingTime"`
}
