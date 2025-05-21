package protocol

// объект с данными профиля
type ProfileObject struct {
	UserID   string `json:"UserID"`
	UserName string `json:"UserName"`

	// при увеличении количества данных профиля будет увеличиваться
}
