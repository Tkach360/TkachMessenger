package protocol

// структура запроса на аутентификацию
type AuthRequest struct {
	UserID   string `json:"UserID"`
	Password []byte `json:"Password"`
}
