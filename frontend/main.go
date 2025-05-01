package main

import (
	"freenahiFront/internal/menu"
	"freenahiFront/internal/settings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
)

const (
	appID   = "github.soragXYZ.freenahi"
	appName = "Freenahi"
)

func main() {

	fyneApp := app.NewWithID(appID)

	settings.InitLogger(fyneApp)

	// Set logs level
	logLevel := fyneApp.Preferences().StringWithFallback(settings.PreferenceLogLevel, "info")
	settings.SetLogLevel(logLevel, fyneApp)

	themeValue := fyneApp.Preferences().StringWithFallback(settings.PreferenceTheme, "light")
	settings.SetTheme(themeValue, fyneApp)

	settings.LogLifecycle(fyneApp)

	w := fyneApp.NewWindow(appName)
	w.CenterOnScreen()
	w.SetFullScreen(settings.GetFullscreen(fyneApp))
	settings.MakeTray(fyneApp, w)

	w.SetMaster()
	w.Resize(fyne.NewSize(800, 800))

	w.SetMainMenu(menu.MakeTopMenu(fyneApp, w))

	content := container.NewStack()

	setTutorial := func(t menu.Tutorial) {
		content.Objects = []fyne.CanvasObject{t.View(w)}
		content.Refresh()
	}

	tutorial := container.NewBorder(nil, nil, nil, nil, content)

	// Split between the navigation menu on the left and the content of the current window on the right
	split := container.NewHSplit(menu.MakeNav(fyneApp, setTutorial, w), tutorial)
	split.Offset = 0.2
	w.SetContent(split)

	// When clicking exit on the window (reduce, fullscreen and exit icons)
	w.SetCloseIntercept(func() {
		exitOnTray := fyneApp.Preferences().BoolWithFallback(settings.PreferenceSystemTray, settings.SystemTrayDefault)
		if exitOnTray {
			w.Hide()
		} else {
			fyneApp.Quit()
		}
	})

	w.ShowAndRun()
}
