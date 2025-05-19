package datamodel

// структура профиля
// содержит данные пользователя приложения
type Profile struct {
	UserID string `json:"UserID"`
	Chats  []Chat `json:"Chats"`
}
