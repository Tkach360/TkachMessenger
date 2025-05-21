package model

import (
    //"fmt"
    "encoding/json"
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

    conn.RegisterHandler(protocol.MESSAGE, model.handleIncomingMessage)

    initMsg := protocol.Message{
        Content: model.profile.UserID,
    }

    conn.SendMessage(initMsg)

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

    if err := m.connection.SendMessage(msg); err != nil {
        return err
    }

    m.addMessageToChat(msg)
    return nil
}

// обрабатка входящиех сообщений
func (m *AppModel) handleIncomingMessage(obj json.RawMessage) {

    var msg protocol.Message
    json.Unmarshal(obj, &msg)

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
