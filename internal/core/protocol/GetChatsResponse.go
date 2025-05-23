package protocol

type GetChatsResponse struct {
	Chats []Chat `json:"Chats"`
}

type Chat struct {
	ID           string `json:"ID"`
	Name         string `json:"Name"`
	CountOfUsers int64  // может быть излишним
	Type         int16
}
