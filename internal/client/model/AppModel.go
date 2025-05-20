package model

import (
    //"fmt"
    "time"

    "fyne.io/fyne/v2/data/binding"

    "github.com/Tkach360/TkachMessenger/internal/client/api"
    "github.com/Tkach360/TkachMessenger/internal/client/core"
    "github.com/Tkach360/TkachMessenger/internal/core/protocol"
    "github.com/Tkach360/TkachMessenger/pkg/tcpclient"
)

// структура модели клиентского приложения
type AppModel struct {
    profile    *core.Profile
    connection *tcpclient.TCPClient
    apiClient  *api.APIClient

    messagesList binding.UntypedList
    chatsList    binding.UntypedList
    currentChat  *core.Chat
}

func NewAppModel(conn *tcpclient.TCPClient) *AppModel {

    apiClient := api.NewAPIClient()
    profile, _ := apiClient.LoadProfile()

    model := &AppModel{
        profile:      profile,
        connection:   conn,
        messagesList: binding.NewUntypedList(),
        chatsList:    binding.NewUntypedList(),
        apiClient:    apiClient,
    }

    conn.RegisterHandler(model.handleIncomingMessage)
    initMsg := protocol.Message{
        ID:      1, // пока тестово
        Content: model.profile.UserID,
    }

    conn.Send(initMsg)
    //json.NewEncoder(model.conn).Encode(initMsg)

    model.initChats()
    return model
}

func (m *AppModel) GetMessagesList() binding.UntypedList {
    return m.messagesList
}

func (m *AppModel) GetChatsList() binding.UntypedList {
    return m.chatsList
}

func (m *AppModel) initChats() {
    for _, chat := range m.profile.Chats {
        m.chatsList.Append(chat)
    }
}

func (m *AppModel) SwitchChat(chatID string) {

    //m.currentChat = chatID
    m.messagesList.Set([]interface{}{})
    for _, chat := range m.profile.Chats {
        if chat.ID == chatID {
            m.currentChat = &chat
            for _, msg := range chat.Messages {
                m.messagesList.Append(msg)
            }
            break
        }
    }
}

// отправить сообщение на сервер
func (m *AppModel) SendMessage(content string) error {

    msg := protocol.Message{
        ChatID:    m.currentChat.ID,
        Sender:    m.profile.UserID,
        Content:   content,
        Timestamp: time.Now().Format(time.RFC3339),
    }

    if err := m.connection.Send(msg); err != nil {
        return err
    }

    m.addMessageToChat(msg)
    return nil
}

// handleIncomingMessage обрабатывает входящие сообщения
func (m *AppModel) handleIncomingMessage(msg protocol.Message) {

    // Добавляем сообщение в соответствующий чат
    for i, chat := range m.profile.Chats {
        if chat.ID == msg.ChatID {

            m.profile.Chats[i].Messages = append(m.profile.Chats[i].Messages, msg)

            // если это текущий чат - обновляем интерфейс
            if msg.ChatID == m.currentChat.ID {
                m.messagesList.Append(msg)
            }
            break
        }
    }
}

func (m *AppModel) addMessageToChat(msg protocol.Message) {
    for i, chat := range m.profile.Chats {
        if chat.ID == msg.ChatID {
            m.profile.Chats[i].AddMessage(msg)
            if msg.ChatID == m.currentChat.ID {
                m.messagesList.Append(msg)
            }
            break
        }
    }
}

// package model

// import (
//     //"bufio"
//     "encoding/json"
//     "fmt"
//     "net"
//     "time"

//     "fyne.io/fyne/v2/data/binding"
//     "github.com/Tkach360/TkachMessenger/internal/client/api"
//     "github.com/Tkach360/TkachMessenger/internal/client/model/datamodel"
// )

// // структура модели клиентского приложения
// type AppModel struct {
//     MessagesBinding binding.UntypedList
//     ChatsBinding    binding.UntypedList

//     Profile       *datamodel.Profile
//     currentChatID string

//     conn net.Conn
// }

// func NewAppModel() *AppModel {
//     model := &AppModel{
//         MessagesBinding: binding.NewUntypedList(),
//         ChatsBinding:    binding.NewUntypedList(),
//     }

//     conn, err := net.Dial("tcp", "localhost:8080")
//     if err != nil {
//         fmt.Println("Ошибка подключения к серверу:", err)
//         return model
//     }
//     model.conn = conn

//     model.Profile = api.LoadProfile()
//     model.SetCurrentChat(model.Profile.Chats[0].ID)

//     for _, chat := range model.Profile.Chats {
//         model.ChatsBinding.Append(chat)
//     }

//     initMsg := struct {
//         UserID string `json:"UserID"`
//     }{
//         UserID: model.Profile.UserID,
//     }

//     json.NewEncoder(model.conn).Encode(initMsg)

//     go model.receiveMessages()

//     return model
// }

// func (m *AppModel) GetMessagesBinding() binding.UntypedList {
//     return m.MessagesBinding
// }

// func (m *AppModel) GetChatsListBinding() binding.UntypedList {
//     return m.ChatsBinding
// }

// // отправить сообщение на сервер
// func (m *AppModel) SendMessage(content string) {
//     if m.conn == nil {
//         fmt.Println("Нет подключения к серверу")
//         return
//     }

//     msg := datamodel.Message{
//         ChatID:    "test",
//         Sender:    m.Profile.UserID, // тут нужно будет заменить на реальный идентификатор
//         Content:   content,
//         Timestamp: time.Now().Format(time.RFC3339), // лучший ли это формат для отображения времени?
//     }

//     json.NewEncoder(m.conn).Encode(msg)
//     fmt.Println("Отправил", msg.Content)

//     // сразу добавляю сообщение в чат, чтобы отображалось
//     for i, chat := range m.Profile.Chats {
//         if chat.ID == msg.ChatID {
//             m.Profile.Chats[i].AddMessage(msg)
//             fmt.Println("сохранил сообщение в чате ", msg.Content)
//             if msg.ChatID == m.currentChatID {
//                 m.MessagesBinding.Append(msg)
//             }
//             break
//         }
//     }
// }

// // В методе receiveMessages():
// func (m *AppModel) receiveMessages() {
//     decoder := json.NewDecoder(m.conn)
//     for {
//         var msg datamodel.Message
//         if err := decoder.Decode(&msg); err != nil {
//             fmt.Println("Соединение с сервером разорвано:", err)
//             m.conn.Close()
//             m.conn = nil
//             return // Выходим из горутины при ошибке
//         }

//         fmt.Println("Получено сообщение:", msg.Content)

//         // Убрал дублирование - было два Append
//         for i, chat := range m.Profile.Chats {
//             if chat.ID == msg.ChatID {
//                 m.Profile.Chats[i].AddMessage(msg)
//                 if msg.ChatID == m.currentChatID {
//                     m.MessagesBinding.Append(msg)
//                 }
//                 break
//             }
//         }
//     }
// }

// func (m *AppModel) SetCurrentChat(chatID string) {
//     m.currentChatID = chatID
//     m.MessagesBinding.Set([]interface{}{})
//     for _, chat := range m.Profile.Chats {
//         if chat.ID == chatID {
//             for _, msg := range chat.Messages {
//                 m.MessagesBinding.Append(msg)
//             }
//             break
//         }
//     }
// }
