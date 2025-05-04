package statusbar

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/hashicorp/go-version"

	"freenahiFront/internal/github"
	"freenahiFront/internal/helper"
	"freenahiFront/internal/settings"
)

const (
	downloadURL                = "https://github.com/soragXYZ/freenahi/releases"
	checkBackendStatusInterval = 5 * time.Second
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

	currentBackendVersionItem := widget.NewLabel("")
	currentBackendVersionItem.Hide()
	currentSpinner := widget.NewActivity()
	currentSpinner.Start()

	remoteBackendVersionItem := widget.NewLabel("")
	remoteBackendVersionItem.Hide()
	remoteSpinner := widget.NewActivity()
	remoteSpinner.Start()

	statusBar.startGoroutines(app, parentWin, remoteBackendVersionItem, currentBackendVersionItem, remoteSpinner, currentSpinner)

	statusBar.backendStatus = NewStatusBarItem(theme.NewWarningThemedResource(theme.MediaRecordIcon()), "Contacting backend...", func() {
		statusBar.showBackendDialog(parentWin, remoteBackendVersionItem, currentBackendVersionItem, remoteSpinner, currentSpinner)
	})

	return statusBar
}

func (a *StatusBar) showBackendDialog(parentWin fyne.Window, remoteBackendVersionItem, currentBackendVersionItem *widget.Label, remoteSpinner, currentSpinner *widget.Activity) {

	content := container.New(
		layout.NewCustomPaddedVBoxLayout(0),
		container.New(
			layout.NewCustomPaddedVBoxLayout(0),
			container.NewHBox(widget.NewLabel("Latest version:"), layout.NewSpacer(), container.NewStack(remoteSpinner, remoteBackendVersionItem)),
			container.NewHBox(widget.NewLabel("Current version:"), layout.NewSpacer(), container.NewStack(currentSpinner, currentBackendVersionItem)),
		),
		container.NewHBox(
			layout.NewSpacer(),
			layout.NewSpacer(),
		),
	)

	d := dialog.NewCustom("Backend version", "Close", content, parentWin)
	d.Resize(fyne.NewSize(d.MinSize().Width*1.3, d.MinSize().Height))
	d.Show()
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
func (a *StatusBar) startGoroutines(app fyne.App, parentWin fyne.Window, remoteBackendVersionItem, currentBackendVersionItem *widget.Label, remoteSpinner, currentSpiner *widget.Activity) {

	go a.showUpdateAvailable(app, parentWin)

	// Check the backend status every xxx seconds
	go func() {
		for {
			backendIp := app.Preferences().StringWithFallback(settings.PreferenceBackendIP, settings.BackendIPDefault)
			backendProtocol := app.Preferences().StringWithFallback(settings.PreferenceBackendProtocol, settings.BackendProtocolDefault)
			backendPort := app.Preferences().StringWithFallback(settings.PreferenceBackendPort, settings.BackendPortDefault)

			var statusBarText string
			statusBarSuccessIcon := theme.NewSuccessThemedResource(theme.MediaRecordIcon())
			statusBarWarnIcon := theme.NewWarningThemedResource(theme.MediaRecordIcon())
			statusBarErrorIcon := theme.NewErrorThemedResource(theme.MediaRecordIcon())

			// Get remote backend version
			// ToDo: get the actual version of the backend when the image is finalized
			remoteBackendVersion, err := version.NewVersion("0.0.2")
			if err != nil {
				helper.Logger.Error().Err(err).Msg("Version error")
			}

			helper.Logger.Trace().Str("Remote backend version", remoteBackendVersion.String()).Msg("")
			remoteBackendVersionItem.Text = remoteBackendVersion.String()

			// Get current backend version
			url := fmt.Sprintf("%s://%s:%s/version/", backendProtocol, backendIp, backendPort)
			resp, err := http.Get(url)

			if e, ok := err.(net.Error); ok && e.Timeout() { // Backend unreachable
				helper.Logger.Error().Err(err).Msg("Timeout")
				statusBarText = "Backend unreachable"
				fyne.Do(func() { a.backendStatus.icon.SetResource(statusBarErrorIcon) })

			} else if err != nil { // Backend Error
				helper.Logger.Error().Err(err).Msg("Cannot run http get request")
				statusBarText = "Backend Error"
				fyne.Do(func() { a.backendStatus.icon.SetResource(statusBarErrorIcon) })

			} else { // Reachable

				body, err := io.ReadAll(resp.Body)
				if err != nil {
					helper.Logger.Error().Err(err).Msg("ReadAll error")
				}

				currentBackendVersion, err := version.NewVersion(string(body))
				if err != nil {
					helper.Logger.Error().Err(err).Msg("Version error")
				}

				helper.Logger.Trace().Str("current backend version", currentBackendVersion.String()).Msg("")
				currentBackendVersionItem.Text = currentBackendVersion.String()

				if currentBackendVersion.LessThan(remoteBackendVersion) {
					statusBarText = "Backend update available"
					fyne.Do(func() { a.backendStatus.icon.SetResource(statusBarWarnIcon) })
				} else {
					statusBarText = "Backend OK"
					fyne.Do(func() { a.backendStatus.icon.SetResource(statusBarSuccessIcon) })
				}
			}

			if len(remoteBackendVersionItem.Text) != 0 {
				fyne.Do(func() {
					remoteBackendVersionItem.Refresh()
					remoteSpinner.Hide()
					remoteBackendVersionItem.Show()
				})
			}

			if len(currentBackendVersionItem.Text) != 0 {
				fyne.Do(func() {
					currentBackendVersionItem.Refresh()
					currentSpiner.Hide()
					currentBackendVersionItem.Show()
				})
			}

			fyne.Do(func() { a.backendStatus.label.SetText(statusBarText) })

			time.Sleep(checkBackendStatusInterval)

		}

	}()

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
	hyperlink := widget.NewHyperlink("Application update available", nil)
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
	fyne.Do(func() {
		a.newVersionAvailable.Add(widget.NewSeparator())
		a.newVersionAvailable.Add(widget.NewIcon(theme.DownloadIcon()))
		a.newVersionAvailable.Add(hyperlink)
	})

}
