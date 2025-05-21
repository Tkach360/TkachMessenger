package protocol

// структура запроса на аутентификацию
type AuthRequest struct {
	UserID   int64  `json:"UserID"`
	Password []byte `json:"Password"`
}
