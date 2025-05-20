package api

import (
    "fmt"
    "os"
    "path/filepath"

    "github.com/Tkach360/TkachMessenger/internal/client/core"
    "github.com/Tkach360/TkachMessenger/internal/core/protocol"
)

type APIClient struct {
    configPath string
}

func NewAPIClient() *APIClient {
    return &APIClient{
        configPath: filepath.Join(os.Getenv("HOME"), ".tkachmessenger"),
    }
}

// загрузка профиля
func (a *APIClient) LoadProfile() (*core.Profile, error) {

    // TODO: пока что это тестовое

    var userName string
    fmt.Scanf("%s", &userName)

    return &core.Profile{
        UserID: userName,
        Chats: []core.Chat{
            {
                ID:       "test",
                Name:     "NewChat",
                Type:     1,
                Messages: make([]protocol.Message, 0),
            },
        },
    }, nil
}

// TODO: сделать сохранение профиля
func (a *APIClient) SaveProfile(profile *core.Profile) error {
    return nil
}
