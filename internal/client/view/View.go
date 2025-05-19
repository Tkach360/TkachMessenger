package view

import (
    "fmt"

    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/data/binding"

    //"fyne.io/fyne/v2/data/binding"
    "fyne.io/fyne/v2/widget"
    "github.com/Tkach360/TkachMessenger/internal/client/controller"
    "github.com/Tkach360/TkachMessenger/internal/client/model/datamodel"

    //"github.com/Tkach360/TkachMessenger/internal/client/model"
    "github.com/Tkach360/TkachMessenger/internal/client/view/msgview"
)

type View struct {
    app        fyne.App
    window     fyne.Window
    border     *fyne.Container
    controller *controller.Controller

    chatContainer      *fyne.Container // контейнер чата, сохраняю его, чтобы не создавать постоянно
    chatsListContainer *fyne.Container // также сохраняю контейнер чатов
}

func NewView(app fyne.App, controller *controller.Controller) *View {
    view := &View{
        app:        app,
        controller: controller,
    }

    view.window = app.NewWindow("TkachMessenger")

    view.chatContainer = view.CreateChatContainer()
    view.chatsListContainer = view.CreateChatsListContainer()

    // TODO: выделить компоновку интерфейса в отдельный файл
    view.border = container.NewBorder(
        view.CreateTopPanel(),
        nil,
        nil,
        nil,
        view.chatContainer,
    )

    view.window.SetContent(view.border)
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
        msgview.GetMessageWidgetSample,
        msgview.FillMessageSample,
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

// создать прокручиваемый список чатов
func (v *View) CreateChatsListScroll() fyne.CanvasObject {
    chatsListBinding := v.controller.Model.GetChatsListBinding()

    chatsList := widget.NewListWithData(
        chatsListBinding,
        msgview.GetChatsListWidgetSample,
        msgview.FillChatsListSample,
    )

    chatsList.OnSelected = func(id widget.ListItemID) {
        item, _ := chatsListBinding.GetItem(id)
        chatUnt, _ := item.(binding.Untyped).Get()
        chat := chatUnt.(datamodel.Chat)

        v.border.Objects[0] = v.chatContainer
        v.border.Refresh()
        v.controller.OpenChat(chat.ID)
        fmt.Println("открыл чат")

        // v.border.Refresh()
    }

    scrollContainer := container.NewScroll(chatsList)
    return scrollContainer
}

// создать контейнер списка чатов
func (v *View) CreateChatsListContainer() *fyne.Container {
    return container.NewBorder(
        nil,
        nil,
        nil,
        nil,
        v.CreateChatsListScroll(),
    )
}

func (v *View) CreateTopPanel() *fyne.Container {

    backBtn := widget.NewButton("<-", func() {
        v.border.Objects[0] = v.chatsListContainer
        v.border.Refresh()
    })

    return container.NewBorder(
        nil,
        nil,
        backBtn,
        nil,
        widget.NewLabel("тут будет имя чата"),
    )
}
