package ui

import (
	"fmt"
	"github.com/goxjs/glfw"
	"github.com/jxo/davinci/vg"
)

type Panel struct {
	WidgetImplement
	title       string
	buttonPanel Widget
	modal       bool
	drag        bool
	draggable   bool
	depth       int
}

type IPanel interface {
	Widget
	RefreshRelativePlacement()
	SetDepth(d int)
}

func NewPanel(parent Widget, title string) *Panel {
	if title == "" {
		title = "Untitled"
	}
	panel := &Panel{
		title:     title,
		draggable: true,
	}
	InitWidget(panel, parent)
	return panel
}

// Title() returns the panel title
func (w *Panel) Title() string {
	return w.title
}

// SetTitle() sets the panel title
func (w *Panel) SetTitle(title string) {
	w.title = title
}

// Modal() returns is this a model dialog?
func (w *Panel) Modal() bool {
	return w.modal
}

// SetModal() set whether or not this is a modal dialog
func (w *Panel) SetModal(m bool) {
	w.modal = m
}

func (w *Panel) Draggable() bool {
	return w.draggable
}

func (w *Panel) SetDraggable(flag bool) {
	w.draggable = flag
}

func (w *Panel) ButtonPanel() Widget {
	if w.buttonPanel == nil {
		w.buttonPanel = NewWidget(w)
		w.buttonPanel.SetLayout(NewBoxLayout(Horizontal, Middle, 0, 4))
	}
	return w.buttonPanel
}

// Dispose() disposes the panel
func (w *Panel) Dispose() {
	var widget Widget = w
	var parent Widget = w.Parent()
	for parent != nil {
		widget = parent
		parent = widget.Parent()
	}
	win := widget.(*Window)
	win.DisposePanel(w)
}

// Center() makes the panel center in the current Window
func (w *Panel) Center() {
	var widget Widget = w
	var parent Widget = w.Parent()
	for parent != nil {
		widget = parent
		parent = widget.Parent()
	}
	win := widget.(*Window)
	win.CenterPanel(w)
}

// RefreshRelativePlacement is internal helper function to maintain nested panel position values; overridden in \ref Popup
func (w *Panel) RefreshRelativePlacement() {
	// overridden in Popup
}

func (w *Panel) MouseButtonEvent(self Widget, x, y int, button glfw.MouseButton, down bool, modifier glfw.ModifierKey) bool {
	if w.WidgetImplement.MouseButtonEvent(self, x, y, button, down, modifier) {
		return true
	}
	if button == glfw.MouseButton1 && w.draggable {
		w.drag = down && (y-w.y) < w.theme.PanelHeaderHeight
		return true
	}
	return false
}

func (w *Panel) MouseDragEvent(self Widget, x, y, relX, relY, button int, modifier glfw.ModifierKey) bool {
	if w.drag && (button&1<<uint(glfw.MouseButton1)) != 0 {
		pW, pH := self.Parent().Size()
		w.x = clampI(w.x+relX, 0, pW-w.w)
		w.y = clampI(w.y+relY, 0, pH-w.h)
		return true
	}
	return false
}

func (w *Panel) ScrollEvent(self Widget, x, y, relX, relY int) bool {
	w.WidgetImplement.ScrollEvent(self, x, y, relX, relY)
	return true
}

func (w *Panel) PreferredSize(self Widget, ctx *vg.Context) (int, int) {
	if w.buttonPanel != nil {
		w.buttonPanel.SetVisible(false)
	}
	width, height := w.WidgetImplement.PreferredSize(self, ctx)
	if w.buttonPanel != nil {
		w.buttonPanel.SetVisible(true)
	}
	ctx.SetFontSize(18.0)
	ctx.SetFontFace(w.theme.FontBold)
	_, bounds := ctx.TextBounds(0, 0, w.title)

	return maxI(width, int(bounds[2]-bounds[0])+20), maxI(height, int(bounds[3]-bounds[1]))
}

func (w *Panel) OnPerformLayout(self Widget, ctx *vg.Context) {
	if w.buttonPanel == nil {
		w.WidgetImplement.OnPerformLayout(self, ctx)
	} else {
		w.buttonPanel.SetVisible(false)
		w.WidgetImplement.OnPerformLayout(self, ctx)
		for _, c := range w.buttonPanel.Children() {
			c.SetFixedSize(22, 22)
			c.SetFontSize(15)
		}
		w.buttonPanel.SetVisible(true)
		w.buttonPanel.SetSize(w.Width(), 22)
		panelW, _ := w.buttonPanel.PreferredSize(w.buttonPanel, ctx)
		w.buttonPanel.SetPosition(w.Width()-(panelW+5), 3)
		w.buttonPanel.OnPerformLayout(w.buttonPanel, ctx)
	}
}

func (w *Panel) Draw(self Widget, ctx *vg.Context) {
	ds := float32(w.theme.PanelDropShadowSize)
	cr := float32(w.theme.PanelCornerRadius)
	hh := float32(w.theme.PanelHeaderHeight)

	// Draw panel
	wx := float32(w.x)
	wy := float32(w.y)
	ww := float32(w.w)
	wh := float32(w.h)
	ctx.Save()
	ctx.BeginPath()
	ctx.RoundedRect(wx, wy, ww, wh, cr)
	if w.mouseFocus {
		ctx.SetFillColor(w.theme.PanelFillFocused)
	} else {
		ctx.SetFillColor(w.theme.PanelFillUnfocused)
	}
	ctx.Fill()

	// Draw a drop shadow
	shadowPaint := vg.BoxGradient(wx, wy, ww, wh, cr*2, ds*2, w.theme.DropShadow, w.theme.Transparent)
	ctx.BeginPath()
	ctx.Rect(wx-ds, wy-ds, ww+ds*2, wh+ds*2)
	ctx.RoundedRect(wx, wy, ww, wh, cr)
	ctx.PathWinding(vg.Hole)
	ctx.SetFillPaint(shadowPaint)
	ctx.Fill()

	if w.title != "" {
		headerPaint := vg.LinearGradient(wx, wy, ww, wh+hh, w.theme.PanelHeaderGradientTop, w.theme.PanelHeaderGradientBot)

		ctx.BeginPath()
		ctx.RoundedRect(wx, wy, ww, hh, cr)
		ctx.SetFillPaint(headerPaint)
		ctx.Fill()

		ctx.BeginPath()
		ctx.RoundedRect(wx, wy, ww, wh, cr)
		ctx.SetStrokeColor(w.theme.PanelHeaderSepTop)
		ctx.Scissor(wx, wy, ww, 0.5)
		ctx.Stroke()
		ctx.ResetScissor()

		ctx.BeginPath()
		ctx.MoveTo(wx+0.5, wy+hh-1.5)
		ctx.LineTo(wx+ww-0.5, wy+hh-1.5)
		ctx.SetStrokeColor(w.theme.PanelHeaderSepTop)
		ctx.Stroke()

		ctx.SetFontSize(18.0)
		ctx.SetFontFace(w.theme.FontBold)
		ctx.SetTextAlign(vg.AlignCenter | vg.AlignMiddle)
		ctx.SetFontBlur(2.0)
		ctx.SetFillColor(w.theme.DropShadow)
		ctx.Text(wx+ww*0.5, wy+hh*0.5, w.title)
		ctx.SetFontBlur(0.0)
		if w.focused {
			ctx.SetFillColor(w.theme.PanelTitleFocused)
		} else {
			ctx.SetFillColor(w.theme.PanelTitleUnfocused)
		}
		ctx.Text(wx+ww*0.5, wy+hh*0.5-1, w.title)
	}
	ctx.Restore()
	w.WidgetImplement.Draw(self, ctx)
}

func (w *Panel) FindPanel() IPanel {
	return w
}

func (w *Panel) String() string {
	return w.StringHelper(fmt.Sprintf("Panel(%d)", w.Depth()), w.title)
}

func (w *Panel) SetDepth(d int) {
	w.depth = d
}

func (w *Panel) Depth() int {
	return w.depth
}
