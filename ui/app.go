package ui

import (
	"github.com/goxjs/gl"
	"github.com/goxjs/glfw"
	"runtime"
    "runtime/pprof"
    "os"
	"time"
    "fmt"
)

var debugFlag bool
var uiWindows map[*glfw.Window]*Window = map[*glfw.Window]*Window {}

var startTime time.Time

// App is a window manager and pumps events to windows
type App struct {
    name string
    active bool
    closeChan chan *Window
}

func NewApp(name string) *App {
	app := &App {
        name : name,
        closeChan : make(chan *Window),
	}

	err := glfw.Init(gl.ContextWatcher)
	if err != nil {
		panic(err)
	}
	startTime = time.Now()

	return app
}

func deleteWindow(w *Window) {
    delete(uiWindows, w.win)
	if w.context != nil {
		w.context.Delete()
		w.context = nil
	}
	if w.win != nil {
		w.win.Destroy()
		w.win = nil
	}
}


// One Window contains one OpenGL render goroutine and event goroutine
func (app *App) NewWindow(width, height int, caption string, resizable, fullScreen bool) *Window {
    if caption == "" {
        caption = app.name
    }

    w := &Window{
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
        w.win, err = glfw.CreateWindow(mode.Width, mode.Height, caption, monitor, nil)
    } else {
        w.win, err = glfw.CreateWindow(width, height, caption, nil, nil)
    }

    if err != nil {
        panic(err)
    }

	uiWindows[w.win] = w
	w.w, w.h = w.win.GetSize()
	w.fbW, w.fbH = w.win.GetFramebufferSize()
	w.visible = true //panel.GetAttrib(glfw.Visible)
	w.theme = NewStandardTheme()

	w.lastInteraction = GetTime()
    //XXX
    /*
    w.win.SetWindowRefreshCallback(func(w *glfw.Window) {
        if w, ok := uiWindows[w]; ok {
            w.Send(PaintEvent{})
        }
    }) */

    w.win.SetCursorPosCallback(func(w *glfw.Window, xpos, ypos float64) {
        if w, ok := uiWindows[w]; ok {
            w.Send(MouseMoveEvent{
                X : xpos,
                Y : ypos,
            })
        }
    })

    w.win.SetMouseButtonCallback(func(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
        if w, ok := uiWindows[w]; ok {
            w.Send(MouseButtonEvent{
                Button: button,
                Action: action,
                Modifier: mods,
            })
        }
    })

    w.win.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scanCode int, action glfw.Action, mods glfw.ModifierKey) {
        if w, ok := uiWindows[w]; ok {
            w.Send(KeyEvent{
                Key : key,
                ScanCode : scanCode,
                Action: action,
                Modifier : mods,
            })
        }
    })

    w.win.SetCharCallback(func(w *glfw.Window, r rune) {
        if w, ok := uiWindows[w]; ok {
            w.Send(RuneEvent {
                R : r,
            })
        }
    })

    w.win.SetDropCallback(func(w *glfw.Window, names []string) {
        if w, ok := uiWindows[w]; ok {
            e := DropEvent {
                Names: names,
            }
            w.Send(e)
        }
    })

    w.win.SetScrollCallback(func(w *glfw.Window, xoff float64, yoff float64) {
        if w, ok := uiWindows[w]; ok {
            e := ScrollEvent {
                Dx : xoff,
                Dy : yoff,
            }
            w.Send(e)
        }
    })

    w.win.SetFramebufferSizeCallback(func(w *glfw.Window, width int, height int) {
        if w, ok := uiWindows[w]; ok {
            e := ResizeEvent {
                Width: width,
                Height: height,
            }
            w.Send(e)
        }
    })

    // OpenGL
    done := make(chan struct {})
    go func() {
        runtime.LockOSThread()
        w.InitializeGfx()

        f, err := os.Create("cpu.prof")
        if err != nil {
            fmt.Println("Profile create error")
            return
        }

        done <- struct {}{}

        fmt.Println("OpenGL up")
        pprof.StartCPUProfile(f)
        defer pprof.StopCPUProfile()
        for {
            switch e := w.renderQueue.NextEvent().(type) {
            case GLEvent:
                e.F(e.T)
            case func():
                e()
            case CloseEvent:
                fmt.Println("OpenGL Quit")
                return
            }
        }
    } ()

    <- done

    // Events processing
    go func() {
        for {
            switch e := w.NextEvent().(type) {
            case PaintEvent:
                w.OnPaint()
            case MouseMoveEvent:
                w.OnMouseMove(&e)
            case MouseButtonEvent:
                w.OnMouseButton(&e)
            case KeyEvent:
                w.OnKeyEvent(&e)
            case DropEvent:
                w.OnDropEvent(&e)
            case ScrollEvent:
                w.OnScrollEvent(&e)
            case ResizeEvent:
                w.OnResizeEvent(&e)
            case CloseEvent:
                if w.OnCloseEvent() {
                    w.renderQueue.Send(CloseEvent{})
                    glfw.PostEmptyEvent()
                    app.closeChan <- w
                    return
                }
            default:
                // do nothing
            }
        }
    }()

    return w
}

func (app *App) Run() {
    app.active = true
    defer glfw.Terminate()

	for {
        select {
        case w := <- app.closeChan:
            deleteWindow(w)

        default:
            // check window want to close
            for _, w := range uiWindows {
                if w.GLFWWindow().ShouldClose() {
                    w.Send(CloseEvent{})
                    continue
                }
            }

            glfw.WaitEvents()

            if len(uiWindows) == 0 {
                return
            }
        }
	}
}

func GetTime() float32 {
	return float32(time.Now().Sub(startTime)/time.Millisecond) * 0.001
}

func SetDebug(d bool) {
	debugFlag = d
}

func InitWidget(child, parent Widget) {
	//w.cursor = Arrow
	if parent != nil {
		parent.AddChild(parent, child)
		child.SetTheme(parent.Theme())
	}
	child.SetVisible(true)
	child.SetEnabled(true)
	child.SetFontSize(-1)
}
