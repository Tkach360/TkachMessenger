package main

import (
	"fyne.io/fyne/v2/app"
	"github.com/Tkach360/TkachMessenger/internal/client/controller"
	"github.com/Tkach360/TkachMessenger/internal/client/model"
	"github.com/Tkach360/TkachMessenger/internal/client/view"
	"github.com/Tkach360/TkachMessenger/pkg/tcpclient"
)

// main клиента только компонует систему и запускает её
func main() {

	tcpClient, _ := tcpclient.NewTCPClient("localhost:8080")
	go tcpClient.Listen() // не забываю запустить прослушивание

	fyneApp := app.New()

	model := model.NewAppModel(tcpClient)
	controller := controller.NewController(model)
	view := view.NewView(fyneApp, controller)
	view.ShowAndRun()
}
