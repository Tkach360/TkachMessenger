package api

import (
	"fmt"

	"fyne.io/fyne/v2/data/binding"
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
				ID:       userName,
				Name:     "NewChat",
				Type:     1,
				Messages: binding.NewUntypedList(),
			},
		},
	}
}
