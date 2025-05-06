package topmenu

import (
	"freenahiFront/internal/github"
	"freenahiFront/internal/helper"
	"net/url"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const (
	docURL      = "https://soragxyz.github.io/freenahi/"
	downloadURL = "https://github.com/soragXYZ/freenahi/releases"
)

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
			helper.Logger.Error().Err(err).Msg("Cannot fetch github version")
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

	doc, _ := url.Parse(docURL)
	download, _ := url.Parse(downloadURL)

	content := container.New(
		layout.NewCustomPaddedVBoxLayout(0),
		container.New(
			layout.NewCustomPaddedVBoxLayout(0),
			container.NewHBox(widget.NewLabel(lang.L("Latest version")), layout.NewSpacer(), container.NewStack(spinner, remoteItem)),
			container.NewHBox(widget.NewLabel(lang.L("Current version")), layout.NewSpacer(), currentVersionItem),
		),
		container.NewHBox(
			layout.NewSpacer(),
			widget.NewHyperlink(lang.L("Website"), doc),
			widget.NewHyperlink(lang.L("Downloads"), download),
			layout.NewSpacer(),
		),
		widget.NewSeparator(),
		container.NewHBox(
			layout.NewSpacer(),
			widget.NewLabel(lang.L("Thanks for using this application!")),
			layout.NewSpacer(),
		),
	)

	d := dialog.NewCustom(lang.L("About"), lang.L("Close"), content, win)
	d.Resize(fyne.NewSize(d.MinSize().Width*1.3, d.MinSize().Height))
	d.Show()
}

func ShowUserDataDialog(app fyne.App, win fyne.Window) {
	type item struct {
		name string
		path string
	}
	items := make([]item, 0)
	items = append(items, item{lang.L("Settings"), filepath.Join(app.Storage().RootURI().Path(), "preferences.json")})
	items = append(items, item{lang.L("Interface Settings"), filepath.Join(filepath.Dir(app.Storage().RootURI().Path()), "settings.json")})
	items = append(items, item{lang.L("Application logs"), filepath.Join(app.Storage().RootURI().Path(), app.Metadata().Name+".log")})

	form := widget.NewForm()

	for _, it := range items {
		form.Append(it.name, makePathEntry(app.Clipboard(), it.path))
	}
	d := dialog.NewCustom(lang.L("User data"), lang.L("Close"), form, win)
	d.Show()
}

func makePathEntry(cb fyne.Clipboard, path string) *fyne.Container {
	cleanedPath := filepath.Clean(path)
	return container.NewHBox(
		widget.NewLabel(cleanedPath),
		layout.NewSpacer(),
		widget.NewButtonWithIcon(lang.L("Copy"), theme.ContentCopyIcon(), func() { cb.SetContent(cleanedPath) }))
}
