package main

import (
	"bufio"
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
}

type Server struct {
	Clients map[net.Conn]bool
	Chats   map[string]*Chat
	//messages  []model.Message
	broadcast chan datamodel.Message
	mutex     sync.Mutex
}

func NewServer() *Server {
	return &Server{
		Clients: make(map[net.Conn]bool),
		Chats:   make(map[string]*Chat),
		//messages:  []model.Message{},
		broadcast: make(chan datamodel.Message),
		mutex:     sync.Mutex{},
	}
}

// горутина для обработки входящих сообщений
func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	// добавление подключенного клиента в список
	s.mutex.Lock()
	s.Clients[conn] = true
	s.mutex.Unlock()
	//defer delete(s.Clients, conn) // так как работаю с Mutex нужно дополнительно использовать Lock и Unlock

	// отправляем историю сообщений новому клиенту
	// вероятно пока что рано это делать
	// for _, msg := range s.messages {
	// 	jsonMsg, _ := json.Marshal(msg)
	// 	fmt.Fprintf(conn, "%s\n", jsonMsg)
	// }

	// получаем сообщения от клиента
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		var msg datamodel.Message
		err := json.Unmarshal(scanner.Bytes(), &msg)
		if err != nil {
			fmt.Println("Ошибка чтения сообщения:", err)
			continue
		}
		fmt.Println("принял ", msg)
		s.broadcast <- msg // пока что отправляем всем клиентам пока нет логики идентификации
	}

	// удаление клиента при отключении
	s.mutex.Lock()
	delete(s.Clients, conn)
	s.mutex.Unlock()
}

// рассылка всем
func (s *Server) runBroadcast() {
	for {
		msg := <-s.broadcast
		fmt.Println("Вытащил")
		s.mutex.Lock()

		fmt.Println("решаю в какой чат отправить")

		// если чат для сообщения есть, то помещаем сообщение в него
		chat, ok := s.Chats[msg.ChatID]
		fmt.Println("поискал")
		if ok {
			fmt.Println("чат есть")
			chat.Messages = append(chat.Messages, msg)
		} else {
			fmt.Println("создаю новый")
			// если такого чата нет, то создаем новый чат
			newChat := &Chat{
				ID:       msg.ChatID,
				Name:     "NewChat",
				Type:     1,
				Messages: make([]datamodel.Message, 0),
			}
			fmt.Println("вроде создал")
			newChat.Messages = append(newChat.Messages, msg)
			fmt.Println("ароде добавил")
			s.Chats[msg.ChatID] = newChat
		}

		fmt.Println("решил")

		s.mutex.Unlock()

		//s.messages = append(s.messages, msg)
		jsonMsg, _ := json.Marshal(msg)
		for client := range s.Clients {
			fmt.Fprintf(client, "%s\n", jsonMsg)
		}
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
