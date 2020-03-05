// +build !js

package main

import (
	"fmt"
	"github.com/jxo/davinci/ui"
    "github.com/jxo/davinci/gfx"
	"github.com/jxo/davinci/ui/sample1/demo"
	"github.com/jxo/davinci/vg"
	"io/ioutil"
	"math"
	"path"
)

func run() {
    var app *ui.App

    gfx.Do(func() {
        app = ui.NewApp("DavinciUI Test")
        w := app.NewWindow(1024, 768, "", true, false)

        var images []ui.Image
        w.SetResizeEventCallback(func(x, y int) bool {
            w.OnPaint()
            return true
        })

        w.SetOnLoad(func() {
            w.NVGContext().CreateFont("japanese", "font/GenShinGothic-P-Regular.ttf")
            images = loadImageDirectory(w.NVGContext(), "icons")
            demo.ButtonDemo(w)
            imageButton, imagePanel, progressBar := demo.BasicWidgetsDemo(w, images)
            progress := progressBar
            demo.SelectedImageDemo(w, imageButton, imagePanel)
            demo.MiscWidgetsDemo(w)
            demo.GridDemo(w)

            w.SetDrawContentsCallback(func() {
                progress.SetValue(float32(math.Mod(float64(ui.GetTime())/10, 1.0)))
            })

            w.PerformLayout()
	        w.DebugPrint()
        })
    })

    gfx.Do(app.Run)
}

func main() {
    gfx.Main(run)
}

func loadImageDirectory(ctx *vg.Context, dir string) []ui.Image {
	var images []ui.Image
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(fmt.Sprintf("loadImageDirectory: read error %v\n", err))
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		ext := path.Ext(file.Name())
		if ext != ".png" {
			continue
		}
		fullPath := path.Join(dir, file.Name())
		img := ctx.CreateImage(fullPath, 0)
		if img == 0 {
			panic("Could not open image data!")
		}
		images = append(images, ui.Image{
			ImageID: img,
			Name:    fullPath[:len(fullPath)-4],
		})
	}
	return images
}
