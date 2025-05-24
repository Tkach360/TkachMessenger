package model

import (
    //"fmt"
    "encoding/json"
    "fmt"
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
    CurrentChat  *core.Chat
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

    conn.RegisterHandler(protocol.MESSAGE, model.handleIncomingMessage)

    // initMsg := protocol.Message{
    //     Content: model.profile.UserID,
    // }

    // conn.SendAsCommunicationObject(protocol.MESSAGE, initMsg)

    auth := protocol.AuthRequest{
        UserID:   model.profile.UserID,
        Password: []byte(model.profile.UserID),
    }

    conn.SendAsCommunicationObject(protocol.AUTH_REQUEST, auth)

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
            m.CurrentChat = &chat
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
        ChatID:    m.CurrentChat.ID,
        Sender:    m.profile.UserID,
        Content:   content,
        Timestamp: time.Now().Format(time.RFC3339),
    }

    if err := m.connection.SendAsCommunicationObject(protocol.MESSAGE, msg); err != nil {
        return err
    }

    m.addMessageToChat(msg)
    return nil
}

// обрабатка входящиех сообщений
func (m *AppModel) handleIncomingMessage(obj []byte) {

    var msg protocol.Message
    json.Unmarshal(obj, &msg)

    // Добавляем сообщение в соответствующий чат
    for i, chat := range m.profile.Chats {
        if chat.ID == msg.ChatID {

            m.profile.Chats[i].Messages = append(m.profile.Chats[i].Messages, msg)

            if m.CurrentChat != nil {
                // если это текущий чат - обновляем интерфейс
                if msg.ChatID == m.CurrentChat.ID {
                    m.messagesList.Append(msg)
                }
            }
            break
        }
    }
}

func (m *AppModel) addMessageToChat(msg protocol.Message) {
    for i, chat := range m.profile.Chats {
        if chat.ID == msg.ChatID {
            m.profile.Chats[i].AddMessage(msg)
            if msg.ChatID == m.CurrentChat.ID {
                m.messagesList.Append(msg)
            }
            break
        }
    }
}

func (m *AppModel) SendChatMessagesRequest(chatID string) {

    m.connection.SendAsCommunicationObject(
        protocol.CHAT_MESSAGE_REQUEST,
        protocol.ChatMessagesRequest{
            ChatID:    chatID,
            Requester: m.profile.UserID,
        },
    )
}

func (m *AppModel) LoadChats() {
    for _, chat := range m.profile.Chats {
        m.SendChatMessagesRequest(chat.ID)
    }

    fmt.Println("отправил все запросы чатов")
}
