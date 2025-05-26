package protocol

import "time"

type ContactObject struct {
	ID              int64     `json:"ID"` // идентификатор пользователя, хранимого в контакте, соответствует БД
	Name            string    `json:"Name"`
	About           string    `json:"About"` // о пользователе (может быть пустым)
	LastOnline      time.Time `json:"LastOnline"`
	ChatID          string    `json:"ChatID"`          // идентификатор личного чата с контактом (может быть пустым)
	ContactsPrivacy []bool    `json:"ContactsPrivacy"` // список ограничений для контакта :TODO

	IsOnline     bool        `json:"-"` // не сериализуется
	PersonalChat *ChatObject `json:"-"` // ссылка на чат с этим контактом внутри приложения
}
