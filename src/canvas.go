package dro

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Canvas struct {
	Tex           *rl.RenderTexture2D
	Position      rl.Vector2
	Scale         float32
	DisplayedSize rl.Vector2
	Size          rl.Vector2
}

func (C *Canvas) Load(t *rl.Texture2D) {
	if t != nil {
		C.Size = rl.Vector2{X: float32(t.Width), Y: float32(t.Height)}
		C.DisplayedSize = C.Size
		tex := rl.LoadRenderTexture(int32(C.Size.X), int32(C.Size.Y))
		C.Tex = &tex
		rl.BeginTextureMode(*C.Tex)
		rl.DrawTexture(*t, 0, 0, rl.White)
		rl.EndTextureMode()
		rl.UnloadTexture(*t)
	} else {
		C.DisplayedSize = rl.NewVector2(C.Size.X*C.Scale, C.Size.Y*C.Scale)
		tex := rl.LoadRenderTexture(int32(C.Size.X), int32(C.Size.Y))
		C.Tex = &tex
		rl.BeginTextureMode(*C.Tex)
		rl.ClearBackground(rl.White)
		rl.EndTextureMode()
	}
}

func (C *Canvas) Draw() {
	rl.DrawTexturePro(C.Tex.Texture, rl.NewRectangle(0, float32(C.Tex.Texture.Height), float32(C.Tex.Texture.Width), -float32(C.Tex.Texture.Height)),
		rl.NewRectangle(C.Position.X, C.Position.Y, C.DisplayedSize.X, C.DisplayedSize.Y), rl.NewVector2(0, 0), 0, rl.White)
}

func (C *Canvas) Pan(s *State) {
	if rl.IsMouseButtonPressed(2) && s.CanGetMouseMovement() {
		s.MouseOld = rl.GetMousePosition()
	} else if rl.IsMouseButtonDown(2) && s.CanGetMouseMovement() {
		C.Position = rl.Vector2Add(C.Position, rl.Vector2Subtract(rl.GetMousePosition(), s.MouseOld))
		s.MouseOld = rl.GetMousePosition()
	}
}

func (C *Canvas) DoScale() {
	if rl.GetMouseWheelMove() > 0 {
		C.Scale *= 2
	}
	if rl.GetMouseWheelMove() < 0 {
		C.Scale *= 0.5
	}
	if rl.GetMouseWheelMove() != 0 {
		C.Position = rl.NewVector2(
			float32(rl.GetMouseX())-(float32(rl.GetMouseX())-C.Position.X)*C.Scale*C.Size.X/C.DisplayedSize.X,
			float32(rl.GetMouseY())-(float32(rl.GetMouseY())-C.Position.Y)*C.Scale*C.Size.Y/C.DisplayedSize.Y)
		C.DisplayedSize.X = C.Size.X * C.Scale
		C.DisplayedSize.Y = C.Size.Y * C.Scale
	}
}

func (C *Canvas) ScreenToCanvas(v rl.Vector2) rl.Vector2 {
	return Vector2Divide(rl.Vector2Subtract(v, C.Position), C.Scale)
}

func (C *Canvas) CanvasToScreen(v rl.Vector2) rl.Vector2 {
	return rl.Vector2Add(rl.Vector2Multiply(v, rl.NewVector2(C.Scale, C.Scale)), C.Position)
}
