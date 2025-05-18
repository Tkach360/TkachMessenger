package msgview

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// получить шаблон сообщения
func GetMessageWidgetSample() fyne.CanvasObject {

	// TODO: при добавлении новых типов сообщений (картинки, гифки...)
	// в этот шаблон нужно поместить новые виджеты
	return container.NewHBox(
		widget.NewLabel(""),
		widget.NewLabel(""),
		widget.NewLabel(""),
	)
}
