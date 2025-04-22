package theme

import (
	"image/color"

	"fyne.io/fyne/v2"
)

type ForcedVariant struct {
	fyne.Theme

	Variant fyne.ThemeVariant
}

func (f *ForcedVariant) Color(name fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
	return f.Theme.Color(name, f.Variant)
}
