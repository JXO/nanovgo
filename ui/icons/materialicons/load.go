package materialicons

import (
	"github.com/jxo/davinci/vg"
)

func LoadFont(ctx *vg.Context) {
	ctx.CreateFontFromMemory("materialicons", MustAsset("font/MaterialIcons-Regular.ttf"), 0)
}

func LoadFontAs(ctx *vg.Context, name string) {
	ctx.CreateFontFromMemory(name, MustAsset("font/MaterialIcons-Regular.ttf"), 0)
}
