package controller

import (
    "fyne.io/fyne/v2/data/binding"
    "github.com/Tkach360/TkachMessenger/internal/client/model"
)

// Controller обрабатывает пользовательские действия и синхронизирует View с Model
type Controller struct {
    Model *model.AppModel
}

func NewController(model *model.AppModel) *Controller {
    return &Controller{
        Model: model,
    }
}

// обработчики действий
func (c *Controller) HandleSendMessage(content string) {

    c.Model.SendMessage(content)

    // if err := c.model.SendMessage(content); err != nil {
    //     c.model.NotifyError(err.Error())
    // }
}

// func (c *Controller) HandleSwitchChat(chatID string) {
//     c.Model.SwitchChat(chatID)
// }

func (c *Controller) GetMessagesBinding() binding.UntypedList {
    return c.Model.GetMessagesList()
}

func (c *Controller) GetChatsBinding() binding.UntypedList {
    return c.Model.GetChatsList()
}

func (c *Controller) OpenChat(chatID string) {
    c.Model.SwitchChat(chatID)
}

func (c *Controller) SetChatsList() {
    c.Model.CurrentChat = nil
}

func (c *Controller) GetChatName() string {
    return c.Model.CurrentChat.Name
}
