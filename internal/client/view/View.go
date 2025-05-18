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

    chatContainer *fyne.Container // контейнер чата, сохраняю его, чтобы не создавать постоянно
}

func NewView(app fyne.App, controller *controller.Controller) *View {
    view := &View{
        app:        app,
        controller: controller,
    }

    view.window = app.NewWindow("TkachMessenger")

    // TODO: выделить компоновку интерфейса в отдельный файл
    content := container.NewBorder(
        nil,
        nil,
        nil,
        nil,
        view.CreateChatContainer(),
    )

    view.window.SetContent(content)
    view.window.Resize(fyne.NewSize(400, 300))

    return view
}

func (v *View) ShowAndRun() {
    v.window.ShowAndRun()
}

// создать прокручиваемый список сообщений
func (v *View) CreateMessagesScroll() fyne.CanvasObject {
    messagesBinding := v.controller.Model.GetMessagesBinding()

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

    scrollContainer := container.NewScroll(messagesList)
    return scrollContainer
}

// создать контейнер ввода сообщения
func (v *View) CreateInputContainer() *fyne.Container {
    input := widget.NewEntry()
    input.SetPlaceHolder("Введите сообщение...")

    sendButton := widget.NewButton("Отправить", func() {
        v.controller.SendMessageInModel(input.Text)
        input.SetText("") // Очищаем поле после отправки
    })

    return container.NewBorder(nil, nil, nil, sendButton, input)
}

// создать контейнер чата
func (v *View) CreateChatContainer() *fyne.Container {
    return container.NewBorder(
        nil,
        v.CreateInputContainer(),
        nil,
        nil,
        v.CreateMessagesScroll(),
    )
}
