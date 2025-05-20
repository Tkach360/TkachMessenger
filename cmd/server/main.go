package main

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"

	"github.com/Tkach360/TkachMessenger/internal/core/protocol"
)

type Client struct {
	UserID string
	Conn   net.Conn
}

type Chat struct {
	ID      string
	Members map[string]net.Conn
	History []protocol.Message
}

type Server struct {
	mu      sync.Mutex
	clients map[string]net.Conn
	chats   map[string]*Chat
}

func NewServer() *Server {
	return &Server{
		clients: make(map[string]net.Conn),
		chats:   make(map[string]*Chat),
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	decoder := json.NewDecoder(conn)

	// Регистрация клиента
	var initMsg protocol.Message
	if err := decoder.Decode(&initMsg); err != nil {
		fmt.Printf("Registration error: %v", err)
		return
	}

	s.mu.Lock()
	s.clients[initMsg.Content] = conn
	s.joinChat(initMsg.Content, "test")
	s.mu.Unlock()

	fmt.Printf("User %s connected\n", initMsg.Content)

	// Обработка входящих сообщений
	for {
		var msg protocol.Message
		if err := decoder.Decode(&msg); err != nil {
			s.handleDisconnect(initMsg.Content)
			return
		}

		s.mu.Lock()
		if chat, exists := s.chats[msg.ChatID]; exists {
			fmt.Println("    принял: ", msg.Content, " от ", msg.Sender)
			chat.History = append(chat.History, msg)
			s.broadcastToChat(msg, chat)
		}
		s.mu.Unlock()
	}
}

func (s *Server) joinChat(userID, chatID string) {
	if _, exists := s.chats[chatID]; !exists {
		s.chats[chatID] = &Chat{
			ID:      chatID,
			Members: make(map[string]net.Conn),
		}
	}
	s.chats[chatID].Members[userID] = s.clients[userID]
}

func (s *Server) broadcastToChat(msg protocol.Message, chat *Chat) {
	jsonMsg, _ := json.Marshal(msg)
	for userID, conn := range chat.Members {
		if userID != msg.Sender {
			fmt.Fprintf(conn, "%s\n", jsonMsg)
			fmt.Println("    отправил: ", userID)
		}
	}
}

func (s *Server) handleDisconnect(userID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.clients, userID)
	for _, chat := range s.chats {
		delete(chat.Members, userID)
	}
	fmt.Printf("User %s disconnected\n", userID)
}

func main() {
	server := NewServer()
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Server startup error:", err)
	}
	defer listener.Close()

	fmt.Println("Server started on :8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Connection error:", err)
			continue
		}
		go server.handleConnection(conn)
	}
}

// package main

// import (
// 	//"bufio"
// 	"encoding/json"
// 	"fmt"
// 	"net"
// 	"sync"

// 	"github.com/Tkach360/TkachMessenger/internal/client/core" // пока что тестово
// )

// type Chat struct {
// 	ID           string
// 	Name         string // если чат личный, то содержит имя собеседника
// 	CountOfUsers int64  // может быть излишним
// 	Type         int16

// 	Messages []core.Message
// 	Users    []string // список пользователей в чате
// }

// type Server struct {
// 	Clients         map[net.Conn]bool
// 	userConnections map[string]net.Conn
// 	Chats           map[string]*Chat

// 	//messages  []model.Message
// 	broadcast chan core.Message
// 	mutex     sync.Mutex
// }

// func NewServer() *Server {
// 	return &Server{
// 		Clients:         make(map[net.Conn]bool),
// 		userConnections: make(map[string]net.Conn),
// 		Chats:           make(map[string]*Chat),
// 		//messages:  []model.Message{},
// 		broadcast: make(chan core.Message),
// 		mutex:     sync.Mutex{},
// 	}
// }

// // В методе handleConnection():
// func (s *Server) handleConnection(conn net.Conn) {
// 	defer conn.Close()
// 	decoder := json.NewDecoder(conn)

// 	// Регистрация пользователя
// 	var initMsg struct{ UserID string }
// 	if err := decoder.Decode(&initMsg); err != nil {
// 		fmt.Println("Ошибка регистрации:", err)
// 		return
// 	}

// 	s.mutex.Lock()
// 	s.Clients[conn] = true
// 	s.userConnections[initMsg.UserID] = conn

// 	// Создаем/обновляем чат
// 	if chat, ok := s.Chats["test"]; ok {
// 		chat.Users = append(chat.Users, initMsg.UserID)
// 	} else {
// 		s.Chats["test"] = &Chat{
// 			ID:    "test",
// 			Users: []string{initMsg.UserID},
// 		}
// 	}
// 	s.mutex.Unlock()

// 	// Обработка сообщений
// 	for {
// 		var msg core.Message
// 		if err := decoder.Decode(&msg); err != nil {
// 			fmt.Println("Клиент отключился:", initMsg.UserID)
// 			break
// 		}
// 		s.broadcast <- msg
// 	}

// 	// Очистка при отключении
// 	s.mutex.Lock()
// 	delete(s.Clients, conn)
// 	delete(s.userConnections, initMsg.UserID)
// 	for _, chat := range s.Chats {
// 		for i, userID := range chat.Users {
// 			if userID == initMsg.UserID {
// 				chat.Users = append(chat.Users[:i], chat.Users[i+1:]...)
// 				break
// 			}
// 		}
// 	}
// 	s.mutex.Unlock()
// }

// // В методе runBroadcast():
// func (s *Server) runBroadcast() {
// 	for msg := range s.broadcast {
// 		s.mutex.Lock()

// 		fmt.Println("принял сообщение: ", msg.Content)

// 		// Обновляем чат
// 		chat, exists := s.Chats[msg.ChatID]
// 		if !exists {
// 			chat = &Chat{
// 				ID:    msg.ChatID,
// 				Users: []string{},
// 			}
// 			s.Chats[msg.ChatID] = chat
// 		}
// 		chat.Messages = append(chat.Messages, msg)

// 		// Отправка сообщений
// 		for _, userID := range chat.Users {
// 			if userID != msg.Sender {
// 				if conn, ok := s.userConnections[userID]; ok {
// 					fmt.Println("   адресат: ", userID)
// 					if err := json.NewEncoder(conn).Encode(msg); err != nil {
// 						fmt.Println("Ошибка отправки:", err)
// 					}
// 					fmt.Println("отправил: ", msg.Content)
// 				}
// 			}
// 		}
// 		s.mutex.Unlock()
// 	}
// }

// func main() {
// 	server := NewServer()

// 	fmt.Println(server.Clients)

// 	go server.runBroadcast()

// 	// запуск собственно TCP сервера
// 	ln, err := net.Listen("tcp", "localhost:8080") // тоже в Config
// 	if err != nil {
// 		fmt.Println("Ошибка запуска сервера:", err)
// 	}
// 	defer ln.Close()

// 	fmt.Println("Сервер запущен на :8080")

// 	for {
// 		conn, err := ln.Accept()
// 		if err != nil {
// 			fmt.Println("Ошибка принятия соединения:", err)
// 			continue
// 		}
// 		fmt.Println("новое соединение")
// 		go server.handleConnection(conn)
// 	}
// }

// func (s *Server) GetUsersOfChat(chatID string) []string {
// 	s.mutex.Lock()
// 	defer s.mutex.Unlock()

// 	if chat, ok := s.Chats[chatID]; ok {
// 		return chat.Users
// 	}
// 	return []string{}
// }
