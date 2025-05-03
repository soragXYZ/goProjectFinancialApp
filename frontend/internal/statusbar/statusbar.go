package statusbar

import (
	"fmt"
	"net/url"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"freenahiFront/internal/github"
	"freenahiFront/internal/helper"
)

const (
	downloadURL = "https://github.com/soragXYZ/freenahi/releases"
)

type StatusBar struct {
	widget.BaseWidget

	newVersionAvailable *fyne.Container
	backendStatus       *StatusBarItem
}

func NewStatusBar(app fyne.App, parentWin fyne.Window) *StatusBar {
	statusBar := &StatusBar{
		newVersionAvailable: container.NewHBox(),
	}
	statusBar.ExtendBaseWidget(statusBar)
	statusBar.startGoroutines(app, parentWin)

	statusBar.backendStatus = NewStatusBarItem(theme.NewErrorThemedResource(theme.MediaRecordIcon()), "Backend Status", func() {
		fmt.Println("tesdst")
	})

	return statusBar
}

func (a *StatusBar) CreateRenderer() fyne.WidgetRenderer {
	c := container.NewVBox(
		widget.NewSeparator(),
		container.NewHBox(
			layout.NewSpacer(),
			a.newVersionAvailable,
			widget.NewSeparator(),
			a.backendStatus,
			widget.NewSeparator(),
		))
	return widget.NewSimpleRenderer(c)
}

// StatusBarItem is a widget with a label and an optional icon, which can be tapped.
type StatusBarItem struct {
	widget.BaseWidget
	icon  *widget.Icon
	label *widget.Label

	// The function that is called when the label is tapped.
	OnTapped func()

	hovered bool
}

var _ fyne.Tappable = (*StatusBarItem)(nil)
var _ desktop.Hoverable = (*StatusBarItem)(nil)

func NewStatusBarItem(res fyne.Resource, text string, tapped func()) *StatusBarItem {
	w := &StatusBarItem{OnTapped: tapped, label: widget.NewLabel(text)}
	if res != nil {
		w.icon = widget.NewIcon(res)
	}
	w.ExtendBaseWidget(w)
	return w
}

func (w *StatusBarItem) Tapped(_ *fyne.PointEvent) {
	if w.OnTapped != nil {
		w.OnTapped()
	}
}

// Cursor returns the cursor type of this widget
func (w *StatusBarItem) Cursor() desktop.Cursor {
	if w.hovered {
		return desktop.PointerCursor
	}
	return desktop.DefaultCursor
}

// MouseIn is a hook that is called if the mouse pointer enters the element.
func (w *StatusBarItem) MouseIn(e *desktop.MouseEvent) {
	w.hovered = true
}

func (w *StatusBarItem) MouseMoved(*desktop.MouseEvent) {
	// needed to satisfy the interface only
}

// MouseOut is a hook that is called if the mouse pointer leaves the element.
func (w *StatusBarItem) MouseOut() {
	w.hovered = false
}

func (w *StatusBarItem) CreateRenderer() fyne.WidgetRenderer {
	c := container.NewHBox()
	if w.icon != nil {
		c.Add(w.icon)
	}
	c.Add(w.label)
	return widget.NewSimpleRenderer(c)
}

// Start asynchronous jobs: check if an update is available, if the backend is reachable, etc...
func (a *StatusBar) startGoroutines(app fyne.App, parentWin fyne.Window) {

	go a.showUpdateAvailable(app, parentWin)
	go a.showBackendStatus()
}

func (a *StatusBar) showBackendStatus() {
	time.Sleep(1 * time.Second)
	// a.backendStatus.label.SetText("Upated")
	a.backendStatus.icon.SetResource(theme.NewSuccessThemedResource(theme.MediaRecordIcon()))

}

func (a *StatusBar) showUpdateAvailable(app fyne.App, parentWin fyne.Window) {
	currentVersion := app.Metadata().Version
	remoteVersion, isRemoteNewer, err := github.AvailableUpdate("ErikKalkoken", "evebuddy", currentVersion)

	if err != nil {
		helper.Logger.Error().Err(err).Msg("Cannot fetch github version")
	}

	// If no update available, show nothing
	if !isRemoteNewer {
		return
	}

	// If an update is available, create a clickable hyperlink and display versions
	hyperlink := widget.NewHyperlink("Update available", nil)
	hyperlink.OnTapped = func() {
		c := container.NewVBox(
			container.NewHBox(widget.NewLabel("Latest version:"), layout.NewSpacer(), widget.NewLabel(remoteVersion)),
			container.NewHBox(widget.NewLabel("Current version:"), layout.NewSpacer(), widget.NewLabel(currentVersion)),
		)

		d := dialog.NewCustomConfirm("Update available", "Download", "Close", c, func(ok bool) {
			if !ok {
				return
			}
			download, _ := url.Parse(downloadURL)
			if err := app.OpenURL(download); err != nil {
				helper.Logger.Error().Err(err).Msg("Cannot open the URL")
			}
		}, parentWin,
		)
		d.Show()
	}
	a.newVersionAvailable.Add(widget.NewSeparator())
	a.newVersionAvailable.Add(widget.NewIcon(theme.DownloadIcon()))
	a.newVersionAvailable.Add(hyperlink)
}
