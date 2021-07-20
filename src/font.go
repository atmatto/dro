package dro

import (
	_ "embed"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var F rl.Font

//go:embed Roboto.ttf
var f []byte

func LoadFont() {
	F = rl.LoadFontFromMemory(".ttf", f, int32(len(f)), 16, nil, 0)
	rl.SetTextureFilter(F.Texture, rl.FilterBilinear)
}
