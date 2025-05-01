package settings

import (
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"time"

	"freenahiFront/internal/github"
	customTheme "freenahiFront/internal/theme"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	kxwidget "github.com/ErikKalkoken/fyne-kx/widget"
	"github.com/rs/zerolog"
)

type settingVariant uint

const (
	settingUndefined settingVariant = iota
	settingText
	settingHeading
	settingSeperator
	settingSwitch
)

const (
	LogLevelDefault    = "info"
	PreferenceLogLevel = "currentLogLevel"

	BackendIPDefault    = "localhost"
	PreferenceBackendIP = "currentBackendIP"

	BackendPortDefault    = "8080"
	PreferenceBackendPort = "currentBackendPort"

	FullscreenDefault    = false
	PreferenceFullscreen = "currentFullscreen"

	SystemTrayDefault    = false
	PreferenceSystemTray = "currentSystemTray"

	ThemeDefault    = "light"
	PreferenceTheme = "currentTheme"
)

var logger zerolog.Logger

func InitLogger(app fyne.App) {
	logFilePath := filepath.Join(app.Storage().RootURI().Path(), app.Metadata().Name+".log")
	runLogFile, _ := os.OpenFile(
		logFilePath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0664,
	)
	stdOut := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.DateTime}
	multi := zerolog.MultiLevelWriter(stdOut, runLogFile)

	logger = zerolog.New(multi).With().Timestamp().Logger()
}

var logLevelName2Level = map[string]zerolog.Level{
	"trace": zerolog.TraceLevel,
	"debug": zerolog.DebugLevel,
	"info":  zerolog.InfoLevel,
	"warn":  zerolog.WarnLevel,
	"error": zerolog.ErrorLevel,
	"fatal": zerolog.FatalLevel,
	"panic": zerolog.PanicLevel,
}

type ContextMenuButton struct {
	widget.Button
	menu *fyne.Menu
}

// NewContextMenuButtonWithIcon is an icon button that shows a context menu. The label is optional.
func NewContextMenuButtonWithIcon(icon fyne.Resource, label string, menu *fyne.Menu) *ContextMenuButton {
	b := &ContextMenuButton{menu: menu}
	b.Text = label
	b.Icon = icon

	b.ExtendBaseWidget(b)
	return b
}

// Open a menu when the button is clicked
func (b *ContextMenuButton) Tapped(e *fyne.PointEvent) {
	widget.ShowPopUpMenuAtPosition(b.menu, fyne.CurrentApp().Driver().CanvasForObject(b), e.AbsolutePosition)
}

// SetMenuItems replaces the menu items.
func (b *ContextMenuButton) SetMenuItems(menuItems []*fyne.MenuItem) {
	b.menu.Items = menuItems
	b.menu.Refresh()
}

type SettingAction struct {
	Label  string
	Action func()
}

// SettingItem represents an item in a setting list.
type SettingItem struct {
	Hint   string      // optional hint text
	Label  string      // label
	Getter func() any  // returns the current value for this setting
	Setter func(v any) // sets the value for this setting

	onSelected func(it SettingItem, refresh func()) // action called when selected
	variant    settingVariant
}

func NewSettingItemOptions(
	label, hint string,
	options []string,
	defaultV string,
	getter func() string,
	setter func(v string),
	win fyne.Window,
) SettingItem {
	return SettingItem{
		Label: label,
		Hint:  hint,
		Getter: func() any {
			return getter()
		},
		Setter: func(v any) {
			setter(v.(string))
		},
		onSelected: func(it SettingItem, refresh func()) {
			sel := widget.NewRadioGroup(options, setter)
			sel.SetSelected(it.Getter().(string))
			d := makeSettingDialog(
				sel,
				it.Label,
				it.Hint,
				func() {
					sel.SetSelected(defaultV)
				},
				refresh,
				win,
			)
			d.Show()
		},
		variant: settingText,
	}
}

func makeSettingDialog(
	setting fyne.CanvasObject,
	label, hint string,
	reset, refresh func(),
	w fyne.Window,
) dialog.Dialog {
	var d dialog.Dialog
	buttons := container.NewHBox(
		widget.NewButton("Close", func() {
			d.Hide()
		}),
		layout.NewSpacer(),
		widget.NewButton("Reset", func() {
			reset()
		}),
	)
	c := container.NewBorder(
		nil,
		container.NewVBox(
			widget.NewLabel(hint),
			buttons,
		),
		nil,
		nil,
		setting,
	)

	d = dialog.NewCustomWithoutButtons(label, c, w)
	_, s := w.Canvas().InteractiveArea()
	d.Resize(fyne.NewSize(s.Width*0.8, 100))
	d.SetOnClosed(refresh)
	return d
}

func NewSettingItemHeading(label string) SettingItem {
	return SettingItem{Label: label, variant: settingHeading}
}

// NewSettingItemUserInput creates a user input setting in a setting list.
func NewSettingItemUserInput(
	label, hint, placeholder, regex, regexError string,
	defaultV string,
	getter func() string,
	setter func(v string),
	win fyne.Window,
) SettingItem {
	return SettingItem{
		Label: label,
		Hint:  hint,
		Getter: func() any {
			return getter()
		},
		Setter: func(v any) {
			setter(v.(string))
		},
		onSelected: func(it SettingItem, refresh func()) {

			userEntry := widget.NewEntry()
			userEntry.SetPlaceHolder(placeholder)
			userEntry.Validator = validation.NewRegexp(regex, regexError)

			items := []*widget.FormItem{
				widget.NewFormItem(hint, userEntry),
			}

			_, s := win.Canvas().InteractiveArea()
			d := dialog.NewForm(label, "Save", "Cancel", items, func(b bool) {
				if !b {
					return
				}
				it.Setter(userEntry.Text)
				refresh()
			}, win)
			d.Resize(fyne.NewSize(s.Width*0.7, 100))
			d.Show()
		},
		variant: settingText,
	}
}

// NewSettingItemSwitch creates a switch setting in a setting list.
func NewSettingItemSwitch(
	label, hint string,
	getter func() bool,
	onChanged func(bool),
) SettingItem {
	return SettingItem{
		Label: label,
		Hint:  hint,
		Getter: func() any {
			return getter()
		},
		Setter: func(v any) {
			onChanged(v.(bool))
		},
		onSelected: func(it SettingItem, refresh func()) {
			it.Setter(!it.Getter().(bool))
			refresh()
		},
		variant: settingSwitch,
	}
}

// NewSettingList returns a new SettingList widget.
func NewSettingList(items []SettingItem) *widget.List {
	w := &widget.List{}
	w.Length = func() int {
		return len(items)
	}
	w.CreateItem = func() fyne.CanvasObject {
		label := widget.NewLabel("Template")
		label.Truncation = fyne.TextTruncateClip
		hint := widget.NewLabel("")
		hint.Truncation = fyne.TextTruncateClip
		c := container.NewPadded(container.NewBorder(
			nil,
			container.New(layout.NewCustomPaddedLayout(0, 0, 0, 0), widget.NewSeparator()),
			nil,
			container.NewVBox(layout.NewSpacer(), container.NewStack(kxwidget.NewSwitch(nil), widget.NewLabel("")), layout.NewSpacer()),
			container.New(layout.NewCustomPaddedVBoxLayout(0), layout.NewSpacer(), label, hint, layout.NewSpacer()),
		))
		return c
	}
	w.UpdateItem = func(id widget.ListItemID, co fyne.CanvasObject) {
		if id >= len(items) {
			return
		}
		it := items[id]
		border := co.(*fyne.Container).Objects[0].(*fyne.Container).Objects
		right := border[2].(*fyne.Container).Objects[1].(*fyne.Container).Objects
		sw := right[0].(*kxwidget.Switch)
		value := right[1].(*widget.Label)
		main := border[0].(*fyne.Container).Objects
		hint := main[2].(*widget.Label)
		if it.Hint != "" {
			hint.SetText(it.Hint)
			hint.Show()
		} else {
			hint.Hide()
		}
		label := main[1].(*widget.Label)
		label.Text = it.Label
		label.TextStyle.Bold = false
		switch it.variant {
		case settingHeading:
			label.TextStyle.Bold = true
			value.Hide()
			sw.Hide()
		case settingSwitch:
			value.Hide()
			sw.OnChanged = func(v bool) {
				it.Setter(v)
			}
			sw.On = it.Getter().(bool)
			sw.Show()
			sw.Refresh()
		case settingText:
			value.SetText(fmt.Sprint(it.Getter()))
			value.Show()
			sw.Hide()
		}
		sep := border[1].(*fyne.Container)
		if it.variant == settingSeperator {
			sep.Show()
			value.Hide()
			sw.Hide()
			label.Hide()
		} else {
			sep.Hide()
			label.Show()
			label.Refresh()
		}
		w.SetItemHeight(id, co.(*fyne.Container).MinSize().Height)
	}
	w.OnSelected = func(id widget.ListItemID) {
		if id >= len(items) {
			w.UnselectAll()
			return
		}
		it := items[id]
		if it.onSelected == nil {
			w.UnselectAll()
			return
		}
		it.onSelected(it, func() {
			w.RefreshItem(id)
		})
		w.UnselectAll()
	}
	w.HideSeparators = true
	w.ExtendBaseWidget(w)
	return w
}

// NewSettingItemSeperator creates a seperator in a setting list.
func NewSettingItemSeperator() SettingItem {
	return SettingItem{variant: settingSeperator}
}

func MakeSettingsPage(title string, content fyne.CanvasObject, actions []SettingAction) fyne.CanvasObject {
	t := widget.NewLabel(title)
	t.TextStyle.Bold = true
	items := make([]*fyne.MenuItem, 0)
	for _, action := range actions {
		items = append(items, fyne.NewMenuItem(action.Label, action.Action))
	}
	options := NewContextMenuButtonWithIcon(theme.MoreHorizontalIcon(), "More", fyne.NewMenu("", items...))
	return container.NewBorder(
		container.NewVBox(container.NewHBox(t, layout.NewSpacer(), options), widget.NewSeparator()),
		nil,
		nil,
		nil,
		container.NewScroll(content),
	)
}

func LogLevelNames() []string {
	x := slices.Collect(maps.Keys(logLevelName2Level))
	return x
}

func SetLogLevel(logLevel string, app fyne.App) {
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
	app.Preferences().SetString(PreferenceLogLevel, logLevel)
	logger.Info().Msgf("Log level set to %s", logLevel)
}

func SetTheme(value string, app fyne.App) {
	switch value {
	case "light":
		app.Settings().SetTheme(&customTheme.ForcedVariant{Theme: theme.DefaultTheme(), Variant: theme.VariantLight})
	case "dark":
		app.Settings().SetTheme(&customTheme.ForcedVariant{Theme: theme.DefaultTheme(), Variant: theme.VariantDark})
	default:
		logger.Fatal().Msgf("Unsupported value '%s' for theme. Should be light or dark", value)
	}
	app.Preferences().SetString(PreferenceTheme, value)
	logger.Info().Msgf("Theme set to %s", value)
}

func SetBackendIP(value string, app fyne.App) {
	app.Preferences().SetString(PreferenceBackendIP, value)
	logger.Info().Msgf("Backend IP set to %s", value)
}

func SetBackendPort(value string, app fyne.App) {
	app.Preferences().SetString(PreferenceBackendPort, value)
	logger.Info().Msgf("Backend port set to %s", value)
}

// GetFullscreen returns the PreferenceFullscreen app preference value
func GetFullscreen(app fyne.App) bool {
	return app.Preferences().BoolWithFallback(PreferenceFullscreen, FullscreenDefault)
}

func SetFullScreen(value bool, app fyne.App, topWin fyne.Window, currentWin fyne.Window) {
	app.Preferences().SetBool(PreferenceFullscreen, value)
	topWin.SetFullScreen(app.Preferences().BoolWithFallback(PreferenceFullscreen, FullscreenDefault))
	currentWin.Show()
	logger.Info().Msgf("Fullscreen set to %s", strconv.FormatBool(value))
}

func SetSystemTray(value bool, app fyne.App) {
	app.Preferences().SetBool(PreferenceSystemTray, value)
	logger.Info().Msgf("System tray set to %s", strconv.FormatBool(value))
}

// Set system tray if desktop (mini icon like wifi, shield, notifs, etc...)
func MakeTray(app fyne.App, win fyne.Window) {
	if desk, isDesktop := app.(desktop.App); isDesktop {
		comebackItem := fyne.NewMenuItem("Open app", func() {})
		comebackItem.Icon = theme.HomeIcon()
		comebackItem.Action = func() {
			win.Show()
			logger.Trace().Msg("Going back to main menu with system tray")
		}
		appName := app.Metadata().Name
		titleItem := fyne.NewMenuItem(appName, nil)
		titleItem.Disabled = true
		menu := fyne.NewMenu("SystemTrayMenu",
			titleItem,
			fyne.NewMenuItemSeparator(),
			comebackItem,
		)
		desk.SetSystemTrayMenu(menu)
	}
}

// Watch events
func LogLifecycle(app fyne.App) {
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

func ShowUserDataDialog(app fyne.App, win fyne.Window) {
	type item struct {
		name string
		path string
	}
	items := make([]item, 0)
	items = append(items, item{"Settings", filepath.Join(app.Storage().RootURI().Path(), "preferences.json")})
	items = append(items, item{"Interface Settings", filepath.Join(filepath.Dir(app.Storage().RootURI().Path()), "settings.json")})
	items = append(items, item{"Application logs", filepath.Join(app.Storage().RootURI().Path(), app.Metadata().Name+".log")})

	form := widget.NewForm()

	for _, it := range items {
		form.Append(it.name, makePathEntry(app.Clipboard(), it.path))
	}
	d := dialog.NewCustom("User data", "Close", form, win)
	d.Show()
}

func ShowAboutDialog(app fyne.App, win fyne.Window) {

	currentVersion := app.Metadata().Version
	remoteItem := widget.NewLabel("?")
	remoteItem.Hide()
	spinner := widget.NewActivity()
	spinner.Start()

	// Stop the spinner and replace it with the remote version when obtained
	go func() {
		// ToDo: store config of owner and repo somewhere else and create a release (evebuddy used for testing)
		remoteVersion, isRemoteNewer, err := github.AvailableUpdate("ErikKalkoken", "evebuddy", currentVersion)
		if err != nil {
			logger.Error().Err(err).Msg("Cannot fetch github version")
			remoteItem.Text = "Error"
			remoteItem.Importance = widget.DangerImportance
		} else {
			remoteItem.Text = remoteVersion

			if isRemoteNewer {
				remoteItem.TextStyle.Bold = true
			}
		}

		fyne.Do(func() {
			remoteItem.Refresh()
			spinner.Hide()
			remoteItem.Show()
		})

	}()

	currentVersionItem := widget.NewLabel(currentVersion)

	content := container.New(
		layout.NewCustomPaddedVBoxLayout(0),
		container.New(
			layout.NewCustomPaddedVBoxLayout(0),
			container.NewHBox(widget.NewLabel("Latest version:"), layout.NewSpacer(), container.NewStack(spinner, remoteItem)),
			container.NewHBox(widget.NewLabel("Current version:"), layout.NewSpacer(), currentVersionItem),
		),
	)

	d := dialog.NewCustom("About", "Close", content, win)
	d.Show()
}

func makePathEntry(cb fyne.Clipboard, path string) *fyne.Container {
	cleanedPath := filepath.Clean(path)
	return container.NewHBox(
		widget.NewLabel(cleanedPath),
		layout.NewSpacer(),
		widget.NewButtonWithIcon("Copy", theme.ContentCopyIcon(), func() { cb.SetContent(cleanedPath) }))
}
