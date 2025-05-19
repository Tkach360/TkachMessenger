package model

import (
    "bufio"
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

    // подключаемся к серверу
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
    jsonMsg, _ := json.Marshal(msg)
    fmt.Fprintf(m.conn, "%s\n", jsonMsg)
    fmt.Println("Отправил", msg.Content)
}

// горутина для приема сообщений от сервера
func (m *AppModel) receiveMessages() {
    scanner := bufio.NewScanner(m.conn)
    for scanner.Scan() {
        var msg datamodel.Message
        err := json.Unmarshal(scanner.Bytes(), &msg)
        if err != nil {
            fmt.Println("Ошибка чтения сообщения:", err)
            continue
        }
        fmt.Println(msg) // лотлвыатидлтывадмто

        for i, chat := range m.Profile.Chats {
            if chat.ID == msg.ChatID {
                m.Profile.Chats[i].AddMessage(msg)
                fmt.Println("сохранил сообщение в чате")
                if msg.ChatID == m.currentChatID {
                    m.MessagesBinding.Append(msg)
                }
                break
            }
        }

        m.MessagesBinding.Append(msg)
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
            // fmt.Println("нашел чат")
            // for i := 0; i < len(chat.Messages); i++ {
            //     msg := chat.Messages[i]

            //     fmt.Println(msg.Content)

            //     m.MessagesBinding.Append(msg)
            // }
            break
        }
    }
}
