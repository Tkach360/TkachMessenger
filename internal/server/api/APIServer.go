package api

import "github.com/Tkach360/TkachMessenger/internal/core/protocol"

type APIServer struct{}

func NewAPIServer() *APIServer {
	return &APIServer{}
}

func (a *APIServer) GetProfile(userID string) *protocol.ProfileObject {
	// пока что тестовое

	if userID == "qwer" {
		return &protocol.ProfileObject{
			UserID:   userID,
			UserName: "Qwer",
		}
	}

	if userID == "rewq" {
		return &protocol.ProfileObject{
			UserID:   userID,
			UserName: "Rewq",
		}
	}

	return &protocol.ProfileObject{
		UserID:   userID,
		UserName: "Кто ты воин",
	}
}

func (a *APIServer) GetUserPassword(userID string) ([]byte, error) {
	// пока тестовое

	if userID == "qwer" {
		return []byte("qwer"), nil
	}
	if userID == "rewq" {
		return []byte("rewq"), nil
	}

	return []byte("кто ты воин"), nil
}
