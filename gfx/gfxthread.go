package gfx

import (
    "runtime"
)

var callInMain = func(f func()) {
    panic("xui.Main(main func()) must be called before xui.Do(f func())")
}

func init() {
    // Make sure the main goroutine is bound to the main thread.
    runtime.LockOSThread()
}

func Main(run func()) {
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

    go func() {
        runtime.LockOSThread()
        run()
        close(callQueue)
    }()

    for f := range callQueue {
        f()
    }
}

func Do(f func()) {
    callInMain(f)
}
