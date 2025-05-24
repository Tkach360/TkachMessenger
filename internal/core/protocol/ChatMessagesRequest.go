package protocol

type ChatMessagesRequest struct {
	ChatID    string `json:"ChatID"`
	Requester string `jsin:"Requester"`
}
