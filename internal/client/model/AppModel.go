package model

import (
    "bufio"
    "encoding/json"
    "fmt"
    "net"
    "time"
)

// сообщение (пока что тествое)
type Message struct {
    Sender    string `json:"sender"`
    Content   string `json:"content"`
    Timestamp string `json:"timestamp"`
}

// структура модели клиентского приложения
type AppModel struct {
    Messages []Message
    conn     net.Conn
}

func NewAppModel() *AppModel {
    model := &AppModel{
        Messages: []Message{},
    }

    // подключаемся к серверу
    conn, err := net.Dial("tcp", "localhost:8080")
    if err != nil {
        fmt.Println("Ошибка подключения к серверу:", err)
        return model
    }
    model.conn = conn

    go model.receiveMessages()

    return model
}

// отправить сообщение на сервер
func (m *AppModel) SendMessage(content string) {
    if m.conn == nil {
        fmt.Println("Нет подключения к серверу")
        return
    }

    msg := Message{
        Sender:    "User", // тут нужно будет заменить на реальный идентификатор
        Content:   content,
        Timestamp: time.Now().Format(time.RFC3339), // лучший ли это формат для отображения времени?
    }
    jsonMsg, _ := json.Marshal(msg)
    fmt.Fprintf(m.conn, "%s\n", jsonMsg)
}

// горутина для приема сообщений от сервера
func (m *AppModel) receiveMessages() {
    scanner := bufio.NewScanner(m.conn)
    for scanner.Scan() {
        var msg Message
        err := json.Unmarshal(scanner.Bytes(), &msg)
        if err != nil {
            fmt.Println("Ошибка чтения сообщения:", err)
            continue
        }
        fmt.Println(msg) // лотлвыатидлтывадмто
        m.Messages = append(m.Messages, msg)
    }
}
