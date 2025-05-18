package view

import (
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/widget"
    "github.com/Tkach360/TkachMessenger/internal/client/controller"
)

// структура объекта представления
// отвечает за весь UI и отображение приложения
type View struct {
    app        fyne.App
    window     fyne.Window
    controller *controller.Controller
    messages   *widget.Label
    input      *widget.Entry
}

func NewView(app fyne.App, controller *controller.Controller) *View {
    view := &View{
        app:        app,
        controller: controller,
    }

    // в дальнейшем нужно чтобы все подобные данные
    // были сохранены в одном объекте Config
    view.window = app.NewWindow("TkachMessenger")

    view.messages = widget.NewLabel("Сообщения будут здесь...")

    view.input = widget.NewEntry()
    view.input.SetPlaceHolder("Введите сообщение...")

    sendButton := widget.NewButton("Отправить", func() {
        view.controller.SendMessageInModel(view.input.Text)
        view.input.SetText("")
        view.updateMessages()
    })

    content := container.NewVBox(
        view.messages,
        view.input,
        sendButton,
    )

    view.window.SetContent(content)
    view.window.Resize(fyne.NewSize(400, 300))

    return view
}

// запуск приложения
func (v *View) ShowAndRun() {
    v.window.ShowAndRun()
}

// обновление отображения сообщений
func (v *View) updateMessages() {
    var text string
    for _, msg := range v.controller.Model.Messages {
        text += msg.Sender + ": " + msg.Content + "\n"
    }
    v.messages.SetText(text)
}
