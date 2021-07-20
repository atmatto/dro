package dro

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Brush struct {
	ForegroundColor rl.Color
	BackgroundColor rl.Color
	Thickness       float32
}

// Returns true if mouse moved enough distance to draw new stroke
func (B *Brush) CanDraw(S *State) bool {
	return S.CanGetMouseMovement() && (rl.Vector2Distance(S.MouseOld, rl.GetMousePosition()) > B.Thickness*0.1 || rl.IsMouseButtonPressed(rl.MouseLeftButton))
}

func (B *Brush) DoDrawing(S *State) *State {
	if B.CanDraw(S) {
		if rl.IsMouseButtonDown(rl.MouseLeftButton) {
			rl.DrawLineEx(Vector2Divide(rl.Vector2Subtract(S.MouseOld, S.Canvas.Position), S.Canvas.Scale), Vector2Divide(rl.Vector2Subtract(rl.GetMousePosition(), S.Canvas.Position), S.Canvas.Scale), B.Thickness, B.ForegroundColor)
			rl.DrawCircleV(Vector2Divide(rl.Vector2Subtract(rl.GetMousePosition(), S.Canvas.Position), S.Canvas.Scale), B.Thickness/2, B.ForegroundColor)
			rl.DrawCircleV(Vector2Divide(rl.Vector2Subtract(S.MouseOld, S.Canvas.Position), S.Canvas.Scale), B.Thickness/2, B.ForegroundColor)
			return S.Change(false)
		}
	}

	return S
}

func (B *Brush) DoThickness() {
	if rl.IsKeyPressed(rl.KeyLeftBracket) {
		if rl.IsKeyDown(rl.KeyLeftControl) {
			B.Thickness -= 5
		} else {
			B.Thickness -= 1
		}
	} else if rl.IsKeyPressed(rl.KeyRightBracket) {
		if rl.IsKeyDown(rl.KeyLeftControl) {
			B.Thickness += 5
		} else {
			B.Thickness += 1
		}
	}
	if B.Thickness < 1 {
		B.Thickness = 1
	}
}
