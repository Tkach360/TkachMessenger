package api

import (
	"fmt"

	"github.com/Tkach360/TkachMessenger/internal/client/model/datamodel"
)

// функция для загрузки профиля
func LoadProfile() *datamodel.Profile {

	var userName string
	fmt.Scanf("%s", &userName)

	return &datamodel.Profile{
		UserID: userName,
		Chats: []datamodel.Chat{
			{
				ID:       "test",
				Name:     "NewChat",
				Type:     1,
				Messages: make([]datamodel.Message, 0),
			},
		},
	}
}
