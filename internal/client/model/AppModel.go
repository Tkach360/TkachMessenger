package model

import (
    //"bufio"
    "encoding/json"
    "fmt"
    "net"
    "time"

    "fyne.io/fyne/v2/data/binding"
    "github.com/Tkach360/TkachMessenger/internal/client/api"
    "github.com/Tkach360/TkachMessenger/internal/client/model/datamodel"
)

// структура модели клиентского приложения
type AppModel struct {
    MessagesBinding binding.UntypedList
    ChatsBinding    binding.UntypedList

    Profile       *datamodel.Profile
    currentChatID string

    conn net.Conn
}

func NewAppModel() *AppModel {
    model := &AppModel{
        MessagesBinding: binding.NewUntypedList(),
        ChatsBinding:    binding.NewUntypedList(),
    }

    conn, err := net.Dial("tcp", "localhost:8080")
    if err != nil {
        fmt.Println("Ошибка подключения к серверу:", err)
        return model
    }
    model.conn = conn

    model.Profile = api.LoadProfile()
    model.SetCurrentChat(model.Profile.Chats[0].ID)

    for _, chat := range model.Profile.Chats {
        model.ChatsBinding.Append(chat)
    }

    initMsg := struct {
        UserID string `json:"UserID"`
    }{
        UserID: model.Profile.UserID,
    }

    json.NewEncoder(model.conn).Encode(initMsg)

    go model.receiveMessages()

    return model
}

func (m *AppModel) GetMessagesBinding() binding.UntypedList {
    return m.MessagesBinding
}

func (m *AppModel) GetChatsListBinding() binding.UntypedList {
    return m.ChatsBinding
}

// отправить сообщение на сервер
func (m *AppModel) SendMessage(content string) {
    if m.conn == nil {
        fmt.Println("Нет подключения к серверу")
        return
    }

    msg := datamodel.Message{
        ChatID:    "test",
        Sender:    m.Profile.UserID, // тут нужно будет заменить на реальный идентификатор
        Content:   content,
        Timestamp: time.Now().Format(time.RFC3339), // лучший ли это формат для отображения времени?
    }

    json.NewEncoder(m.conn).Encode(msg)
    fmt.Println("Отправил", msg.Content)

    // сразу добавляю сообщение в чат, чтобы отображалось
    for i, chat := range m.Profile.Chats {
        if chat.ID == msg.ChatID {
            m.Profile.Chats[i].AddMessage(msg)
            fmt.Println("сохранил сообщение в чате ", msg.Content)
            if msg.ChatID == m.currentChatID {
                m.MessagesBinding.Append(msg)
            }
            break
        }
    }
}

// В методе receiveMessages():
func (m *AppModel) receiveMessages() {
    decoder := json.NewDecoder(m.conn)
    for {
        var msg datamodel.Message
        if err := decoder.Decode(&msg); err != nil {
            fmt.Println("Соединение с сервером разорвано:", err)
            m.conn.Close()
            m.conn = nil
            return // Выходим из горутины при ошибке
        }

        fmt.Println("Получено сообщение:", msg.Content)

        // Убрал дублирование - было два Append
        for i, chat := range m.Profile.Chats {
            if chat.ID == msg.ChatID {
                m.Profile.Chats[i].AddMessage(msg)
                if msg.ChatID == m.currentChatID {
                    m.MessagesBinding.Append(msg)
                }
                break
            }
        }
    }
}

func (m *AppModel) SetCurrentChat(chatID string) {
    m.currentChatID = chatID
    m.MessagesBinding.Set([]interface{}{})
    for _, chat := range m.Profile.Chats {
        if chat.ID == chatID {
            for _, msg := range chat.Messages {
                m.MessagesBinding.Append(msg)
            }
            break
        }
    }
}
