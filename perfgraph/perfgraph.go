package perfgraph

import (
	"fmt"
	"github.com/jxo/davinci/vg"
	"time"
)

const (
	graphHistoryCount = 100
	gpuQueryCount     = 5
)

type RenderStyle int

const (
	RenderFPS RenderStyle = iota
	RenderMS
	RenderPercent
)

var backgroundColor = vg.RGBA(0, 0, 0, 128)
var graphColor = vg.RGBA(255, 192, 0, 128)
var titleTextColor = vg.RGBA(255, 192, 0, 128)
var fpsTextColor = vg.RGBA(240, 240, 240, 255)
var percentTextColor = vg.RGBA(240, 240, 240, 255)
var msTextColor = vg.RGBA(240, 240, 240, 160)

// PerfGraph shows FPS counter on DaVinci application
type PerfGraph struct {
	name     string
	fontFace string
	style    RenderStyle
	values   [graphHistoryCount]float32
	head     int

	startTime      time.Time
	lastUpdateTime time.Time
}

type GPUTimer struct {
	supported int
	cur, ret  int
	queries   [gpuQueryCount]uint
}

// NewPerfGraph creates PerfGraph instance
func NewPerfGraph(name, fontFace string, style RenderStyle) *PerfGraph {
	return &PerfGraph{
		name:           name,
		fontFace:       fontFace,
		style:          style,
		startTime:      time.Now(),
		lastUpdateTime: time.Now(),
	}
}

// UpdateGraph updates timer it is needed to show graph
func (pg *PerfGraph) UpdateGraph() (timeFromStart, frameTime float32) {
	timeNow := time.Now()
	timeFromStart = float32(timeNow.Sub(pg.startTime)/time.Millisecond) * 0.001
	frameTime = float32(timeNow.Sub(pg.lastUpdateTime)/time.Millisecond) * 0.001
	pg.lastUpdateTime = timeNow

	pg.head = (pg.head + 1) % graphHistoryCount
	pg.values[pg.head] = frameTime
	return
}

// RenderGraph shows graph
func (pg *PerfGraph) RenderGraph(ctx *vg.Context, x, y float32) {
	avg := pg.GetGraphAverage()
	var w float32 = 200
	var h float32 = 35
	var v, vx, vy float32

	ctx.BeginPath()
	ctx.Rect(x, y, w, h)
	ctx.SetFillColor(backgroundColor)
	ctx.Fill()

	ctx.BeginPath()
	ctx.MoveTo(x, y+h)
	if pg.style == RenderFPS {
		for i := 0; i < graphHistoryCount; i++ {
			v = float32(1.0) / float32(0.00001+pg.values[(pg.head+i)%graphHistoryCount])
			if v > 80.0 {
				v = 80.0
			}
			vx = x + float32(i)/float32(graphHistoryCount-1)*w
			vy = y + h - ((v / 80.0) * h)
			ctx.LineTo(vx, vy)
		}
	} else if pg.style == RenderPercent {
		for i := 0; i < graphHistoryCount; i++ {
			v = float32(pg.values[(pg.head+i)%graphHistoryCount])
			if v > 100.0 {
				v = 100.0
			}
			vx = x + float32(i)/float32(graphHistoryCount-1)*w
			vy = y + h - v*0.01*h
			ctx.LineTo(vx, vy)
		}
	} else {
		for i := 0; i < graphHistoryCount; i++ {
			v = float32(pg.values[(pg.head+i)%graphHistoryCount] * 1000)
			if v > 20.0 {
				v = 20.0
			}
			vx = x + float32(i)/float32(graphHistoryCount-1)*w
			vy = y + h - v*0.05*h
			ctx.LineTo(vx, vy)
		}
	}
	ctx.LineTo(x+w, y+h)
	ctx.SetFillColor(graphColor)
	ctx.Fill()

	ctx.SetFontFace(pg.fontFace)

	if len(pg.name) > 0 {
		ctx.SetFontSize(14.0)
		ctx.SetTextAlign(vg.AlignLeft | vg.AlignTop)
		ctx.SetFillColor(titleTextColor)
		ctx.Text(x+3, y+1, pg.name)
	}

	if pg.style == RenderFPS {
		ctx.SetFontSize(18.0)
		ctx.SetTextAlign(vg.AlignRight | vg.AlignTop)
		ctx.SetFillColor(fpsTextColor)
		ctx.Text(x+w-3, y+1, fmt.Sprintf("%.2f FPS", 1.0/avg))

		ctx.SetFontSize(15.0)
		ctx.SetTextAlign(vg.AlignRight | vg.AlignBottom)
		ctx.SetFillColor(msTextColor)
		ctx.Text(x+w-3, y+h+1, fmt.Sprintf("%.2f ms", avg*1000.0))
	} else if pg.style == RenderPercent {
		ctx.SetFontSize(18.0)
		ctx.SetTextAlign(vg.AlignRight | vg.AlignTop)
		ctx.SetFillColor(percentTextColor)
		ctx.Text(x+w-3, y+1, fmt.Sprintf("%.1f %%", avg))
	} else {
		ctx.SetFontSize(18.0)
		ctx.SetTextAlign(vg.AlignRight | vg.AlignTop)
		ctx.SetFillColor(msTextColor)
		ctx.Text(x+w-3, y+1, fmt.Sprintf("%.2f ms", avg*1000.0))
	}
}

// GetGraphAverage returns average value of graph.
func (pg *PerfGraph) GetGraphAverage() float32 {
	var average float32
	for _, value := range pg.values {
		average += value
	}
	return average / float32(graphHistoryCount)
}

func NewGPUTimer() *GPUTimer {
	return &GPUTimer{}
}

func (timer *GPUTimer) Start() {

}

func (timer *GPUTimer) Stop(times []float32) {

}
