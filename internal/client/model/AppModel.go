package model

import (
    "bufio"
    "encoding/json"
    "fmt"
    "net"
    "time"

    "fyne.io/fyne/v2/data/binding"
)

// структура модели клиентского приложения
type AppModel struct {
    MessagesBinding binding.StringList
    conn            net.Conn
}

func NewAppModel() *AppModel {
    model := &AppModel{
        MessagesBinding: binding.NewStringList(),
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

func (m *AppModel) GetMessagesBinding() binding.StringList {
    return m.MessagesBinding
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
        //m.MessagesBinding = append(m.MessagesBinding, msg)

        // TODO: сделать панель сообщения не простым тестовым полем, а чем нибудь более красивым
        displayMsg := fmt.Sprintf("%s [%s]: %s", msg.Sender, msg.Timestamp, msg.Content) // формат сообщения
        m.MessagesBinding.Append(displayMsg)
    }
}
