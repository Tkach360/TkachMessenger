package core

// структура профиля
// содержит данные пользователя приложения
type Profile struct {
	UserID   string `json:"UserID"`
	UserName string `json:"UserName"`
	Chats    []Chat `json:"Chats"`
}
