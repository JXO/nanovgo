package vg

const (
	vgInitFontImageSize = 512
	vgMaxFontImageSize  = 2048
	vgMaxFontImages     = 4

	vgInitCommandsSize = 256
	vgInitPointsSize   = 128
	vgInitPathsSize    = 16
	vgInitVertsSize    = 256
	vgMaxStates        = 32
)

type vgCommands int

const (
	vgMOVETO vgCommands = iota
	vgLINETO
	vgBEZIERTO
	vgCLOSE
	vgWINDING
)

type vgPointFlags int

const (
	vgPtCORNER     vgPointFlags = 0x01
	vgPtLEFT       vgPointFlags = 0x02
	vgPtBEVEL      vgPointFlags = 0x04
	vgPrINNERBEVEL vgPointFlags = 0x08
)

type vgTextureType int

const (
	vgTextureALPHA vgTextureType = 1
	vgTextureRGBA  vgTextureType = 2
)

type vgCodePointSize int

const (
	vgNEWLINE vgCodePointSize = iota
	vgSPACE
	vgCHAR
)
