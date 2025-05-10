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

	w.SetContent(container.NewBorder(
		nil,
		statusbar.NewStatusBar(fyneApp, w),
		nil,
		nil,
		menu.NewLeftMenu(fyneApp, w),
	))

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

// package main

// import (
// 	"fmt"
// 	"slices"
// 	"time"

// 	"fyne.io/fyne/v2"
// 	"fyne.io/fyne/v2/app"
// 	"fyne.io/fyne/v2/widget"
// )

// var data = [][]string{[]string{"top left", "top right", "ok"},
// 	[]string{"bottom left", "bottom right", "ok"}}

// func main() {
// 	myApp := app.New()
// 	myWindow := myApp.NewWindow("Table Widget")

// 	list := widget.NewTable(
// 		func() (int, int) {
// 			return len(data), len(data[0])
// 		},
// 		func() fyne.CanvasObject {
// 			return widget.NewLabel("wide content")
// 		},
// 		func(i widget.TableCellID, o fyne.CanvasObject) {
// 			o.(*widget.Label).SetText(data[i.Row][i.Col])
// 		})

// 	list.OnSelected = func(id widget.TableCellID) {

// 		go func() {
// 			time.Sleep(1 * time.Second)
// 			fmt.Println("deleted")
// 			fmt.Println(len(data))
// 			fyne.Do(func() {
// 				list.Unselect(id)
// 				data = slices.Delete(data, id.Row, id.Row+1)
// 				list.Refresh()
// 			})
// 			fmt.Println(len(data))

// 		}()
// 	}

// 	myWindow.SetContent(list)
// 	myWindow.CenterOnScreen()
// 	myWindow.Canvas().Refresh(list)
// 	myWindow.Resize(fyne.NewSize(800, 800))
// 	myWindow.ShowAndRun()
// }
