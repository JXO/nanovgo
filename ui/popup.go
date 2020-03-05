package ui

import (
	"fmt"
	"github.com/jxo/davinci/vg"
)

type Popup struct {
	Panel
	parentPanel IPanel
	anchorX      int
	anchorY      int
	anchorHeight int
	vScroll      *VScrollPanel
	panel        Widget
}

func NewPopup(parent Widget, parentPanel IPanel) *Popup {
	popup := &Popup{
		parentPanel: parentPanel,
		anchorHeight: 30,
	}
	InitWidget(popup, parent)
	popup.vScroll = NewVScrollPanel(popup)
	popup.panel = NewVScrollPanelChild(popup.vScroll)
	return popup
}

// SetAnchorPosition() sets the anchor position in the parent panel; the placement of the popup is relative to it
func (p *Popup) SetAnchorPosition(x, y int) {
	p.anchorX = x
	p.anchorY = y
}

// AnchorPosition() returns the anchor position in the parent panel; the placement of the popup is relative to it
func (p *Popup) AnchorPosition() (int, int) {
	return p.anchorX, p.anchorY
}

// SetAnchorHeight() sets the anchor height; this determines the vertical shift relative to the anchor position
func (p *Popup) SetAnchorHeight(h int) {
	p.anchorHeight = h
}

// AnchorHeight() returns the anchor height; this determines the vertical shift relative to the anchor position
func (p *Popup) AnchorHeight() int {
	return p.anchorHeight
}

// SetParentPanel() sets the parent panel of the popup
func (p *Popup) SetParentPanel(w *Panel) {
	p.parentPanel = w
}

// ParentPanel() returns the parent panel of the popup
func (p *Popup) ParentPanel() IPanel {
	return p.parentPanel
}

func (p *Popup) OnPerformLayout(self Widget, ctx *vg.Context) {
	if p.layout != nil || len(p.children) != 1 {
		p.WidgetImplement.OnPerformLayout(self, ctx)
	} else {
		p.children[0].SetPosition(0, 0)
		p.children[0].SetSize(p.w, p.h)
		p.children[0].OnPerformLayout(p.children[0], ctx)
	}
}

func (p *Popup) IsPositionAbsolute() bool {
	return true
}

func (p *Popup) Draw(self Widget, ctx *vg.Context) {
	p.RefreshRelativePlacement()

	if !p.visible {
		return
	}
	ds := float32(p.theme.PanelDropShadowSize)
	cr := float32(p.theme.PanelCornerRadius)

	px := float32(p.x)
	py := float32(p.y)
	pw := float32(p.w)
	ph := float32(p.h)
	ah := float32(p.anchorHeight)

	/* Draw a drop shadow */
	shadowPaint := vg.BoxGradient(px, py, pw, ph, cr*2, ds*2, p.theme.DropShadow, p.theme.Transparent)
	ctx.BeginPath()
	ctx.Rect(px-ds, py-ds, pw+ds*2, ph+ds*2)
	ctx.RoundedRect(px, py, pw, ph, cr)
	ctx.PathWinding(vg.Hole)
	ctx.SetFillPaint(shadowPaint)
	ctx.Fill()

	/* Draw panel */
	ctx.BeginPath()
	ctx.RoundedRect(px, py, pw, ph, cr)

	ctx.MoveTo(px-15, py+ah)
	ctx.LineTo(px+1, py+ah-15)
	ctx.LineTo(px+1, py+ah+15)

	ctx.SetFillColor(p.theme.PanelPopup)

	ctx.Fill()

	p.WidgetImplement.Draw(self, ctx)
}

// RefreshRelativePlacement is internal helper function to maintain nested panel position values; overridden in \ref Popup
func (p *Popup) RefreshRelativePlacement() {
	p.parentPanel.RefreshRelativePlacement()
	p.visible = p.visible && p.parentPanel.VisibleRecursive()
	x, y := p.parentPanel.Position()
	p.x = x + p.anchorX
	p.y = y + p.anchorY - p.anchorHeight
}

func (p *Popup) FindPanel() IPanel {
	return p
}

func (p *Popup) String() string {
	return p.StringHelper(fmt.Sprintf("Popup(%d)", p.Depth()), "")
}
