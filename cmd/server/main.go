package main

import (
	//"bufio"
	"encoding/json"
	"fmt"
	"net"
	"sync"

	"github.com/Tkach360/TkachMessenger/internal/client/model/datamodel"
)

type Chat struct {
	ID           string
	Name         string // если чат личный, то содержит имя собеседника
	CountOfUsers int64  // может быть излишним
	Type         int16

	Messages []datamodel.Message
	Users    []string // список пользователей в чате
}

type Server struct {
	Clients         map[net.Conn]bool
	userConnections map[string]net.Conn
	Chats           map[string]*Chat

	//messages  []model.Message
	broadcast chan datamodel.Message
	mutex     sync.Mutex
}

func NewServer() *Server {
	return &Server{
		Clients:         make(map[net.Conn]bool),
		userConnections: make(map[string]net.Conn),
		Chats:           make(map[string]*Chat),
		//messages:  []model.Message{},
		broadcast: make(chan datamodel.Message),
		mutex:     sync.Mutex{},
	}
}

// В методе handleConnection():
func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	decoder := json.NewDecoder(conn)

	// Регистрация пользователя
	var initMsg struct{ UserID string }
	if err := decoder.Decode(&initMsg); err != nil {
		fmt.Println("Ошибка регистрации:", err)
		return
	}

	s.mutex.Lock()
	s.Clients[conn] = true
	s.userConnections[initMsg.UserID] = conn

	// Создаем/обновляем чат
	if chat, ok := s.Chats["test"]; ok {
		chat.Users = append(chat.Users, initMsg.UserID)
	} else {
		s.Chats["test"] = &Chat{
			ID:    "test",
			Users: []string{initMsg.UserID},
		}
	}
	s.mutex.Unlock()

	// Обработка сообщений
	for {
		var msg datamodel.Message
		if err := decoder.Decode(&msg); err != nil {
			fmt.Println("Клиент отключился:", initMsg.UserID)
			break
		}
		s.broadcast <- msg
	}

	// Очистка при отключении
	s.mutex.Lock()
	delete(s.Clients, conn)
	delete(s.userConnections, initMsg.UserID)
	for _, chat := range s.Chats {
		for i, userID := range chat.Users {
			if userID == initMsg.UserID {
				chat.Users = append(chat.Users[:i], chat.Users[i+1:]...)
				break
			}
		}
	}
	s.mutex.Unlock()
}

// В методе runBroadcast():
func (s *Server) runBroadcast() {
	for msg := range s.broadcast {
		s.mutex.Lock()

		// Обновляем чат
		chat, exists := s.Chats[msg.ChatID]
		if !exists {
			chat = &Chat{
				ID:    msg.ChatID,
				Users: []string{},
			}
			s.Chats[msg.ChatID] = chat
		}
		chat.Messages = append(chat.Messages, msg)

		// Отправка сообщений
		for _, userID := range chat.Users {
			if userID != msg.Sender {
				if conn, ok := s.userConnections[userID]; ok {
					if err := json.NewEncoder(conn).Encode(msg); err != nil {
						fmt.Println("Ошибка отправки:", err)
					}
				}
			}
		}
		s.mutex.Unlock()
	}
}

func main() {
	server := NewServer()

	fmt.Println(server.Clients)

	go server.runBroadcast()

	// запуск собственно TCP сервера
	ln, err := net.Listen("tcp", ":8080") // тоже в Config
	if err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
	}
	defer ln.Close()

	fmt.Println("Сервер запущен на :8080")

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Ошибка принятия соединения:", err)
			continue
		}
		fmt.Println("новое соединение")
		go server.handleConnection(conn)
	}
}

func (s *Server) GetUsersOfChat(chatID string) []string {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if chat, ok := s.Chats[chatID]; ok {
		return chat.Users
	}
	return []string{}
}
