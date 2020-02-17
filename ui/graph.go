package ui

import (
	"github.com/jxo/davinci/vg"
)

type Graph struct {
	WidgetImplement

	caption, header, footer                     string
	backgroundColor, foregroundColor, textColor vg.Color
	values                                      []float32
}

func NewGraph(parent Widget, captions ...string) *Graph {
	var caption string
	switch len(captions) {
	case 0:
		caption = "Untitled"
	case 1:
		caption = captions[0]
	default:
		panic("NewGraph can accept only one extra parameter (label)")
	}
	graph := &Graph{
		caption:         caption,
		backgroundColor: vg.MONO(20, 128),
		foregroundColor: vg.RGBA(255, 192, 0, 128),
		textColor:       vg.MONO(240, 192),
	}
	InitWidget(graph, parent)
	return graph
}

func (g *Graph) Caption() string {
	return g.caption
}

func (g *Graph) SetCaption(caption string) {
	g.caption = caption
}

func (g *Graph) Header() string {
	return g.header
}

func (g *Graph) SetHeader(header string) {
	g.header = header
}

func (g *Graph) Footer() string {
	return g.footer
}

func (g *Graph) SetFooter(footer string) {
	g.footer = footer
}

func (g *Graph) BackgroundColor() vg.Color {
	return g.backgroundColor
}

func (g *Graph) SetBackgroundColor(color vg.Color) {
	g.backgroundColor = color
}

func (g *Graph) ForegroundColor() vg.Color {
	return g.foregroundColor
}

func (g *Graph) SetForegroundColor(color vg.Color) {
	g.foregroundColor = color
}

func (g *Graph) TextColor() vg.Color {
	return g.textColor
}

func (g *Graph) SetTextColor(color vg.Color) {
	g.textColor = color
}

func (g *Graph) Values() []float32 {
	return g.values
}

func (g *Graph) SetValues(values []float32) {
	g.values = values
}

func (g *Graph) PreferredSize(self Widget, ctx *vg.Context) (int, int) {
	return 180, 45

}

func (g *Graph) Draw(self Widget, ctx *vg.Context) {
	g.WidgetImplement.Draw(self, ctx)

	x := float32(g.x)
	y := float32(g.y)
	w := float32(g.w)
	h := float32(g.h)

	ctx.BeginPath()
	ctx.Rect(x, y, w, h)
	ctx.SetFillColor(g.backgroundColor)
	ctx.Fill()

	if len(g.values) < 2 {
		return
	}

	ctx.BeginPath()
	ctx.MoveTo(x, y+h)
	dx := float32(len(g.values) - 1)
	for i, v := range g.values {
		vx := x + float32(i)*w/dx
		vy := y + (1.0-v)*h
		ctx.LineTo(vx, vy)
	}

	ctx.LineTo(x+w, y+h)
	ctx.SetStrokeColor(vg.MONO(100, 255))
	ctx.Stroke()
	ctx.SetFillColor(g.foregroundColor)
	ctx.Fill()

	ctx.SetFontFace(g.theme.FontNormal)
	ctx.SetFillColor(g.textColor)
	if g.caption != "" {
		ctx.SetFontSize(14)
		ctx.SetTextAlign(vg.AlignLeft | vg.AlignTop)
		ctx.Text(x+3, y+1, g.caption)
	}

	if g.header != "" {
		ctx.SetFontSize(18)
		ctx.SetTextAlign(vg.AlignRight | vg.AlignTop)
		ctx.Text(x+w-3, y+1, g.header)
	}

	if g.footer != "" {
		ctx.SetFontSize(15)
		ctx.SetTextAlign(vg.AlignRight | vg.AlignBottom)
		ctx.Text(x+w-3, y+h-1, g.footer)
	}

	ctx.BeginPath()
	ctx.Rect(x, y, w, h)
	ctx.SetStrokeColor(vg.MONO(100, 255))
	ctx.Stroke()
}

func (g *Graph) String() string {
	return g.StringHelper("Graph", g.caption)
}
