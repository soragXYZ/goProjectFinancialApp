package menu

import (
	"fmt"
	"net/url"

	"fyne.io/fyne/v2"
	fyneSettings "fyne.io/fyne/v2/cmd/fyne_settings/settings"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"freenahiFront/internal/account"
	"freenahiFront/internal/animation"
	"freenahiFront/internal/collection"
	"freenahiFront/internal/settings"
	customTheme "freenahiFront/internal/theme"
	"freenahiFront/internal/welcome"
)

// Tutorial defines the data structure for a tutorial
type Tutorial struct {
	Title string
	View  func(w fyne.Window) fyne.CanvasObject
}

var (
	// Tutorials defines the metadata for each tutorial
	Tutorials = map[string]Tutorial{
		"welcome": {
			"Welcome",
			welcome.WelcomeScreen,
		},
		"animations": {
			"Animations",
			animation.MakeAnimationScreen,
		},
		"collections": {
			"Collections",
			collection.CollectionScreen,
		},
		"list": {
			"List",
			collection.MakeListTab,
		},
		"table": {
			"Table",
			collection.MakeTableTab,
		},
		"tree": {
			"Tree",
			collection.MakeTreeTab,
		},
		"accounts": {
			"Accounts",
			account.AccountScreen,
		},
	}

	// TutorialIndex  defines how our tutorials should be laid out in the index tree
	TutorialIndex = map[string][]string{
		"":            {"welcome", "collections", "animations", "accounts"},
		"collections": {"list", "table", "tree"},
	}
)

func MakeTopMenu(app fyne.App) *fyne.MainMenu {
	uiFyneSettings := func() {
		w := app.NewWindow("Fyne Settings")
		w.SetContent(fyneSettings.NewSettings().LoadAppearanceScreen(w))
		w.Resize(fyne.NewSize(440, 520))
		w.Show()
	}

	oula := settings.SettingAction{
		Label: "Oula",
		Action: func() {
			fmt.Println("Function oula")
		},
	}
	testBis := settings.SettingAction{
		Label: "Open new window",
		Action: func() {
			win := app.NewWindow("Test Bis Win")
			win.SetContent(widget.NewLabel("Test bis entered"))
			win.Resize(fyne.NewSize(800, 800))
			win.Show()
		},
	}
	actions := []settings.SettingAction{oula, testBis}

	generalSettings := func() {
		win := app.NewWindow("General Settings")

		tabs := container.NewAppTabs(
			container.NewTabItem("Tab 1", widget.NewLabel("Hello")),
			container.NewTabItem("General", settings.MakeSettingsPage("General", widget.NewLabel("World"), actions)),
		)

		tabs.SetTabLocation(container.TabLocationLeading)
		win.SetContent(tabs)
		win.Resize(fyne.NewSize(800, 800))
		win.Show()
	}

	helpMenu := fyne.NewMenu("Settings",
		fyne.NewMenuItem("Interface Settings", uiFyneSettings),
		fyne.NewMenuItem("General Settings", generalSettings),
		fyne.NewMenuItemSeparator(),
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

func MakeNav(app fyne.App, setTutorial func(tutorial Tutorial), win fyne.Window) fyne.CanvasObject {

	tree := &widget.Tree{
		ChildUIDs: func(uid string) []string {
			return TutorialIndex[uid]
		},
		IsBranch: func(uid string) bool {
			children, ok := TutorialIndex[uid]

			return ok && len(children) > 0
		},
		CreateNode: func(branch bool) fyne.CanvasObject {
			return widget.NewLabel("Collection Widgets")
		},
		UpdateNode: func(uid string, branch bool, obj fyne.CanvasObject) {
			t, ok := Tutorials[uid]
			if !ok {
				fyne.LogError("Missing tutorial panel: "+uid, nil)
				return
			}
			obj.(*widget.Label).SetText(t.Title)
		},
		OnSelected: func(uid string) {
			if t, ok := Tutorials[uid]; ok {
				setTutorial(t)
			}
		},
	}

	// Default to the welcome Menu
	tree.Select("welcome")

	themes := container.NewGridWithColumns(3,
		widget.NewButton("Dark", func() {
			app.Settings().SetTheme(&customTheme.ForcedVariant{Theme: theme.DefaultTheme(), Variant: theme.VariantDark})
		}),
		widget.NewButton("Light", func() {
			app.Settings().SetTheme(&customTheme.ForcedVariant{Theme: theme.DefaultTheme(), Variant: theme.VariantLight})
		}),
		widget.NewButton("Fullscreen", func() {
			win.SetFullScreen(!win.FullScreen())
		}),
	)

	return container.NewBorder(nil, themes, nil, nil, tree)
}
