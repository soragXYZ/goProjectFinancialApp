package helper

import (
	"image/color"

	"fyne.io/fyne/v2"
)

// OnChangeFuncs is a slice of functions that can be registered
// to run when the user switches tutorial.
var OnChangeFuncs []func()

// Tutorial defines the data structure for a tutorial
type Tutorial struct {
	Title, Intro string
	View         func(w fyne.Window) fyne.CanvasObject
}

type GeneratedTutorial struct {
	title string

	content []string
	code    []func() fyne.CanvasObject
}

var (
	// Tutorials defines the metadata for each tutorial
	Tutorials = map[string]Tutorial{
		"welcome": {
			"Welcome",
			"",
			welcomeScreen,
		},
		"animations": {
			"Animations",
			"See how to animate components.",
			makeAnimationScreen,
		},
		"collections": {
			"Collections",
			"Collection widgets provide an efficient way to present lots of content.\n" +
				"The List, Table, and Tree provide a cache and re-use mechanism that make it possible to scroll through thousands of elements.\n" +
				"Use this for large data sets or for collections that can expand as users scroll.",
			collectionScreen,
		},
		"list": {
			"List",
			"A vertical arrangement of cached elements with the same styling.",
			makeListTab,
		},
		"table": {
			"Table",
			"A two dimensional cached collection of cells.",
			makeTableTab,
		},
		"tree": {
			"Tree",
			"A tree based arrangement of cached elements with the same styling.",
			makeTreeTab,
		},
	}

	// TutorialIndex  defines how our tutorials should be laid out in the index tree
	TutorialIndex = map[string][]string{
		"":            {"welcome", "collections", "animations"},
		"collections": {"list", "table", "tree"},
	}
)

type ForcedVariant struct {
	fyne.Theme

	Variant fyne.ThemeVariant
}

func (f *ForcedVariant) Color(name fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
	return f.Theme.Color(name, f.Variant)
}
