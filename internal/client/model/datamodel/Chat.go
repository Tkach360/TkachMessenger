package datamodel

// структура чата
type Chat struct {
	ID           string
	Name         string // если чат личный, то содержит имя собеседника
	CountOfUsers int64  // может быть излишним
	Type         int16

	Messages []Message
	//Users []User // до поры до времени
}

func NewChat(ID string, Name string, CountOfUsers int64, Type int16) Chat {
	return Chat{
		ID:           ID,
		Name:         Name,
		CountOfUsers: CountOfUsers,
		Type:         Type,
		Messages:     make([]Message, 0),
	}
}

func (c *Chat) AddMessage(msg Message) {
	c.Messages = append(c.Messages, msg)
}

func (c *Chat) GetMessages() *[]Message {
	return &c.Messages
}
