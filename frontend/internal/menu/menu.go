package menu

import (
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"freenahiFront/internal/animation"
	"freenahiFront/internal/collection"
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
	}

	// TutorialIndex  defines how our tutorials should be laid out in the index tree
	TutorialIndex = map[string][]string{
		"":            {"welcome", "collections", "animations"},
		"collections": {"list", "table", "tree"},
	}
)

func MakeTopMenu(app fyne.App) *fyne.MainMenu {
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
