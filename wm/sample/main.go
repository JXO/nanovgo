package main

import (
    "github.com/jxo/davinci/wm"
)

func run(app *wm.App) {
    app.Do(func () {
        wm.NewScreen(800, 600, "", true, false)

        app.Run()
    })
}

func main() {
    wm.Main("test", run)
}
