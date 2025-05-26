package core

import "github.com/Tkach360/TkachMessenger/internal/core/protocol"

// структура чата
type Chat struct {
	ID           string
	Name         string // если чат личный, то содержит имя собеседника
	CountOfUsers int64  // может быть излишним
	Type         int16

	Messages []protocol.Message
}

func NewChat(ID string, Name string, CountOfUsers int64, Type int16) Chat {
	return Chat{
		ID:           ID,
		Name:         Name,
		CountOfUsers: CountOfUsers,
		Type:         Type,
		Messages:     make([]protocol.Message, 0),
	}
}

func (c *Chat) AddMessage(msg protocol.Message) {
	c.Messages = append(c.Messages, msg)
}

func (c *Chat) GetMessages() *[]protocol.Message {
	return &c.Messages
}
