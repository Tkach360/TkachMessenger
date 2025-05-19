package datamodel

// сообщение (пока что тествое)
type Message struct {
	ID        int64  `json:"ID"`
	ChatID    string `json:"ChatID"`
	Sender    string `json:"sender"`
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"`
}
