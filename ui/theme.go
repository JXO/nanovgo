package ui

import (
	"github.com/jxo/davinci/vg"
)

type Theme struct {
	StandardFontSize     int
	ButtonFontSize       int
	TextBoxFontSize      int
	PanelCornerRadius   int
	PanelHeaderHeight   int
	PanelDropShadowSize int
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

	/* Panel-related */
	PanelFillUnfocused  vg.Color
	PanelFillFocused    vg.Color
	PanelTitleUnfocused vg.Color
	PanelTitleFocused   vg.Color

	PanelHeaderGradientTop vg.Color
	PanelHeaderGradientBot vg.Color
	PanelHeaderSepTop      vg.Color
	PanelHeaderSepBot      vg.Color

	PanelPopup            vg.Color
	PanelPopupTransparent vg.Color

	FontNormal string
	FontBold   string
	FontIcons  string
}

func NewStandardTheme() *Theme {
	return &Theme{
		StandardFontSize:     16,
		ButtonFontSize:       20,
		TextBoxFontSize:      20,
		PanelCornerRadius:   2,
		PanelHeaderHeight:   30,
		PanelDropShadowSize: 10,
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

		PanelFillUnfocused:  vg.MONO(43, 230),
		PanelFillFocused:    vg.MONO(45, 230),
		PanelTitleUnfocused: vg.MONO(220, 160),
		PanelTitleFocused:   vg.MONO(255, 190),

		PanelHeaderGradientTop: vg.MONO(74, 255),
		PanelHeaderGradientBot: vg.MONO(58, 255),
		PanelHeaderSepTop:      vg.MONO(92, 255),
		PanelHeaderSepBot:      vg.MONO(29, 255),

		PanelPopup:            vg.MONO(50, 255),
		PanelPopupTransparent: vg.MONO(50, 0),

		FontNormal: "sans",
		FontBold:   "sans-bold",
		FontIcons:  "icons",
	}
}
