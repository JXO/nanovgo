package wm

import (
	"github.com/goxjs/gl"
	"github.com/goxjs/glfw"
	"time"
    "runtime"
)

var debugFlag bool
var startTime time.Time
var theApp *App
var callInMain = func(f func()) {
    panic("wm.Main(...) must be called before wm.Do(f func())")
}


func init() {
    // Make sure the main goroutine is bound to the main thread.
    runtime.LockOSThread()
}

// App is a window manager and pumps events to windows
type App struct {
    name string
    active bool
    closeChan chan *Screen
}

func NewApp(name string) *App {
    if theApp != nil {

    }

	theApp = &App {
        name : name,
        closeChan : make(chan *Screen),
	}

	err := glfw.Init(gl.ContextWatcher)
	if err != nil {
		panic(err)
	}
	startTime = time.Now()

	return theApp
}

func (app *App) Run() {
    app.active = true
    defer glfw.Terminate()

	for {
        select {
        case s := <- app.closeChan:
            deleteScreen(s)

        default:
            // check window want to close
            for _, s := range uiScreens {
                if s.Window().ShouldClose() {
                    s.Send(CloseEvent{})
                    continue
                }
            }

            glfw.WaitEvents()

            if len(uiScreens) == 0 {
                return
            }
        }
	}
}


func Main(name string, main func(app *App)) {
    // Queue of functions that are thread-sensitive
    callQueue := make(chan func())

    // Properly intialize callInMain for use by xui.Do(..)
    callInMain = func(f func()) {
        done := make(chan bool, 1)
        callQueue <- func() {
            f()
            done <- true
        }
        <-done
    }

    NewApp(name)

    go func() {
        runtime.LockOSThread()
        main(theApp)
        close(callQueue)
    }()

    for f := range callQueue {
        f()
    }
}

func (app *App) Do(f func()) {
    callInMain(f)
}

func Do(f func()) {
    callInMain(f)
}


func GetTime() float32 {
	return float32(time.Now().Sub(startTime)/time.Millisecond) * 0.001
}

func SetDebug(d bool) {
	debugFlag = d
}

