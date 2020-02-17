package ui

import (
	"github.com/jxo/davinci/vg"
)

type Theme struct {
	StandardFontSize     int
	ButtonFontSize       int
	TextBoxFontSize      int
	WindowCornerRadius   int
	WindowHeaderHeight   int
	WindowDropShadowSize int
	ButtonCornerRadius   int

	DropShadow        vg.Color
	Transparent       vg.Color
	BorderDark        vg.Color
	BorderLight       vg.Color
	BorderMedium      vg.Color
	TextColor         vg.Color
	DisabledTextColor vg.Color
	TextColorShadow   vg.Color
	IconColor         vg.Color

	ButtonGradientTopFocused   vg.Color
	ButtonGradientBotFocused   vg.Color
	ButtonGradientTopUnfocused vg.Color
	ButtonGradientBotUnfocused vg.Color
	ButtonGradientTopPushed    vg.Color
	ButtonGradientBotPushed    vg.Color

	/* Window-related */
	WindowFillUnfocused  vg.Color
	WindowFillFocused    vg.Color
	WindowTitleUnfocused vg.Color
	WindowTitleFocused   vg.Color

	WindowHeaderGradientTop vg.Color
	WindowHeaderGradientBot vg.Color
	WindowHeaderSepTop      vg.Color
	WindowHeaderSepBot      vg.Color

	WindowPopup            vg.Color
	WindowPopupTransparent vg.Color

	FontNormal string
	FontBold   string
	FontIcons  string
}

func NewStandardTheme(ctx *vg.Context) *Theme {
	ctx.CreateFontFromMemory("sans", MustAsset("fonts/Roboto-Regular.ttf"), 0)
	ctx.CreateFontFromMemory("sans-bold", MustAsset("fonts/Roboto-Bold.ttf"), 0)
	ctx.CreateFontFromMemory("icons", MustAsset("fonts/entypo.ttf"), 0)
	return &Theme{
		StandardFontSize:     16,
		ButtonFontSize:       20,
		TextBoxFontSize:      20,
		WindowCornerRadius:   2,
		WindowHeaderHeight:   30,
		WindowDropShadowSize: 10,
		ButtonCornerRadius:   2,

		DropShadow:        vg.MONO(0, 128),
		Transparent:       vg.MONO(0, 0),
		BorderDark:        vg.MONO(29, 255),
		BorderLight:       vg.MONO(92, 255),
		BorderMedium:      vg.MONO(35, 255),
		TextColor:         vg.MONO(255, 160),
		DisabledTextColor: vg.MONO(255, 80),
		TextColorShadow:   vg.MONO(0, 160),
		IconColor:         vg.MONO(255, 160),

		ButtonGradientTopFocused:   vg.MONO(64, 255),
		ButtonGradientBotFocused:   vg.MONO(48, 255),
		ButtonGradientTopUnfocused: vg.MONO(74, 255),
		ButtonGradientBotUnfocused: vg.MONO(58, 255),
		ButtonGradientTopPushed:    vg.MONO(41, 255),
		ButtonGradientBotPushed:    vg.MONO(29, 255),

		WindowFillUnfocused:  vg.MONO(43, 230),
		WindowFillFocused:    vg.MONO(45, 230),
		WindowTitleUnfocused: vg.MONO(220, 160),
		WindowTitleFocused:   vg.MONO(255, 190),

		WindowHeaderGradientTop: vg.MONO(74, 255),
		WindowHeaderGradientBot: vg.MONO(58, 255),
		WindowHeaderSepTop:      vg.MONO(92, 255),
		WindowHeaderSepBot:      vg.MONO(29, 255),

		WindowPopup:            vg.MONO(50, 255),
		WindowPopupTransparent: vg.MONO(50, 0),

		FontNormal: "sans",
		FontBold:   "sans-bold",
		FontIcons:  "icons",
	}
}
