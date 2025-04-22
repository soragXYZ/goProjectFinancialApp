package main

import (
	"fmt"
	"log"

	"freenahiFront/internal/menu"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
)

const (
	appID = "github.soragXYZ.freenahi"
)

// Set system tray if desktop (mini icon like wifi, shield, notifs, etc...)
func makeTray(app fyne.App) {
	if desk, isDesktop := app.(desktop.App); isDesktop {
		h := fyne.NewMenuItem("Come back to Freenahi", func() {})
		h.Icon = theme.HomeIcon()
		h.Action = func() { log.Println("System tray menu tapped for Welcome") }
		menu := fyne.NewMenu("SystemTrayMenu", h)
		desk.SetSystemTrayMenu(menu)
	}
}

// Watch events
func logLifecycle(app fyne.App) {
	app.Lifecycle().SetOnStarted(func() {
		log.Println("Lifecycle: Started")
	})
	app.Lifecycle().SetOnStopped(func() {
		log.Println("Lifecycle: Stopped")
	})
	app.Lifecycle().SetOnEnteredForeground(func() {
		log.Println("Lifecycle: Entered Foreground")
	})
	app.Lifecycle().SetOnExitedForeground(func() {
		log.Println("Lifecycle: Exited Foreground")
	})
}

func main() {
	fyneApp := app.NewWithID(appID)
	logLifecycle(fyneApp)
	makeTray(fyneApp)

	w := fyneApp.NewWindow("Main Window")
	w.SetMaster()
	w.Resize(fyne.NewSize(800, 800))

	w.SetMainMenu(menu.MakeTopMenu(fyneApp))

	content := container.NewStack()

	setTutorial := func(t menu.Tutorial) {
		content.Objects = []fyne.CanvasObject{t.View(w)}
		content.Refresh()
	}

	tutorial := container.NewBorder(nil, nil, nil, nil, content)

	// Split
	split := container.NewHSplit(menu.MakeNav(fyneApp, setTutorial, w), tutorial)
	split.Offset = 0.2
	w.SetContent(split)

	// Exit cross on the window (with reduce and fullscreen)
	w.SetCloseIntercept(func() {
		fmt.Println("Tried to quit")
		fyneApp.Quit()
	})

	w.ShowAndRun()
}
