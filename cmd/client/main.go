package main

import (
	"fyne.io/fyne/v2/app"
	"github.com/Tkach360/TkachMessenger/internal/client/controller"
	"github.com/Tkach360/TkachMessenger/internal/client/model"
	"github.com/Tkach360/TkachMessenger/internal/client/view"
)

// main клиента только компонует систему и запускает её
func main() {
	fyneApp := app.New()
	model := model.NewAppModel()
	controller := controller.NewController(model)
	view := view.NewView(fyneApp, controller)
	view.ShowAndRun()
}
