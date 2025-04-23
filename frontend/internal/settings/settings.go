package settings

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	kxwidget "github.com/ErikKalkoken/fyne-kx/widget"
)

type settingVariant uint

const (
	settingUndefined settingVariant = iota
	settingCustom
	settingHeading
	settingSeperator
	settingSwitch
)

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

// SettingList is a custom list widget for settings.
type SettingList struct {
	widget.List

	SelectDelay time.Duration
}

func NewSettingItemOptions(
	label, hint string,
	options []string,
	defaultV string,
	getter func() string,
	setter func(v string),
	window func() fyne.Window,
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
			w := window()
			d := makeSettingDialog(
				sel,
				it.Label,
				it.Hint,
				func() {
					sel.SetSelected(defaultV)
				},
				refresh,
				w,
			)
			d.Show()
		},
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
	}
}

// NewSettingList returns a new SettingList widget.
func NewSettingList(items []SettingItem) *SettingList {
	w := &SettingList{SelectDelay: 200 * time.Millisecond}
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
		// hint := main[2].(*Label)
		// if it.Hint != "" {
		// 	hint.SetText(it.Hint)
		// 	hint.Show()
		// } else {
		// 	hint.Hide()
		// }
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
		case settingCustom:
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
		go func() {
			time.Sleep(w.SelectDelay)
			w.UnselectAll()
		}()
	}
	w.HideSeparators = true
	w.ExtendBaseWidget(w)
	return w
}

func MakeSettingsPage(title string, content fyne.CanvasObject, actions []SettingAction) fyne.CanvasObject {
	t := widget.NewLabel(title)
	t.TextStyle.Bold = true
	items := make([]*fyne.MenuItem, 0)
	for _, action := range actions {
		fmt.Println(action.Label)
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
