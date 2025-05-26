package msgview

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/Tkach360/TkachMessenger/internal/client/core"
)

// функция для заполнения шаблона сообщения
func FillChatsListSample(item binding.DataItem, obj fyne.CanvasObject) {
	c, _ := item.(binding.Untyped).Get()
	chat := c.(core.Chat)

	box := obj.(*fyne.Container)
	nameLabel := box.Objects[0].(*widget.Label)

	nameLabel.SetText(chat.Name)
}
