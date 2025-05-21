package protocol

// ответ на аутентификацию
type AuthResponse struct {
	Status  bool          `json:"Status"`
	Profile ProfileObject `json:"Profile"`
}
