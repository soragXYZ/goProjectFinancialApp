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

func makeTopMenu(app fyne.App) *fyne.MainMenu {
	helpMenu := fyne.NewMenu("Help",
		fyne.NewMenuItem("Documentation", func() {
			u, _ := url.Parse("https://soragxyz.github.io/freenahi/")
			_ = app.OpenURL(u)
		}),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Contribute", func() {
			u, _ := url.Parse("https://soragxyz.github.io/freenahi/other/contribute/")
			_ = app.OpenURL(u)
		}),
		// a quit item will be appended to our first menu, cannot remove it
	)

	// Add new entries here if needed
	return fyne.NewMainMenu(
		helpMenu,
	)
}

func makeNav(app fyne.App, setTutorial func(tutorial helper.Tutorial), win fyne.Window) fyne.CanvasObject {

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
				setTutorial(t)
			}
		},
	}

	// Default to the welcome Menu
	tree.Select("welcome")

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

func main() {
	fyneApp := app.NewWithID(appID)
	logLifecycle(fyneApp)
	makeTray(fyneApp)

	w := fyneApp.NewWindow("Main Window")
	w.SetMaster()
	w.Resize(fyne.NewSize(800, 800))

	w.SetMainMenu(makeTopMenu(fyneApp))

	content := container.NewStack()

	setTutorial := func(t helper.Tutorial) {
		content.Objects = []fyne.CanvasObject{t.View(w)}
		content.Refresh()
	}

	tutorial := container.NewBorder(nil, nil, nil, nil, content)

	// Split
	split := container.NewHSplit(makeNav(fyneApp, setTutorial, w), tutorial)
	split.Offset = 0.2
	w.SetContent(split)

	// Exit cross on the window (with reduce and fullscreen)
	w.SetCloseIntercept(func() {
		fmt.Println("Tried to quit")
		fyneApp.Quit()
	})

	w.ShowAndRun()
}
