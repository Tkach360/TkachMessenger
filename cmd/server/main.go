package main

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"

	"github.com/Tkach360/TkachMessenger/internal/core/protocol"
	"github.com/Tkach360/TkachMessenger/internal/server/api"
	"github.com/Tkach360/TkachMessenger/pkg/tcpclient"
)

type Server struct {
	mu        sync.Mutex
	apiServer api.APIServer
	clients   map[string]net.Conn
}

func NewServer() *Server {
	return &Server{
		clients:   make(map[string]net.Conn),
		apiServer: *api.NewAPIServer(),
	}
}

func (s *Server) handleConnection(conn net.Conn) {

	// создаём TCPClient из существующего соединения
	client := tcpclient.NewTCPClientFromConn(conn)

	/*
	   создание отдельных TCPClient для каждого соединения вроде как оправдано
	   так как в другом случае, если бы был один глобальный TCPClient (тогда уж TCPServer)
	   то пришлось бы что-то делать с encoder и decoder, так как будь они общие для
	   всех соединения, то возникали бы гонки данных. К тому же это сильно упрощает логику:
	   логика работы с соединением находится в TCPClient, так что она изолирована,
	   а на сервере я могу использовать что-либо для потоков (что собственно и делаю)

	   Слабым местом является разве что TCPClient.handlers, так как для всех
	   соединений будут одни и те же обработчики

	   Возможно имеет смысл сделать отдельный TCPServer
	*/

	var userID string

	client.RegisterHandler(
		protocol.AUTH_REQUEST,
		func(obj []byte) {
			s.mu.Lock()
			defer s.mu.Unlock()

			var auth protocol.AuthRequest
			json.Unmarshal(obj, &auth)

			userID = auth.UserID
			s.clients[auth.UserID] = conn
			fmt.Printf("User %s connected\n", auth.UserID)
		})

	client.RegisterHandler(protocol.MESSAGE, s.handleMessage)
	client.RegisterHandler(protocol.CHAT_MESSAGE_REQUEST, s.handleChatMessagesRequest)

	// регистрация on-close хендлера
	client.SetOnClose(func() {
		s.mu.Lock()

		defer s.mu.Unlock()
		if userID != "" {
			s.handleDisconnect(userID)
		}
	})

	client.Listen()
}

func (s *Server) handleMessage(obj []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var msg protocol.Message
	json.Unmarshal(obj, &msg)
	s.SendToUsers(msg)
	s.apiServer.SaveMessage(msg)
}

func (s *Server) handleAuth(obj []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var auth protocol.AuthRequest
	json.Unmarshal(obj, &auth)

}

// отправка пользовательского сообщения пользователям
func (s *Server) SendToUsers(msg protocol.Message) {

	jsonMsg, _ := json.Marshal(msg)

	obj := protocol.CommunicationObject{
		Type:    protocol.MESSAGE,
		Content: jsonMsg,
	}

	jsonObj, _ := json.Marshal(obj)

	if users, err := s.apiServer.GetChatUsersID(msg.ChatID); err == nil {
		for _, user := range users {
			if user != msg.Sender {
				conn, ok := s.clients[user]
				if ok {
					fmt.Fprintf(conn, "%s\n", jsonObj)
					fmt.Println("    отправил: ", user)
				}
			}
		}
	}
}

func (s *Server) SendToUser(userID string, msg protocol.Message) {

	conn := s.clients[userID]

	jsonMsg, _ := json.Marshal(msg)

	obj := protocol.CommunicationObject{
		Type:    protocol.MESSAGE,
		Content: jsonMsg,
	}

	jsonObj, _ := json.Marshal(obj)

	fmt.Fprintf(conn, "%s\n", jsonObj)
	fmt.Println("    отправил: ", userID)

}

// обработчик закрытия соединения
func (s *Server) handleDisconnect(userID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.clients, userID)
	fmt.Printf("User %s disconnected\n", userID)
}

func (s *Server) handleChatMessagesRequest(obj []byte) {

	var cmr protocol.ChatMessagesRequest
	json.Unmarshal(obj, &cmr)

	var test protocol.Message

	msgs, _ := s.apiServer.GetAllMessages(cmr.ChatID)
	for _, msg := range msgs {
		s.SendToUser(cmr.Requester, msg)
		test = msg
	}

	fmt.Println("Обработал ChatMessagesRequest, последнее сообщение: ", test.Content)
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
