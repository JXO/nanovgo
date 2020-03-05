package ui

import (
	"bytes"
	"github.com/goxjs/gl"
	"github.com/goxjs/glfw"
	"github.com/jxo/davinci/vg"
	"runtime"
    "time"
//    "fmt"
)

type Window struct {
	WidgetImplement
    EventQueue

    renderQueue            EventQueue
	win *glfw.Window
	context                *vg.Context
	cursors                [3]int
	cursor                 Cursor
	focusPath              []Widget
	fbW, fbH               int
	pixelRatio             float32
	mouseState int
    modifiers glfw.ModifierKey
	mousePosX, mousePosY   int
	dragActive             bool
	dragWidget             Widget
	lastInteraction        float32
	backgroundColor        vg.Color
	caption                string
	shutdownGLFWOnDestruct bool

    onLoadCallback func()
	drawContentsCallback func()
	dropEventCallback    func([]string) bool
	resizeEventCallback  func(x, y int) bool
}

func (w *Window) MakeCurrent() {
    w.win.MakeContextCurrent()
}

func (w *Window) InitializeGfx() {
    w.MakeCurrent()
	var err error
	w.context, err = vg.NewContext(vg.StencilStrokes | vg.AntiAlias)
	if err != nil {
		panic(err)
	}

	w.context.CreateFontFromMemory("sans", MustAsset("fonts/Roboto-Regular.ttf"), 0)
	w.context.CreateFontFromMemory("sans-bold", MustAsset("fonts/Roboto-Bold.ttf"), 0)
	w.context.CreateFontFromMemory("icons", MustAsset("fonts/entypo.ttf"), 0)
}

// Caption() gets the panel title bar caption
func (w *Window) Caption() string {
	return w.caption
}

// SetCaption() sets the panel title bar caption
func (w *Window) SetCaption(caption string) {
	if w.caption != caption {
		w.win.SetTitle(caption)
		w.caption = caption
	}
}

// BackgroundColor() returns the screen's background color
func (w *Window) BackgroundColor() vg.Color {
	return w.backgroundColor
}

// SetBackgroundColor() sets the screen's background color
func (w *Window) SetBackgroundColor(color vg.Color) {
	w.backgroundColor = color
	w.backgroundColor.A = 1.0
}

// SetVisible() sets the top-level panel visibility (no effect on full-screen panels)
func (w *Window) SetVisible(flag bool) {
	if w.visible != flag {
		w.visible = flag
		if flag {
			w.win.Show()
		} else {
			w.win.Hide()
		}
	}
}

// SetSize() sets window size
func (w *Window) SetSize(width, height int) {
	w.WidgetImplement.SetSize(width, height)
	w.win.SetSize(width, height)
}

// DrawAll() draws the Window contents
func (w *Window) OnPaint() {
    w.renderQueue.Send(GLEvent {
        T : time.Now(),
        F : func (t time.Time) {
            if d := time.Since(t).Milliseconds(); d > 5 {
                //fmt.Println(d)
                return
            }

            w.MakeCurrent()
            gl.ClearColor(w.backgroundColor.R, w.backgroundColor.G, w.backgroundColor.B, 1.0)
            gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT | gl.STENCIL_BUFFER_BIT)

            if w.drawContentsCallback != nil {
                w.drawContentsCallback()
            }

            w.drawWidgets()
            w.win.SwapBuffers()

        },
    })
}

// SetResizeEventCallback() sets panel resize event handler
func (w *Window) SetResizeEventCallback(callback func(x, y int) bool) {
	w.resizeEventCallback = callback
}

// SetDrawContentsCallback() sets event handler to use OpenGL draw call
func (w *Window) SetDrawContentsCallback(callback func()) {
	w.drawContentsCallback = callback
}

// SetDropEventCallback() sets event handler to handle a file drop event
func (w *Window) SetDropEventCallback(callback func([] string) bool) {
	w.dropEventCallback = callback
}

func (w *Window) SetOnLoad(callback func()) {
    w.onLoadCallback = callback
    if callback != nil {
        w.renderQueue.Send(callback)
    }
}

// KeyboardEvent() is a default key event handler
func (w *Window) OnKeyEvent(e *KeyEvent) bool {
	if len(w.focusPath) > 1 {
		for i := len(w.focusPath) - 2; i >= 0; i-- {
			path := w.focusPath[i]
			if path.Focused() && path.KeyboardEvent(path, e.Key, e.ScanCode, e.Action, e.Modifier) {
				return true
			}
		}
	}

	return false
}

// KeyboardCharacterEvent() is a text input event handler: codepoint is native endian UTF-32 format
func (w *Window) OnRuneEvent(e *RuneEvent) bool {
	if len(w.focusPath) > 1 {
		for i := len(w.focusPath) - 2; i >= 0; i-- {
			path := w.focusPath[i]
			if path.Focused() && path.KeyboardCharacterEvent(path, e.R) {
				return true
			}
		}
	}
	return false
}

/*
// IMEPreeditEvent() handles preedit text changes of IME (default implementation: do nothing)
func (w *Window) IMEPreeditEvent(self Widget, text []rune, blocks []int, focusedBlock int) bool {
	if len(w.focusPath) > 1 {
		for i := len(w.focusPath) - 2; i >= 0; i-- {
			path := w.focusPath[i]
			if path.Focused() && path.IMEPreeditEvent(path, text, blocks, focusedBlock) {
				return true
			}
		}
	}
	return false
}

// IMEStatusEvent() handles IME status change event (default implementation: do nothing)
func (s *Window) IMEStatusEvent(self Widget) bool {
	if len(w.focusPath) > 1 {
		for i := len(w.focusPath) - 2; i >= 0; i-- {
			path := w.focusPath[i]
			if path.Focused() && path.IMEStatusEvent(path) {
				return true
			}
		}
	}
	return false
}
*/

// MousePosition() returns the last observed mouse position value
func (w *Window) MousePosition() (int, int) {
	return w.mousePosX, w.mousePosY
}

// GLFWWindow() returns a pointer to the underlying GLFW panel data structure
func (w *Window) GLFWWindow() *glfw.Window {
	return w.win
}

// NVGContext() returns a pointer to the underlying nanoVGo draw context
func (w *Window) NVGContext() *vg.Context {
	return w.context
}

func (w *Window) SetShutdownGLFWOnDestruct(v bool) {
	w.shutdownGLFWOnDestruct = v
}

func (w *Window) ShutdownGLFWOnDestruct() bool {
	return w.shutdownGLFWOnDestruct
}

// UpdateFocus is an internal helper function
func (w *Window) UpdateFocus(widget Widget) {
	for _, w := range w.focusPath {
		if w.Focused() {
			w.FocusEvent(widget, false)
		}
	}
	w.focusPath = w.focusPath[:0]
	var panel *Panel
	for widget != nil {
		w.focusPath = append(w.focusPath, widget)
		if _, ok := widget.(*Panel); ok {
			panel = widget.(*Panel)
		}
		widget = widget.Parent()
	}
	for _, w := range w.focusPath {
		w.FocusEvent(w, true)
	}

	if panel != nil {
		w.MovePanelToFront(panel)
	}

    w.Send(PaintEvent{})
}

// DisposePanel is an internal helper function
func (w *Window) DisposePanel(panel *Panel) {
	find := false
	for _, w := range w.focusPath {
		if w == panel {
			find = true
			break
		}
	}

	if find {
		w.focusPath = w.focusPath[:0]
	}

	if w.dragWidget == panel {
		w.dragWidget = nil
	}

	panel.Parent().RemoveChild(panel)
}

// CenterPanel is an internal helper function
func (w *Window) CenterPanel(panel *Panel) {
	width, height := panel.Size()
	if width == 0 && height == 0 {
		panel.SetSize(panel.PreferredSize(panel, w.context))
		panel.OnPerformLayout(panel, w.context)
	}

	ww, wh := panel.Size()
	pw, ph := panel.Parent().Size()
	panel.SetPosition((pw-ww)/2, (ph-wh)/2)
}

// MovePanelToFront is an internal helper function
func (w *Window) MovePanelToFront(panel IPanel) {
	parent := panel.Parent()
	maxDepth := 0
	for _, child := range parent.Children() {
		depth := child.Depth()
		if child != panel && maxDepth < depth {
			maxDepth = depth
		}
	}
	panel.SetDepth(maxDepth + 1)
	changed := true
	for changed {
		baseDepth := 0
		for _, child := range parent.Children() {
			if child == panel {
				baseDepth = child.Depth()
			}
		}
		changed = false
		for _, child := range parent.Children() {
			pw, ok := child.(*Popup)
			if ok && pw.ParentPanel() == panel && pw.Depth() < baseDepth {
				w.MovePanelToFront(pw)
				changed = true
				break
			}
		}
	}
}

//XXX
func (w *Window) PreeditCursorPos() (int, int, int) {
    return -1, -1, -1
	//return w.win.GetPreeditCursorPos()
}

// XXX
func (w *Window) SetPreeditCursorPos(x, y, h int) {
	//w.win.SetPreeditCursorPos(x, y, h)
}

func (w *Window) drawWidgets() {
	if !w.visible {
		return
	}

	w.fbW, w.fbH = w.win.GetFramebufferSize()
	w.w, w.h = w.win.GetSize()
	gl.Viewport(0, 0, w.fbW, w.fbH)

	w.pixelRatio = float32(w.fbW) / float32(w.w)
	w.context.BeginFrame(w.w, w.h, w.pixelRatio)
	w.Draw(w, w.context)
	elapsed := GetTime() - w.lastInteraction

	if elapsed > 0.5 {
		// Draw tooltips
		widget := w.FindWidget(w, w.mousePosX, w.mousePosY)
		if widget != nil && widget.Tooltip() != "" {
			var tooltipWidth float32 = 150
			ctx := w.context
			ctx.SetFontFace(w.theme.FontNormal)
			ctx.SetFontSize(15.0)
			ctx.SetTextAlign(vg.AlignCenter | vg.AlignTop)
			ctx.SetTextLineHeight(1.1)
			posX, posY := widget.AbsolutePosition()
			posX += widget.Width() / 2
			posY += widget.Height() + 10
			bounds := ctx.TextBoxBounds(float32(posX), float32(posY), tooltipWidth, widget.Tooltip())
			ctx.SetGlobalAlpha(minF(1.0, 2*(elapsed-0.5)) * 0.8)
			ctx.BeginPath()
			ctx.SetFillColor(vg.MONO(0, 255))
			h := (bounds[2] - bounds[0]) / 2
			ctx.RoundedRect(bounds[0]-4-h, bounds[1]-4, bounds[2]-bounds[0]+8, bounds[3]-bounds[1]+8, 3)
			px := (bounds[2]+bounds[0])/2 - h
			ctx.MoveTo(px, bounds[1]-10)
			ctx.LineTo(px+7, bounds[1]+1)
			ctx.LineTo(px-7, bounds[1]+1)
			ctx.Fill()

			ctx.SetFillColor(vg.MONO(255, 255))
			ctx.SetFontBlur(0.0)
			ctx.TextBox(float32(posX)-h, float32(posY), tooltipWidth, widget.Tooltip())

		}
	}

	w.context.EndFrame()
}

func (w *Window) mouseInWindow(x, y int) bool {
    if x < 0 || x > w.w || y < 0 || y > w.h {
        return false
    }
    return true
}

func (w *Window) OnMouseMove(e *MouseMoveEvent) bool {
	px := int(e.X) - 1
	py := int(e.Y) - 2

    if !w.mouseInWindow(px, py) {
        return false
    }

	ret := false
	w.lastInteraction = GetTime()
	if !w.dragActive {
		widget := w.FindWidget(w, int(e.X), int(e.Y))
		if widget != nil && widget.Cursor() != w.cursor {
			//w.cursor = widget.Cursor()
			//w.win.SetCursor()
		}
	} else {
		ax, ay := w.dragWidget.Parent().AbsolutePosition()
		ret = w.dragWidget.MouseDragEvent(w.dragWidget, px-ax, py-ay, px-w.mousePosX, py-w.mousePosY, w.mouseState, w.modifiers)
	}
	if !ret {
		ret = w.MouseMotionEvent(w, px, py, px-w.mousePosX, py-w.mousePosY, w.mouseState, w.modifiers)
    } else {
        w.Send(PaintEvent{})
    }
	w.mousePosX = px
	w.mousePosY = py
	return ret
}

func (w *Window) OnMouseButton(e *MouseButtonEvent) bool {
    if !w.mouseInWindow(w.mousePosX, w.mousePosY) {
        return false
    }

	w.lastInteraction = GetTime()
	if len(w.focusPath) > 1 {
		panel, ok := w.focusPath[len(w.focusPath)-2].(*Panel)
		if ok && panel.Modal() {
			if !panel.Contains(w.mousePosX, w.mousePosY) {
				return false
			}
		}
	}

	if e.Action == glfw.Press {
		w.mouseState |= 1 << uint(e.Button)
	} else {
		w.mouseState &= ^(1 << uint(e.Button))
	}

	dropWidget := w.FindWidget(w, w.mousePosX, w.mousePosY)
	if w.dragActive && e.Action == glfw.Release && dropWidget != w.dragWidget {
		ax, ay := w.dragWidget.Parent().AbsolutePosition()
		w.dragWidget.MouseButtonEvent(w.dragWidget, w.mousePosX-ax, w.mousePosY-ay, e.Button, false, e.Modifier)
	}

	if dropWidget != nil && dropWidget.Cursor() != w.cursor {
		//w.cursor = widget.Cursor()
		//w.win.SetCursor()
	}

	if e.Action == glfw.Press && e.Button == glfw.MouseButton1 {
		w.dragWidget = w.FindWidget(w, w.mousePosX, w.mousePosY)
		if w.dragWidget == w {
			w.dragWidget = nil
		}
		w.dragActive = w.dragWidget != nil
		if !w.dragActive {
			w.UpdateFocus(nil)
		}
	} else {
		w.dragActive = false
		w.dragWidget = nil
	}

	return w.MouseButtonEvent(w, w.mousePosX, w.mousePosY, e.Button, e.Action == glfw.Press, e.Modifier)
}

func (w *Window) KeyboardEvent(self Widget, key glfw.Key, scanCode int, action glfw.Action, modifier glfw.ModifierKey) bool {
	w.lastInteraction = GetTime()
	return w.KeyboardEvent(w, key, scanCode, action, modifier)
}

func (w *Window) KeyboardCharacterEvent(self Widget, codePoint rune) bool {
	w.lastInteraction = GetTime()
	return w.KeyboardCharacterEvent(w, codePoint)
}

func (w *Window) preeditCallbackEvent(text []rune, blocks []int, focusedBlock int) {
	w.lastInteraction = GetTime()
	w.IMEPreeditEvent(w, text, blocks, focusedBlock)
}

func (w *Window) imeStatusCallbackEvent() {
	w.lastInteraction = GetTime()
	w.IMEStatusEvent(w)
}

func (w *Window) OnDropEvent(e *DropEvent) bool {
	if w.dropEventCallback != nil {
		return w.dropEventCallback(e.Names)
	}
	return false
}

func (w *Window) OnScrollEvent(e *ScrollEvent) bool {
	w.lastInteraction = GetTime()

    //XXX
	if runtime.GOOS == "Windows" {
		e.Dx *= 32
		e.Dy *= 32
	}

	if len(w.focusPath) > 1 {
		panel, ok := w.focusPath[len(w.focusPath)-2].(*Panel)
		if ok && panel.Modal() {
			if !panel.Contains(w.mousePosX, w.mousePosY) {
				return false
			}
		}
	}

	return w.ScrollEvent(w, w.mousePosX, w.mousePosY, int(e.Dx), int(e.Dy))
}

func (w *Window) OnResizeEvent(e *ResizeEvent) bool {
	fbW, fbH := w.win.GetFramebufferSize()

	if (fbW == 0 && fbH == 0) && (e.Width== 0 && e.Height == 0) {
		return false
	}
	w.fbW = fbW
	w.fbH = fbH
	w.w = e.Width
	w.h = e.Height
	w.lastInteraction = GetTime()
	if w.resizeEventCallback != nil {
		return w.resizeEventCallback(int(float32(fbW)/w.pixelRatio), int(float32(fbH)/w.pixelRatio))
	}
	return false
}

func (w *Window) OnCloseEvent() bool {
    return true
}

func (w *Window) PerformLayout() {
	w.OnPerformLayout(w, w.context)
}

func (w *Window) String() string {
	return w.StringHelper("Window", "")
}

func (w *Window) IsClipped(cx, cy, cw, ch int) bool {
	if cy+ch < 0 {
		return true
	}
	if cy > w.h {
		return true
	}
	if cx+cw < 0 {
		return true
	}
	if cx > w.w {
		return true
	}
	return false
}

func traverse(buffer *bytes.Buffer, w Widget, indent int) {
	for i := 0; i < indent; i++ {
		buffer.WriteString("  ")
	}
	buffer.WriteString(w.String())
	buffer.WriteByte('\n')
	for _, c := range w.Children() {
		traverse(buffer, c, indent+1)
	}
}

func (w *Window) DebugPrint() {
	var buffer bytes.Buffer
	buffer.WriteString(w.String())
	buffer.WriteByte('\n')
	for _, c := range w.Children() {
		traverse(&buffer, c, 1)
	}
	println(buffer.String())
}
