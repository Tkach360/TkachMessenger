package protocol

type ChatObject struct {
	ID           int64  `json:"ID"`
	Name         string `json:"Name"`
	CountOfUsers int64  `json:"CountOfUsers"` // может быть излишним
	Type         int8   `json:"Type"`
}
