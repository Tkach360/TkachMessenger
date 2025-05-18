package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
)

// нужно будет выделить в отдельный файл
type Message struct {
	Sender    string `json:"sender"`
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"`
}

type Server struct {
	Clients   map[net.Conn]bool
	messages  []Message
	broadcast chan Message
}

func NewServer() *Server {
	return &Server{
		Clients:   make(map[net.Conn]bool),
		messages:  []Message{},
		broadcast: make(chan Message),
	}
}

// горутина для обработки входящих сообщений
func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	// добавление подключенного клиента в список
	s.Clients[conn] = true
	defer delete(s.Clients, conn)

	// отправляем историю сообщений новому клиенту
	// вероятно пока что рано это делать
	for _, msg := range s.messages {
		jsonMsg, _ := json.Marshal(msg)
		fmt.Fprintf(conn, "%s\n", jsonMsg)
	}

	// получаем сообщения от клиента
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		var msg Message
		err := json.Unmarshal(scanner.Bytes(), &msg)
		if err != nil {
			fmt.Println("Ошибка чтения сообщения:", err)
			continue
		}
		s.broadcast <- msg // пока что отправляем всем клиентам пока нет логики идентификации
	}
}

// рассылка всем
func (s *Server) runBroadcast() {
	for {
		msg := <-s.broadcast
		s.messages = append(s.messages, msg)
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
