package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"freenahiFront/internal/menu"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"github.com/rs/zerolog"
)

const (
	appID = "github.soragXYZ.freenahi"
	// preferenceBackendIP = "currentBackendIP"
	preferenceLogLevel = "currentLogLevel"
)

var logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.DateTime}).With().Timestamp().Logger()

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
		logger.Trace().Msg("Started application")
	})
	app.Lifecycle().SetOnStopped(func() {
		logger.Trace().Msg("Stopped application")
	})
	app.Lifecycle().SetOnEnteredForeground(func() {
		logger.Trace().Msg("Entered foreground")
	})
	app.Lifecycle().SetOnExitedForeground(func() {
		logger.Trace().Msg("Exited foreground")
	})
}

func main() {

	fyneApp := app.NewWithID(appID)

	// Set logs level
	logLevel := fyneApp.Preferences().StringWithFallback(preferenceLogLevel, "info")
	logger.Trace().Str("logLevel", logLevel).Msg("")

	switch logLevel {
	case "trace":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "panic":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	default:
		logger.Fatal().Msgf("Unsupported value '%s' for log level. Should be trace, debug, info, warn, error, fatal or panic", logLevel)
	}

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
