package dark

import (
    "github.com/jxo/davinci/vg"
)

type WidgetTheme struct {
    // color of widget box outline
    OutLineColor vg.Color
    // color of widget item (meaning changes depending on class)
    ItemColor vg.Color
    // fill color of widget box
    InnerColor vg.Color
    // fill color of widget box when active
    InnerSelectedColor vg.Color
    // color of text label
    TextColor vg.Color
    // color of text label when active
    TextSelectedColor vg.Color
    // delta modifier for upper part of gradient (-100 to 100)
    ShadeTop int
    ShadeDown int
}

type NodeTheme struct {
    // inner color of selected node (and down arrow)
    NodeSelectedColor vg.Color
    // outline of wires
    WireColor vg.Color
    // color of text label when active
    TextSelectedColor vg.Color
    // inner color of active node (and dragged wire)
    ActiveNodeColor vg.Color
    // color of selected wire
    WireSelectedColor vg.Color
    // color of background of node
    NodeBackdropColor vg.Color

    // how much a noodle curves (0 to 10)
    NoodleCurving int
}

// describes the theme used to draw widgets
type Theme struct {
    // the background color of panels and panels
    BackgroundColor vg.Color
    // alpha of disabled widget groups
    // can be used in conjunction with vg.GlobalAlpha()
    DisabledAlpha float32
    Regular  WidgetTheme
    Tool WidgetTheme
    Radio WidgetTheme
    TextField WidgetTheme
    Option WidgetTheme
    Choice WidgetTheme
    NumberField WidgetTheme
    Slider WidgetTheme
    ScrollBar WidgetTheme
    ToolTip WidgetTheme
    MenuTheme WidgetTheme
    MenuItem  WidgetTheme
    Node NodeTheme

    WidgetHeight int

    // tool button width (if icon only)
    ToolWidth int

    NodePortRadius int
    NodeMarginTop int
    NodeMarginDown int
    NodeMarginSide int
    NodeTitleHeight int
    NodeArrowAreaWidth int

    // size of splitter corner click area
    SplitterAreaSize int

    // width of vertical scrollbar
    ScrollBarWidth int
    ScrollBarHeight int

    // vertical spacing
    VSpacing int
    // vertical spacing between groups
    VSpacingGroup int
    // horizontal spacing
    HSpacing int

    // shade intensity of beveled panels
    BevelShade int
    // shade intensity of beveled insets
    InsetBevelShade int
    // shade intensity of hovered inner boxes
    HoverShade int
    // shade intensity of splitter bevels
    SplitterShade int
}

// How text on a control is aligned: vg.Align

// States altering the styling of a widget
type WidgetState uint8
const (
    // not interacting
    WSDefault WidgetState = iota
    // the mouse is hovering over the control
    WSHover
    // the widget is activated (pressed) or in an active state (toggled)
    WSActive
)

// flags indiacting which corners are sharp (for grouping widgets)
type CornerFlag int
const (
    // all corners are round
    CornerNone CornerFlag = 0

    // sharp top left corner
    CornerTopLeft CornerFlag = 1 << iota
    CornerTopRight
    CornerDownRight
    CornerDownLeft

    CornerALL CornerFlag = 0xF
    // top border is sharp
    CornerTop CornerFlag = 3
    CornerDown CornerFlag = 0xC
    CornerLeft CornerFlag = 9
    CornerRight CornerFlag = 6
)

// XXX export vg.clampF
func clampF(a, min, max float32) float32 {
	if a < min {
		return min
	}
	if a > max {
		return max
	}
	return a
}

func offsetColor(c vg.Color, delta int) vg.Color {
    offset := float32(delta) / 255.0
    if delta == 0 {
        return c
    }

    return vg.RGBf(
        clampF(c.R + offset, 0, 1),
        clampF(c.G + offset, 0, 1),
        clampF(c.B + offset, 0, 1))
}

// Draw a node port at the given position filled with the given color
func nodePort(ctx *vg.Context, theme *Theme, x, y float32, s WidgetState, c vg.Color) {
    ctx.BeginPath()
    ctx.Circle(x, y, float32(theme.NodePortRadius))
    ctx.SetStrokeColor(theme.Node.WireColor)
    ctx.SetStrokeWidth(1.0)
    ctx.Stroke()
    if s != WSDefault {
        c = offsetColor(c, theme.HoverShade)
    }
    ctx.SetFillColor(c)
    ctx.Fill()
}

func nodeWireColor(theme *Theme, s WidgetState) (c vg.Color) {
    switch s {
    case WSHover:
        c = theme.WireSelectedColor
    case WSActive:
        c = theme.ActiveNodeColor
    case WSDefault:
        fallthrough
    default:
        c = vg.RGBf(0.5, 0.5, 0.5)
    }
    return
}

// Draw a node wire originating at (x0, y0) and floating to (x1, y1),
// with a colored gradient based on the states s0 and s1:
// Default: default wire color
// Hover: selected wire color
// Active: dragged wire color
func nodeWire(ctx *vg.Context, theme *Theme, x0, y0, x1, y1 float32, s0, s1 WidgetState) {
    c0 := nodeWireColor(theme, s0)
    c1 := nodeWireColor(theme, s1)


}
