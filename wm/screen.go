package wm

import (
	//"github.com/goxjs/gl"
	"github.com/goxjs/glfw"
	"github.com/jxo/davinci/vg"
	"runtime"
    "fmt"
)

var uiScreens map[*glfw.Window]*Screen = map[*glfw.Window]*Screen {}

type Screen struct {
    EventQueue

	w, h int
	fbW, fbH int
    visible bool

	win *glfw.Window
	ctx *vg.Context
	caption string

	cursors                [3]int
	pixelRatio             float32
	mouseState int
    modifiers glfw.ModifierKey
	mousePosX, mousePosY   int
	dragActive             bool
	lastInteraction        float32
	backgroundColor        vg.Color

    onLoadCallback func(s *Screen)
}

func deleteScreen(s *Screen) {
    delete(uiScreens, s.win)
	if s.ctx != nil {
		s.ctx.Delete()
		s.ctx = nil
	}
	if s.win != nil {
		s.win.Destroy()
		s.win = nil
	}
}

// One Screen contains one OpenGL/Event goroutine
func NewScreen(width, height int, caption string, resizable, fullScreen bool) *Screen {
    if theApp == nil {
        panic("NewScreen(...) should be called after NewApp(...)")
    }

    if caption == "" {
        caption = theApp.name
    }

    s := &Screen{
        //cursor:  glfw.CursorNormal,
        caption: caption,
    }

    if runtime.GOARCH == "js" {
        glfw.WindowHint(glfw.Hint(0x00021101), 1) // enable stencil for vg
    }

    glfw.WindowHint(glfw.Samples, 4)
    //glfw.WindowHint(glfw.RedBits, 8)
    //glfw.WindowHint(glfw.GreenBits, 8)
    //glfw.WindowHint(glfw.BlueBits, 8)
    glfw.WindowHint(glfw.AlphaBits, 8)
    //glfw.WindowHint(glfw.StencilBits, 8)
    //glfw.WindowHint(glfw.DepthBits, 8)
    //glfw.WindowHint(glfw.Visible, 0)
    if resizable {
        glfw.WindowHint(glfw.Resizable, 1)
    } else {
        glfw.WindowHint(glfw.Resizable, 0)
    }

    var err error
    if fullScreen {
        monitor := glfw.GetPrimaryMonitor()
        mode := monitor.GetVideoMode()
        s.win, err = glfw.CreateWindow(mode.Width, mode.Height, caption, monitor, nil)
    } else {
        s.win, err = glfw.CreateWindow(width, height, caption, nil, nil)
    }

    if err != nil {
        panic(err)
    }

	uiScreens[s.win] = s
	s.w, s.h = s.win.GetSize()
	s.fbW, s.fbH = s.win.GetFramebufferSize()
	s.visible = true //s.win.GetAttrib(glfw.Visible)
//	s.stheme = NewStandardTheme()

	s.lastInteraction = GetTime()
    //XXX
    /*
    s.win.SetScreenRefreshCallback(func(w *glfw.Window) {
        if s, ok := uiScreens[w]; ok {
            s.Send(PaintEvent{})
        }
    }) */

    s.win.SetCursorPosCallback(func(w *glfw.Window, xpos, ypos float64) {
        if s, ok := uiScreens[w]; ok {
            s.Send(MouseMoveEvent{
                X : xpos,
                Y : ypos,
            })
        }
    })

    s.win.SetMouseButtonCallback(func(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
        if s, ok := uiScreens[w]; ok {
            s.Send(MouseButtonEvent{
                Button: button,
                Action: action,
                Modifier: mods,
            })
        }
    })

    s.win.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scanCode int, action glfw.Action, mods glfw.ModifierKey) {
        if s, ok := uiScreens[w]; ok {
            s.Send(KeyEvent{
                Key : key,
                ScanCode : scanCode,
                Action: action,
                Modifier : mods,
            })
        }
    })

    s.win.SetCharCallback(func(w *glfw.Window, r rune) {
        if s, ok := uiScreens[w]; ok {
            s.Send(RuneEvent {
                R : r,
            })
        }
    })

    s.win.SetDropCallback(func(w *glfw.Window, names []string) {
        if s, ok := uiScreens[w]; ok {
            s.Send(DropEvent {
                Names: names,
            })
        }
    })

    s.win.SetScrollCallback(func(w *glfw.Window, xoff float64, yoff float64) {
        if s, ok := uiScreens[w]; ok {
            e := ScrollEvent {
                Dx : xoff,
                Dy : yoff,
            }
            s.Send(e)
        }
    })

    s.win.SetFramebufferSizeCallback(func(w *glfw.Window, width int, height int) {
        if s, ok := uiScreens[w]; ok {
            s.Send(ResizeEvent {
                Width: width,
                Height: height,
            })
        }
    })

    // OpenGL and Events processing
    go func() {
        runtime.LockOSThread()
        s.initializeGfx()

        if debugFlag {
            fmt.Println("OpenGL up")
        }

        if s.onLoadCallback != nil {
            s.onLoadCallback(s)
        }

        // Screen Event Processing
        for {
            switch e := s.NextEvent().(type) {
            case PaintEvent:
                s.OnPaint()
            case MouseMoveEvent:
                s.OnMouseMove(&e)
            case MouseButtonEvent:
                s.OnMouseButton(&e)
            case KeyEvent:
                s.OnKeyEvent(&e)
            case DropEvent:
                s.OnDropEvent(&e)
            case ScrollEvent:
                s.OnScrollEvent(&e)
            case ResizeEvent:
                s.OnResizeEvent(&e)
            case CloseEvent:
                if s.OnCloseEvent() {
                    glfw.PostEmptyEvent()
                    theApp.closeChan <- s
                    return
                }
            default:
                // do nothing
            }
        }
    }()

    return s
}

func (s *Screen) OnPaint() {
    fmt.Println("Paint")
}

func (s *Screen) OnMouseMove(e *MouseMoveEvent) {
    fmt.Println("MouseMove")
}

func (s *Screen) OnMouseButton(e *MouseButtonEvent) {
    fmt.Println("MouseButton")
}

func (s *Screen) OnKeyEvent (e *KeyEvent) {
    fmt.Println("Key")
}

func (s *Screen) OnDropEvent(e *DropEvent) {
    fmt.Println("Drop")
}

func (s *Screen) OnScrollEvent(e *ScrollEvent) {
    fmt.Println("Scroll")
}

func (s *Screen) OnResizeEvent(e *ResizeEvent) {
    fmt.Println("resize")
}

func (s *Screen) OnCloseEvent() bool {
    return true
}

func (s *Screen) Window() *glfw.Window {
    return s.win
}

func (s *Screen) MakeCurrent() {
    s.win.MakeContextCurrent()
}

func (s *Screen) initializeGfx() {
    s.MakeCurrent()
	var err error
	s.ctx, err = vg.NewContext(vg.StencilStrokes | vg.AntiAlias)
	if err != nil {
		panic(err)
	}
}
