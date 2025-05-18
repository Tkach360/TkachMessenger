package controller

import "github.com/Tkach360/TkachMessenger/internal/client/model"

// структура объекта контроллера
// отвечает за передачу данных из представления в модель
type Controller struct {
    Model *model.AppModel
}

func NewController(model *model.AppModel) *Controller {
    return &Controller{
        Model: model,
    }
}

// функция передачи сообщения в модель
func (c *Controller) SendMessageInModel(content string) {
    if content != "" {
        c.Model.SendMessage(content)
    }
}
