package msgview

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/Tkach360/TkachMessenger/internal/client/model"
)

// функция для заполнения шаблона сообщения
func FillMessageSample(item binding.DataItem, obj fyne.CanvasObject) {
	m, _ := item.(binding.Untyped).Get()
	msg := m.(model.Message)

	box := obj.(*fyne.Container)
	timeLabel := box.Objects[0].(*widget.Label)
	senderLabel := box.Objects[1].(*widget.Label)
	textLabel := box.Objects[2].(*widget.Label)

	timeLabel.SetText(msg.Timestamp)
	senderLabel.SetText(msg.Sender)

	// в последствие, если произойдет выделение контента сообщения
	// тут должен произойти вызов content.FillMessageSample
	// чтобы он заполнил данными виджет сообщения
	textLabel.SetText(msg.Content)
}
