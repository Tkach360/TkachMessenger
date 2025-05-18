package view

import (
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/data/binding"
    "fyne.io/fyne/v2/widget"
    "github.com/Tkach360/TkachMessenger/internal/client/controller"
)

type View struct {
    app        fyne.App
    window     fyne.Window
    controller *controller.Controller
    input      *widget.Entry
}

func NewView(app fyne.App, controller *controller.Controller) *View {
    view := &View{
        app:        app,
        controller: controller,
    }

    view.window = app.NewWindow("Messenger")

    messagesBinding := controller.Model.GetMessagesBinding()

    messagesList := widget.NewListWithData(
        messagesBinding,
        func() fyne.CanvasObject {
            return widget.NewLabel("")
        },
        func(item binding.DataItem, object fyne.CanvasObject) {
            label := object.(*widget.Label)
            str, _ := item.(binding.String).Get()
            label.SetText(str)
        },
    )

    view.input = widget.NewEntry()
    view.input.SetPlaceHolder("Введите сообщение...")

    sendButton := widget.NewButton("Отправить", func() {
        view.controller.SendMessageInModel(view.input.Text)
        view.input.SetText("") // Очищаем поле после отправки
    })

    // TODO: выделить компоновку интерфейса в отдельный файл
    content := container.NewBorder(
        nil,
        container.NewVBox(view.input, sendButton), // панель снизу
        nil,
        nil,
        messagesList, // по центру располагается список сообщений
    )

    view.window.SetContent(content)
    view.window.Resize(fyne.NewSize(400, 300))

    return view
}

func (v *View) ShowAndRun() {
    v.window.ShowAndRun()
}
