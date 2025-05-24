package api

import (
	"errors"
	"fmt"

	"github.com/Tkach360/TkachMessenger/internal/core/protocol"
)

type APIServer struct {
	data map[string][]protocol.Message
}

func NewAPIServer() *APIServer {
	return &APIServer{
		data: make(map[string][]protocol.Message),
	}
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
	if userID == "ttt" {
		return []byte("ttt"), nil
	}

	return []byte("кто ты воин"), nil
}

// получить всех пользователей чата по ID чата
func (a *APIServer) GetChatUsersID(chatID string) ([]string, error) {
	// пока что тестовое

	return []string{"qwer", "rewq", "ttt", "кто ты воин"}, nil
}

// сохранить сообщение
func (a *APIServer) SaveMessage(msg protocol.Message) error {

	if _, ok := a.data[msg.ChatID]; ok {
		a.data[msg.ChatID] = append(a.data[msg.ChatID], msg)
	} else {
		a.data[msg.ChatID] = make([]protocol.Message, 0)
		a.data[msg.ChatID] = append(a.data[msg.ChatID], msg)
	}

	fmt.Println("Сохранил сообщение, размер: ", len(a.data[msg.ChatID]))

	return nil
}

func (a *APIServer) GetAllMessages(chatID string) ([]protocol.Message, error) {
	if msgs, ok := a.data[chatID]; ok {
		return msgs, nil
	} else {
		return nil, errors.New("Нет такого чата")
	}
}

// добавить пользователя в чат
func (a *APIServer) JoinUserInChat(userID string, chatID string) error {
	// пока что тестовое
	return nil
}
