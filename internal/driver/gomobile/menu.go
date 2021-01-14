package gomobile

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/container"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

type menuLabel struct {
	widget.BaseWidget

	menu   *fyne.Menu
	bar    *fyne.Container
	canvas *mobileCanvas
}

func (m *menuLabel) Tapped(*fyne.PointEvent) {
	pos := fyne.CurrentApp().Driver().AbsolutePositionForObject(m)
	menu := widget.NewPopUpMenu(m.menu, m.canvas)
	menu.ShowAtPosition(fyne.NewPos(pos.X+m.Size().Width, pos.Y))

	menuDismiss := menu.OnDismiss // this dismisses the menu stack
	menu.OnDismiss = func() {
		menuDismiss()
		m.bar.Hide() // dismiss the overlay menu bar
		m.canvas.setMenu(nil)
	}
}

func (m *menuLabel) CreateRenderer() fyne.WidgetRenderer {
	label := widget.NewLabel(m.menu.Label)
	box := container.NewHBox(layout.NewSpacer(), label, layout.NewSpacer(), widget.NewIcon(theme.MenuExpandIcon()))

	return &menuLabelRenderer{menu: m, content: box}
}

func newMenuLabel(item *fyne.Menu, parent *fyne.Container, c *mobileCanvas) *menuLabel {
	return &menuLabel{menu: item, bar: parent, canvas: c}
}

func (c *mobileCanvas) showMenu(menu *fyne.MainMenu) {
	var panel *fyne.Container
	top := container.NewHBox(widget.NewButtonWithIcon("", theme.CancelIcon(), func() {
		panel.Hide()
		c.setMenu(nil)
	}))
	panel = container.NewVBox(top)
	for _, item := range menu.Items {
		panel.Add(newMenuLabel(item, panel, c))
	}
	shadow := canvas.NewHorizontalGradient(theme.ShadowColor(), color.Transparent)
	c.setMenu(container.NewWithoutLayout(panel, shadow))

	safePos, safeSize := c.InteractiveArea()
	panel.Move(safePos)
	panel.Resize(fyne.NewSize(panel.MinSize().Width+theme.Padding(), safeSize.Height))
	shadow.Resize(fyne.NewSize(theme.Padding()/2, safeSize.Height))
	shadow.Move(fyne.NewPos(panel.Size().Width+safePos.X, safePos.Y))
}

func (d *mobileDriver) findMenu(win *window) *fyne.MainMenu {
	if win.menu != nil {
		return win.menu
	}

	matched := false
	for x := len(d.windows) - 1; x >= 0; x-- {
		w := d.windows[x]
		if !matched {
			if w == win {
				matched = true
			}
			continue
		}

		if w.(*window).menu != nil {
			return w.(*window).menu
		}
	}

	return nil
}

type menuLabelRenderer struct {
	menu    *menuLabel
	content *fyne.Container
}

func (m *menuLabelRenderer) BackgroundColor() color.Color {
	return theme.BackgroundColor()
}

func (m *menuLabelRenderer) Destroy() {
}

func (m *menuLabelRenderer) Layout(size fyne.Size) {
	m.content.Resize(size)
}

func (m *menuLabelRenderer) MinSize() fyne.Size {
	return m.content.MinSize()
}

func (m *menuLabelRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{m.content}
}

func (m *menuLabelRenderer) Refresh() {
	m.content.Refresh()
}
