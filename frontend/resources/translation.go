package resources

import (
	"embed"
	"freenahiFront/internal/helper"
	"slices"

	"fyne.io/fyne/v2"
)

//go:embed translations
var Translations embed.FS

const (
	LanguageDefault    = "English"
	PreferenceLanguage = "currentLanguage"
)

type TranslationInfo struct {
	Name                string
	DisplayName         string
	TranslationFileName string
}

var TranslationsInfo = []TranslationInfo{
	{Name: "en", DisplayName: "English"},
	{Name: "fr", DisplayName: "Français"},
}

func GetTranslationNames() []string {
	values := []string{}
	for _, value := range TranslationsInfo {
		values = append(values, value.DisplayName)
	}

	return values
}

func SetLanguage(value string, app fyne.App) {
	app.Preferences().SetString(PreferenceLanguage, value)
	helper.Logger.Info().Msgf("Language set to %s", value)
}

// Returns the index of the value in TranslationsInfo list.
// Mainly used to get a key from a value in the TranslationsInfo "map"
// Ex: from 'English' to 'en' or from "Français" to "fr"
func GetLanguageIndex(app fyne.App) int {
	return slices.IndexFunc(TranslationsInfo, func(t TranslationInfo) bool {
		return t.DisplayName == app.Preferences().StringWithFallback(PreferenceLanguage, LanguageDefault)
	})
}
