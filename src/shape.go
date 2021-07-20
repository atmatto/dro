package dro

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	S_Line = iota
	S_Rectangle
	S_RectangleLines
	S_Ellipse
	S_EllipseLines
	ShapeNumber
)

type dRectangle struct {
	X      float32
	Y      float32
	Width  float32
	Height float32
	DirX   int
	DirY   int
}

// Normalize will make sure width and height aren't less than zero.
func (r dRectangle) Normalize() dRectangle {
	if r.Width < 0 {
		r.Width *= -1
		r.X -= r.Width
	}
	if r.Height < 0 {
		r.Height *= -1
		r.Y -= r.Height
	}
	return r
}

func NewdRectangle(start rl.Vector2, end rl.Vector2) dRectangle {
	dx := 0
	dy := 0
	if start.X < end.X {
		dx = 1
	} else {
		dx = -1
	}
	if start.Y < end.Y {
		dy = 1
	} else {
		dy = -1
	}
	r := dRectangle{
		X:      start.X,
		Y:      start.Y,
		Width:  end.X - start.X,
		Height: end.Y - start.Y,
		DirX:   dx,
		DirY:   dy,
	}

	return r.Normalize()
}

func (r dRectangle) TestDirection() {
	switch r.DirX {
	case 1:
		if r.DirY == 1 {
			fmt.Println("⭨")
			NewdRectangle(rl.NewVector2(30, 30), rl.NewVector2(60, 60)).DrawShape(S_Line, 1, rl.Black)
		} else {
			fmt.Println("⭧")
			NewdRectangle(rl.NewVector2(30, 60), rl.NewVector2(60, 30)).DrawShape(S_Line, 1, rl.Black)
		}
	case -1:
		if r.DirY == 1 {
			fmt.Println("⭩")
			NewdRectangle(rl.NewVector2(60, 30), rl.NewVector2(30, 60)).DrawShape(S_Line, 1, rl.Black)
		} else {
			fmt.Println("⭦")
			NewdRectangle(rl.NewVector2(60, 60), rl.NewVector2(30, 30)).DrawShape(S_Line, 1, rl.Black)
		}
	}
}

// Returns start and end positions
func (r dRectangle) Vectors() (rl.Vector2, rl.Vector2) {
	switch r.DirX {
	case 1:
		if r.DirY == 1 {
			// ⭨
			return rl.NewVector2(r.X, r.Y), rl.NewVector2(r.X+r.Width, r.Y+r.Height)
		} else {
			// ⭧
			return rl.NewVector2(r.X, r.Y+r.Height), rl.NewVector2(r.X+r.Width, r.Y)
		}
	case -1:
		if r.DirY == 1 {
			// ⭩
			return rl.NewVector2(r.X+r.Width, r.Y), rl.NewVector2(r.X, r.Y+r.Height)
		} else {
			// ⭦
			return rl.NewVector2(r.X+r.Width, r.Y+r.Height), rl.NewVector2(r.X, r.Y)
		}
	}
	panic("Couldn't get dRectangle vectors")
}

func (r dRectangle) Rec() rl.Rectangle {
	return rl.NewRectangle(r.X, r.Y, r.Width, r.Height)
}

func (r dRectangle) DrawShape(shape int, thick float32, color rl.Color) {
	switch shape {
	case S_Line:
		s, e := r.Vectors()
		rl.DrawEllipse(int32(s.X), int32(s.Y), thick/2, thick/2, color)
		rl.DrawEllipse(int32(e.X), int32(e.Y), thick/2, thick/2, color)
		rl.DrawLineEx(s, e, thick, color)
	case S_Rectangle:
		rl.DrawRectangleRec(r.Rec(), color)
	case S_RectangleLines:
		rl.DrawRectangleLinesEx(r.Rec(), thick, color)
	case S_Ellipse:
		rl.DrawEllipse(int32(r.X+(r.Width/2)), int32(r.Y+(r.Height/2)), r.Width/2, r.Height/2, color)
	case S_EllipseLines:
		// TODO: ellipse lines thickness
		rl.DrawEllipseLines(int32(r.X+(r.Width/2)), int32(r.Y+(r.Height/2)), r.Width/2, r.Height/2, color)
	}
}

func (r dRectangle) DrawRenderTex(tex rl.RenderTexture2D) {
	//fmt.Println(rl.NewRectangle(0, float32(tex.Texture.Height), float32(tex.Texture.Width), -float32(tex.Texture.Height)))
	trec := r.Rec()
	//rl.DrawRectangleRec(trec, rl.Black)
	//rl.DrawTexture(tex.Texture, int32(trec.X), int32(trec.Y), rl.White)
	rl.DrawTexturePro(tex.Texture, rl.NewRectangle(0, float32(tex.Texture.Height), float32(tex.Texture.Width), -float32(tex.Texture.Height)), trec, rl.NewVector2(0, 0), 0, rl.White)
}

func (r dRectangle) CanvasToScreen(c *Canvas) dRectangle {
	s, e := r.Vectors()
	s = c.CanvasToScreen(s)
	e = c.CanvasToScreen(e)
	return NewdRectangle(s, e)
}

func (r dRectangle) ScreenToCanvas(c *Canvas) dRectangle {
	s, e := r.Vectors()
	s = c.ScreenToCanvas(s)
	e = c.ScreenToCanvas(e)
	return NewdRectangle(s, e)
}

func (r dRectangle) Offset(dx float32, dy float32, dw float32, dh float32) dRectangle {
	r.X += dx
	r.Y += dy
	r.Width += dw
	r.Height += dh
	return r
}

func OpenShapeSelector(S *State) *State {
	if S.Mode == M_Dialog && S.Dialog == D_Shapes {
		return CloseShapeSelector(S, -1)
	}
	if S.Mode != M_Dialog {
		S.LastMode = S.Mode
	}
	S.Mode = M_Dialog
	S.Dialog = D_Shapes
	return S
}

// If shape == -1 then no shape was selected
func CloseShapeSelector(S *State, shape int) *State {
	if shape != -1 {
		S.Mode = M_Selection | M_Shape
		S.Shape = shape
		S.Dialog = 0
		S.MouseOld = rl.NewVector2(-1, -1)
		return S
	}
	S.Mode = S.LastMode
	S.Dialog = 0
	S.TimeoutT = rl.GetTime()
	return S
}

const (
	top          = 20
	outerMargin  = 20
	innerMargin  = 10
	selectorSize = 60
)

func AssignSelector(n int) rl.Rectangle {
	inRow := (rl.GetScreenWidth() - 2*outerMargin) / (innerMargin + selectorSize)
	Row := n / inRow
	Collumn := n % inRow
	dynMargin := ((rl.GetScreenWidth() - 2*outerMargin) - inRow*selectorSize) / (inRow - 1)
	Position := rl.NewVector2(float32(outerMargin+Collumn*dynMargin+Collumn*selectorSize), float32(top+outerMargin+selectorSize*Row+innerMargin*Row))

	return rl.NewRectangle(Position.X, Position.Y, selectorSize, selectorSize)
}

func DoShapeSelector(S *State) *State {
	for i := 0; i < ShapeNumber; i++ {
		rec := AssignSelector(i)
		rl.DrawRectangleRec(rec, rl.NewColor(248, 248, 248, 255))
		rl.DrawRectangleLinesEx(rec, 1, rl.NewColor(187, 187, 187, 255))
		NewdRectangle(rl.NewVector2(rec.X+10, rec.Y+10), rl.NewVector2(rec.X+rec.Width-10, rec.Y+rec.Height-10)).DrawShape(i, 1, rl.Black)

		if rl.CheckCollisionPointRec(rl.GetMousePosition(), rec) {
			rl.DrawRectangleLinesEx(rec, 1, rl.NewColor(100, 100, 100, 255))
			if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
				S = CloseShapeSelector(S, i)
			}
		}
	}

	return S
}
