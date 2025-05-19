package datamodel

import "fyne.io/fyne/v2/data/binding"

// структура чата
type Chat struct {
	ID           string
	Name         string // если чат личный, то содержит имя собеседника
	CountOfUsers int64  // может быть излишним
	Type         int16

	Messages binding.UntypedList
	//Users []User // до поры до времени
}

func NewChat(ID string, Name string, CountOfUsers int64, Type int16) Chat {
	return Chat{
		ID:           ID,
		Name:         Name,
		CountOfUsers: CountOfUsers,
		Type:         Type,
		Messages:     binding.NewUntypedList(),
	}
}

func (c *Chat) AddMessage(msgs *Message) {
	c.Messages.Append(msgs)
}

func (c *Chat) GetMessages() *binding.UntypedList {
	return &c.Messages
}
