package dro

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

// TODO : Switch to HSV picker

type Picker struct {
	H float32 // Hue [0..360]
	S float32 // Saturation [0..1]
	V float32 // Value [0..1]
	A float32 // Alpha [0..1]
}

func ColorToPicker(C rl.Color) Picker {
	hsv := rl.ColorToHSV(C)
	return Picker{
		H: hsv.X,
		S: hsv.Y,
		V: hsv.Z,
		A: float32(C.A) / 255,
	}
}

// Update skips fields less than zero
func (p *Picker) Update(h float32, s float32, v float32, a float32) {
	if h >= 0 {
		p.H = h
	}
	if s >= 0 {
		p.S = s
	}
	if v >= 0 {
		p.V = v
	}
	if a >= 0 {
		p.A = a
	}
}

func (p Picker) ToColor() rl.Color {
	return rl.ColorAlpha(rl.ColorFromHSV(p.H, p.S, p.V), p.A)
}

var (
	oldFg       rl.Color
	oldBg       rl.Color
	newFg       rl.Color
	newBg       rl.Color
	picker      Picker
	currentCtrl = "" // Which control is currently used (when lmb is held down)
)

func OpenColorDialog(S *State) *State {
	if S.Mode == M_Dialog && S.Dialog == D_Color {
		return CloseColorDialog(S)
	}
	if S.Mode != M_Dialog {
		S.LastMode = S.Mode
	}
	S.Mode = M_Dialog
	S.Dialog = D_Color
	oldFg = S.Brush.ForegroundColor
	oldBg = S.Brush.BackgroundColor
	newFg = oldFg
	newBg = oldBg
	picker = ColorToPicker(newFg)
	currentCtrl = "-"
	return S
}

func CloseColorDialog(S *State) *State {
	S.Brush.ForegroundColor = newFg
	S.Brush.BackgroundColor = newBg
	S.Mode = S.LastMode
	S.Dialog = 0
	return S
}

func DoColorDialog(S *State) *State {
	var top int32 = 35
	var left int32 = 15
	var theight int32 = int32(rl.GetScreenHeight() - 70)
	var twidth int32 = int32(rl.GetScreenWidth() - 30)

	if theight > 400 {
		theight = 400
	}

	height := theight - 15 - 30 - 10
	width := twidth - 2*15 - 2*30
	if width > height {
		width = height
	} else {
		height = width
	}

	// SL
	rl.DrawRectangleLines(left-1, top-1, width+2, height+2, rl.Black)
	for x := 0; int32(x) < width; x++ {
		for y := 0; int32(y) < height; y++ {
			rl.DrawPixel(int32(x)+left, int32(y)+top, rl.ColorFromHSV(picker.H, float32(x)/float32(width), 1-(float32(y)/float32(height))))
		}
	}
	if rl.CheckCollisionPointRec(rl.GetMousePosition(), rl.NewRectangle(float32(left), float32(top), float32(width), float32(height))) && rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		currentCtrl = "sl"
	}
	if rl.IsMouseButtonDown(rl.MouseLeftButton) && currentCtrl == "sl" {
		s := float32(rl.GetMouseX()-left) / float32(width)
		v := float32(rl.GetMouseY()-top) / float32(height)
		s = rl.Clamp(s, 0, 1)
		v = rl.Clamp(v, 0, 1)
		picker.Update(-1, s, 1-v, -1)
	}

	rl.DrawRectangleLines(int32(float32(picker.S)*float32(width))+left-6, int32(float32(1-picker.V)*float32(height))+top-6, 12, 12, rl.Black)
	rl.DrawRectangleLines(int32(float32(picker.S)*float32(width))+left-5, int32(float32(1-picker.V)*float32(height))+top-5, 10, 10, rl.White)

	// H
	rl.DrawRectangleLines(left+width+14, top-1, 32, height+2, rl.Black)
	for i := 1; int32(i) <= height; i++ {
		rl.DrawLine(left+width+15, top+int32(i), left+width+15+30, top+int32(i), rl.ColorFromHSV(float32(i)/float32(height)*360, 1, 1))
	}
	if rl.CheckCollisionPointRec(rl.GetMousePosition(), rl.NewRectangle(float32(left+width+15), float32(top), float32(30), float32(height))) && rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		currentCtrl = "h"
	}
	if rl.IsMouseButtonDown(rl.MouseLeftButton) && currentCtrl == "h" {
		h := float32(rl.GetMouseY()-top) / float32(height) * 360
		h = rl.Clamp(h, 0, 360)
		picker.Update(h, -1, -1, -1)
	}
	rl.DrawRectangleLines(left+width+14-5, int32(float32(picker.H)/360.0*float32(height))+top-6, 40, 12, rl.Black)
	rl.DrawRectangleLines(left+width+14-4, int32(float32(picker.H)/360.0*float32(height))+top-5, 38, 10, rl.White)

	// A
	rl.DrawRectangle(left+width+15+15+30, top, 15, height, rl.White)
	rl.DrawRectangle(left+width+15+30+30, top, 15, height, rl.Black)
	rl.DrawRectangleLines(left+width+14+15+30, top-1, 32, height+2, rl.Black)
	for i := 1; int32(i) <= height; i++ {
		rl.DrawLine(left+width+15+15+30, top+int32(i), left+width+15+15+30+30, top+int32(i), rl.ColorAlpha(picker.ToColor(), float32(i)/float32(height)))
	}
	if rl.CheckCollisionPointRec(rl.GetMousePosition(), rl.NewRectangle(float32(left+width+15+15+30), float32(top), float32(30), float32(height))) && rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		currentCtrl = "a"
	}
	if rl.IsMouseButtonDown(rl.MouseLeftButton) && currentCtrl == "a" {
		a := float32(rl.GetMouseY()-top) / float32(height)
		a = rl.Clamp(a, 0, 1)
		picker.Update(-1, -1, -1, a)
	}
	rl.DrawRectangleLines(left+width+15+15+30-5, int32(float32(picker.A)*float32(height))+top-6, 40, 12, rl.Black)
	rl.DrawRectangleLines(left+width+15+15+30-4, int32(float32(picker.A)*float32(height))+top-5, 38, 10, rl.White)

	rl.DrawTextEx(F, "Saturation, Lightness", rl.NewVector2(float32(left), float32(top+height+2)), 16, 1, rl.Black)
	rl.DrawTextEx(F, "Hue", rl.NewVector2(float32(left+width+15), float32(top+height+2)), 16, 1, rl.Black)
	rl.DrawTextEx(F, "Alpha", rl.NewVector2(float32(left+width+15+15+30), float32(top+height+2)), 16, 1, rl.Black)

	rl.DrawTextEx(F, "OLD FG", rl.NewVector2(float32(left), float32(top+height+40)), 16, 1, rl.Black)
	l := int32(rl.MeasureTextEx(F, "OLD FG", 16, 1).X + 10)
	rl.DrawRectangle(left+l, top+height+25, 60, 28, oldFg)
	rl.DrawRectangleLines(left+l-1, top+height+24, 62, 30, rl.Black)
	if rl.CheckCollisionPointRec(rl.GetMousePosition(), rl.NewRectangle(float32(left+l), float32(top+height+25), float32(60), float32(28))) {
		if currentCtrl == "" {
			rl.DrawRectangleLines(left+l-1, top+height+24, 62, 30, rl.White)
		}
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			currentCtrl = "set"
			picker = ColorToPicker(oldFg)
		}
	}
	l += 75
	rl.DrawTextEx(F, "BG", rl.NewVector2(float32(left+l), float32(top+height+40)), 16, 1, rl.Black)
	l += int32(rl.MeasureTextEx(F, "BG", 16, 1).X + 10)
	rl.DrawRectangle(left+l, top+height+25, 60, 28, oldBg)
	rl.DrawRectangleLines(left+l-1, top+height+24, 62, 30, rl.Black)
	if rl.CheckCollisionPointRec(rl.GetMousePosition(), rl.NewRectangle(float32(left+l), float32(top+height+25), float32(60), float32(28))) {
		if currentCtrl == "" {
			rl.DrawRectangleLines(left+l-1, top+height+24, 62, 30, rl.White)
		}
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			currentCtrl = "set"
			picker = ColorToPicker(oldBg)
		}
	}
	l += 80

	rl.DrawTextEx(F, "NEW FG", rl.NewVector2(float32(left+l), float32(top+height+40)), 16, 1, rl.Black)
	l += int32(rl.MeasureTextEx(F, "NEW FG", 16, 1).X + 10)
	rl.DrawRectangle(left+l, top+height+25, 60, 28, newFg)
	rl.DrawRectangleLines(left+l-1, top+height+24, 62, 30, rl.Black)
	if rl.CheckCollisionPointRec(rl.GetMousePosition(), rl.NewRectangle(float32(left+l), float32(top+height+25), float32(60), float32(28))) {
		if currentCtrl == "" {
			rl.DrawRectangle(left+l, top+height+25, 60, 28, picker.ToColor())
			rl.DrawRectangleLines(left+l-1, top+height+24, 62, 30, rl.White)
		}
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			currentCtrl = "set"
			newFg = picker.ToColor()
		}
	}
	l += 75
	rl.DrawTextEx(F, "BG", rl.NewVector2(float32(left+l), float32(top+height+40)), 16, 1, rl.Black)
	l += int32(rl.MeasureTextEx(F, "BG", 16, 1).X + 10)
	rl.DrawRectangle(left+l, top+height+25, 60, 28, newBg)
	rl.DrawRectangleLines(left+l-1, top+height+24, 62, 30, rl.Black)
	if rl.CheckCollisionPointRec(rl.GetMousePosition(), rl.NewRectangle(float32(left+l), float32(top+height+25), float32(60), float32(28))) {
		if currentCtrl == "" {
			rl.DrawRectangle(left+l, top+height+25, 60, 28, picker.ToColor())
			rl.DrawRectangleLines(left+l-1, top+height+24, 62, 30, rl.White)
		}
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			currentCtrl = "set"
			newBg = picker.ToColor()
		}
	}

	if rl.IsMouseButtonPressed(rl.MouseLeftButton) && currentCtrl == "" {
		return CloseColorDialog(S)
	}

	if rl.IsMouseButtonReleased(rl.MouseLeftButton) {
		currentCtrl = ""
	}

	return S
}
