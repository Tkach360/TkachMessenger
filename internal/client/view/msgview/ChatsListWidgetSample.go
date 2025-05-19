package msgview

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// получить шаблон виджета чата в списке чатов
func GetChatsListWidgetSample() fyne.CanvasObject {

	// TODO: сделать красивывм что-ли
	return container.NewHBox(
		widget.NewLabel(""),
		widget.NewLabel(""),
	)
}
