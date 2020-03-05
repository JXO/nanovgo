package ui

import (
	"fmt"
	"github.com/jxo/davinci/vg"
)

type ColorPicker struct {
	PopupButton

	callback   func(color vg.Color)
	colorWheel *ColorWheel
	pickButton *Button
}

func NewColorPicker(parent Widget, colors ...vg.Color) *ColorPicker {
	var color vg.Color
	switch len(colors) {
	case 0:
		color = vg.RGBAf(1.0, 0.0, 0.0, 1.0)
	case 1:
		color = colors[0]
	default:
		panic("NewColorPicker can accept only one extra parameter (color)")
	}

	colorPicker := &ColorPicker{}

	// init PopupButton member
	colorPicker.chevronIcon = IconRightOpen
	colorPicker.SetIconPosition(ButtonIconLeftCentered)
	colorPicker.SetFlags(ToggleButtonType | PopupButtonType)
	parentPanel := parent.FindPanel()

	colorPicker.popup = NewPopup(parentPanel.Parent(), parentPanel)
	colorPicker.popup.panel.SetLayout(NewGroupLayout())

	colorPicker.colorWheel = NewColorWheel(colorPicker.popup.panel)

	colorPicker.pickButton = NewButton(colorPicker.popup.panel, "Pick")
	colorPicker.pickButton.SetFixedSize(100, 25)

	InitWidget(colorPicker, parent)

	colorPicker.SetColor(color)

	colorPicker.PopupButton.SetChangeCallback(func(flag bool) {
		colorPicker.SetColor(colorPicker.BackgroundColor())
		if colorPicker.callback != nil {
			colorPicker.callback(colorPicker.BackgroundColor())
		}
	})

	colorPicker.colorWheel.SetCallback(func(color vg.Color) {
		colorPicker.pickButton.SetBackgroundColor(color)
		colorPicker.pickButton.SetTextColor(color.ContrastingColor())
	})

	colorPicker.pickButton.SetCallback(func() {
		color := colorPicker.colorWheel.Color()
		colorPicker.SetPushed(false)
		colorPicker.SetColor(color)
		if colorPicker.callback != nil {
			colorPicker.callback(colorPicker.BackgroundColor())
		}
	})

	return colorPicker
}

func (c *ColorPicker) SetCallback(callback func(color vg.Color)) {
	c.callback = callback
}

func (c *ColorPicker) Color() vg.Color {
	return c.BackgroundColor()
}

func (c *ColorPicker) SetColor(color vg.Color) {
	if !c.pushed {
		fgColor := color.ContrastingColor()
		c.SetBackgroundColor(color)
		c.SetTextColor(fgColor)
		c.colorWheel.SetColor(color)
		c.pickButton.SetBackgroundColor(color)
		c.pickButton.SetTextColor(fgColor)
	}
}

func (c *ColorPicker) String() string {
	cw := c.colorWheel
	return c.StringHelper("ColorPicker", fmt.Sprintf("h:%f s:%f l:%f", cw.hue, cw.saturation, cw.lightness))
}
