package wm

import (
	"github.com/goxjs/glfw"
    "time"
)

type PaintEvent struct { }

type MouseMoveEvent struct {
    X, Y float64
}

type MouseButtonEvent struct {
    Button   glfw.MouseButton
    Action   glfw.Action
    Modifier glfw.ModifierKey
}

type RuneEvent struct {
    R rune
}

type KeyEvent struct {
    Key glfw.Key
    ScanCode int
    Action glfw.Action
    Modifier glfw.ModifierKey
}

type DropEvent struct {
    Names []string
}

type ScrollEvent struct {
    Dx, Dy float64
}

type ResizeEvent struct {
    Width, Height int
}

type CloseEvent struct { }

type GLEvent struct {
    T time.Time
    F func(t time.Time)
}
