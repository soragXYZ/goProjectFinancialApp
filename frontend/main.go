package main

import (
	"freenahiFront/internal/helper"
	"freenahiFront/internal/menu"
	"freenahiFront/internal/settings"
	"freenahiFront/internal/statusbar"
	"freenahiFront/resources"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/lang"
)

const (
	appID   = "github.soragXYZ.freenahi"
	appName = "Freenahi"
)

func main() {

	fyneApp := app.NewWithID(appID)

	helper.InitLogger(fyneApp)

	// Set logs level
	logLevel := fyneApp.Preferences().StringWithFallback(helper.PreferenceLogLevel, "info")
	helper.SetLogLevel(logLevel, fyneApp)

	// Set the language
	// This part should be moved in SetLanguage func if it is possible to change the lang without a restart
	languageIndex := resources.GetLanguageIndex(fyneApp)

	translation := resources.TranslationsInfo[languageIndex]
	translationResource, err := resources.Translations.ReadFile("translations/" + translation.Name + ".json")
	if err != nil {
		helper.Logger.Fatal().Msgf("Error loading translation file: %s", err.Error())
	}

	// "trick" Fyne into loading translations for configured language
	// by pretending it's the translation for the system locale
	name := lang.SystemLocale().LanguageString()
	lang.AddTranslations(fyne.NewStaticResource(name+".json", translationResource))
	helper.Logger.Info().Msgf("Language set to %s", translation.DisplayName)

	themeValue := fyneApp.Preferences().StringWithFallback(settings.PreferenceTheme, settings.ThemeDefault)
	settings.SetTheme(themeValue, fyneApp)

	helper.LogLifecycle(fyneApp)

	w := fyneApp.NewWindow(appName)
	w.CenterOnScreen()
	w.SetFullScreen(settings.GetFullscreen(fyneApp))
	helper.MakeTray(fyneApp, w)

	w.SetMaster()
	w.Resize(fyne.NewSize(800, 800))

	w.SetMainMenu(menu.NewTopMenu(fyneApp, w))

	content := container.NewStack()

	setTutorial := func(t menu.Tutorial) {
		content.Objects = []fyne.CanvasObject{t.View(w)}
		content.Refresh()
	}

	tutorial := container.NewBorder(nil, nil, nil, nil, content)

	// Split between the navigation menu on the left and the content of the current window on the right
	split := container.NewHSplit(menu.NewLeftMenu(fyneApp, setTutorial, w), tutorial)
	split.Offset = 0.2

	statusBar := statusbar.NewStatusBar(fyneApp, w)

	view := container.NewBorder(nil, statusBar, nil, nil, split)
	w.SetContent(view)

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
