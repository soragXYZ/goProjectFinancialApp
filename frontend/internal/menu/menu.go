package menu

import (
	"net/url"

	"fyne.io/fyne/v2"
	fyneSettings "fyne.io/fyne/v2/cmd/fyne_settings/settings"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"freenahiFront/internal/account"
	"freenahiFront/internal/animation"
	"freenahiFront/internal/collection"
	"freenahiFront/internal/settings"
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

func MakeTopMenu(app fyne.App, topWindow fyne.Window) *fyne.MainMenu {
	uiFyneSettings := func() {
		w := app.NewWindow("Fyne Settings")
		w.SetContent(fyneSettings.NewSettings().LoadAppearanceScreen(w))
		w.Resize(fyne.NewSize(440, 520))
		w.Show()
	}

	generalSettings := func() {
		win := app.NewWindow("General Settings")

		theme := settings.NewSettingItemOptions(
			"Theme",
			"Set theme color to dark or light",
			[]string{"light", "dark"},
			settings.ThemeDefault,
			func() string {
				return app.Preferences().StringWithFallback(settings.PreferenceTheme, settings.ThemeDefault)
			},
			func(v string) {
				settings.SetTheme(v, app)
			},
			win,
		)
		fullscreen := settings.NewSettingItemSwitch(
			"Fullscreen",
			"App will go fullscreen.",
			func() bool {
				return app.Preferences().BoolWithFallback(settings.PreferenceFullscreen, settings.FullscreenDefault)
			},
			func(v bool) {
				settings.SetFullScreen(v, app, topWindow, win)
			},
		)
		logLevel := settings.NewSettingItemOptions(
			"Log level",
			"Set current log level",
			settings.LogLevelNames(),
			settings.LogLevelDefault,
			func() string {
				return app.Preferences().StringWithFallback(settings.PreferenceLogLevel, settings.LogLevelDefault)
			},
			func(v string) {
				settings.SetLogLevel(v, app)
			},
			win,
		)
		backendIP := settings.NewSettingItemUserInput(
			"Backend IP",
			"Set the IP of the backend",
			"Must be IPv4. Ex: 192.168.1.1",
			`^((25[0-5]|(2[0-4]|1\d|[1-9]|)\d)\.?\b){4}$|localhost$`, // IPv4 or localhost regex
			"userEntry can only contain letters, numbers, '.', and ':'",
			settings.BackendIPDefault,
			func() string {
				return app.Preferences().StringWithFallback(settings.PreferenceBackendIP, settings.BackendIPDefault)
			},
			func(v string) {
				settings.SetBackendIP(v, app)
			},
			win,
		)
		backendPort := settings.NewSettingItemUserInput(
			"Backend Port",
			"Set the port of the backend",
			"Must be a port. Ex: 8080",
			"^([1-9][0-9]{0,3}|[1-5][0-9]{4}|6[0-4][0-9]{3}|65[0-4][0-9]{2}|655[0-2][0-9]|6553[0-5])$", // 0-65535 port regex
			"userEntry can only contain letters, numbers, '.', and ':'",
			settings.BackendPortDefault,
			func() string {
				return app.Preferences().StringWithFallback(settings.PreferenceBackendPort, settings.BackendPortDefault)
			},
			func(v string) {
				settings.SetBackendPort(v, app)
			},
			win,
		)

		items := []settings.SettingItem{
			settings.NewSettingItemHeading("Visual"),
			theme,
			fullscreen,
			settings.NewSettingItemSeperator(),
			settings.NewSettingItemHeading("Application"),
			logLevel,
			settings.NewSettingItemSeperator(),
			settings.NewSettingItemHeading("Backend"),
			backendIP,
			backendPort,
		}

		list := settings.NewSettingList(items)

		reset := settings.SettingAction{
			Label: "Reset to default",
			Action: func() {
				settings.SetTheme(settings.ThemeDefault, app)
				settings.SetFullScreen(settings.FullscreenDefault, app, topWindow, win)
				settings.SetLogLevel(settings.LogLevelDefault, app)
				settings.SetBackendIP(settings.BackendIPDefault, app)
				settings.SetBackendPort(settings.BackendPortDefault, app)
				list.Refresh()
			},
		}

		actions := []settings.SettingAction{reset}

		tabs := container.NewAppTabs(
			container.NewTabItem("General", settings.MakeSettingsPage("General", list, actions)),
			container.NewTabItem("Tab 1", widget.NewLabel("Hello")),
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
	return container.NewBorder(nil, nil, nil, nil, tree)
}
