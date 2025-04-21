package main

import (
	"fmt"
	"log"
	"net/url"

	"freenahiFront/helper"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
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

// Create the option menu
func makeMenu(fyneApp fyne.App, w fyne.Window) *fyne.MainMenu {
	helpMenu := fyne.NewMenu("Help",
		fyne.NewMenuItem("Documentation", func() {
			u, _ := url.Parse("https://soragxyz.github.io/freenahi/")
			_ = fyneApp.OpenURL(u)
		}),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Contribute", func() {
			u, _ := url.Parse("https://soragxyz.github.io/freenahi/other/contribute/")
			_ = fyneApp.OpenURL(u)
		}),
		// a quit item will be appended to our first menu, cannot remove it ?
	)

	// Add new entries here if needed
	main := fyne.NewMainMenu(
		helpMenu,
	)
	return main
}

func makeNav(setTutorial func(tutorial helper.Tutorial), loadPrevious bool, win fyne.Window) fyne.CanvasObject {
	app := fyne.CurrentApp()
	tree := &widget.Tree{
		ChildUIDs: func(uid string) []string {
			return helper.TutorialIndex[uid]
		},
		IsBranch: func(uid string) bool {
			children, ok := helper.TutorialIndex[uid]

			return ok && len(children) > 0
		},
		CreateNode: func(branch bool) fyne.CanvasObject {
			return widget.NewLabel("Collection Widgets")
		},
		UpdateNode: func(uid string, branch bool, obj fyne.CanvasObject) {
			t, ok := helper.Tutorials[uid]
			if !ok {
				fyne.LogError("Missing tutorial panel: "+uid, nil)
				return
			}
			obj.(*widget.Label).SetText(t.Title)
		},
		OnSelected: func(uid string) {
			if t, ok := helper.Tutorials[uid]; ok {
				for _, f := range helper.OnChangeFuncs {
					f()
				}
				helper.OnChangeFuncs = nil // Loading a page registers a new cleanup.

				app.Preferences().SetString(preferenceCurrentTutorial, uid)
				setTutorial(t)
			}
		},
	}

	if loadPrevious {
		currentPref := app.Preferences().StringWithFallback(preferenceCurrentTutorial, "welcome")
		tree.Select(currentPref)
	}

	themes := container.NewGridWithColumns(3,
		widget.NewButton("Dark", func() {
			app.Settings().SetTheme(&helper.ForcedVariant{Theme: theme.DefaultTheme(), Variant: theme.VariantDark})
		}),
		widget.NewButton("Light", func() {
			app.Settings().SetTheme(&helper.ForcedVariant{Theme: theme.DefaultTheme(), Variant: theme.VariantLight})
		}),
		widget.NewButton("Fullscreen", func() {
			win.SetFullScreen(!win.FullScreen())
		}),
	)

	return container.NewBorder(nil, themes, nil, nil, tree)
}

var topWindow fyne.Window

const preferenceCurrentTutorial = "welcome"

func main() {
	fyneApp := app.NewWithID(appID)
	logLifecycle(fyneApp)
	makeTray(fyneApp)

	w := fyneApp.NewWindow("Main Window")
	topWindow = w

	w.SetMainMenu(makeMenu(fyneApp, w))
	w.SetMaster()
	w.SetContent(widget.NewLabel("Hello World! Messing with the front"))
	w.Resize(fyne.NewSize(800, 800))

	content := container.NewStack()
	title := widget.NewLabel("Component name")
	intro := widget.NewLabel("An introduction would probably go\nhere, as well as a")
	top := container.NewVBox(title, widget.NewSeparator(), intro)
	setTutorial := func(t helper.Tutorial) {
		if fyne.CurrentDevice().IsMobile() {
			child := fyneApp.NewWindow(t.Title)
			topWindow = child
			child.SetContent(t.View(topWindow))
			child.Show()
			child.SetOnClosed(func() {
				topWindow = w
			})
			return
		}

		title.SetText(t.Title)
		isMarkdown := len(t.Intro) == 0
		if !isMarkdown {
			intro.SetText(t.Intro)
		}

		if t.Title == "Welcome" || isMarkdown {
			top.Hide()
		} else {
			top.Show()
		}

		content.Objects = []fyne.CanvasObject{t.View(w)}
		content.Refresh()
	}

	tutorial := container.NewBorder(nil, nil, nil, nil, content)
	if fyne.CurrentDevice().IsMobile() {
		w.SetContent(makeNav(setTutorial, false, w))
	} else {
		split := container.NewHSplit(makeNav(setTutorial, true, w), tutorial)
		split.Offset = 0.2
		w.SetContent(split)
	}

	// Exit cross on the window (with reduce and fullscreen)
	w.SetCloseIntercept(func() {
		fmt.Println("Tried to quit")
		fyneApp.Quit()
	})

	w.ShowAndRun()
}
