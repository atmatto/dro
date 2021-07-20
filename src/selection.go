package dro

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	s = rl.NewVector2(-1, -1)
	e = rl.NewVector2(-1, -1)
)

func DoSelectFragment(S *State) *State {
	if rl.IsMouseButtonUp(rl.MouseLeftButton) {
		if S.Mode&M_Shape != 0 {
			S.Alert("Draw shape or click ESC to cancel.")
		} else {
			S.Alert("Select fragment or click ESC to cancel.")
		}
	}
	if rl.IsMouseButtonPressed(rl.MouseLeftButton) && S.CanGetMouseMovement() {
		s = rl.GetMousePosition()
	}
	e = rl.GetMousePosition()
	if rl.IsKeyDown(rl.KeyLeftShift) {
		if e.X-s.X > e.Y-s.Y {
			e = rl.Vector2Add(s, rl.NewVector2(e.X-s.X, e.X-s.X))
		} else {
			e = rl.Vector2Add(s, rl.NewVector2(e.Y-s.Y, e.Y-s.Y))
		}
	}
	sel := NewdRectangle(s, e)
	if s.X != -1 && s.Y != -1 {
		if S.Mode&M_Shape != 0 {
			sel.DrawShape(S.Shape, S.Brush.Thickness*S.Canvas.Scale, S.Brush.ForegroundColor)
		} else {
			sel.DrawShape(S_RectangleLines, 1, rl.Black)
		}
		if rl.IsMouseButtonReleased(rl.MouseLeftButton) {
			return FragmentSelected(S)
		}
	}
	if rl.IsKeyPressed(rl.KeyEscape) {
		S.Mode = S.LastMode
	}
	return S
}

func FragmentSelected(S *State) *State {
	S.Alert("")
	s = S.Canvas.ScreenToCanvas(s)
	e = S.Canvas.ScreenToCanvas(e)
	S.Selection = NewdRectangle(s, e)
	s = rl.NewVector2(-1, -1)
	e = rl.NewVector2(-1, -1)
	if S.Mode&M_Shape == 0 {
		if S.Selection.X < 0 {
			S.Selection.Width += S.Selection.X
			S.Selection.X = 0
		}
		if S.Selection.Y < 0 {
			S.Selection.Height += S.Selection.Y
			S.Selection.Y = 0
		}
		if S.Selection.Width+S.Selection.X > S.Canvas.Size.X {
			S.Selection.Width = S.Canvas.Size.X - S.Selection.X
		}
		if S.Selection.Height+S.Selection.Y > S.Canvas.Size.Y {
			S.Selection.Height = S.Canvas.Size.Y - S.Selection.Y
		}
		S.SelectionTex = rl.LoadRenderTexture(int32(S.Selection.Width), int32(S.Selection.Height))
		rl.EndDrawing()
		rl.BeginTextureMode(S.SelectionTex)
		srec := S.Selection.Rec()
		fmt.Println(S.Selection, srec)
		srec = rl.NewRectangle(srec.X, S.Canvas.Size.Y-srec.Y-srec.Height, srec.Width, -srec.Height)
		rl.DrawTexturePro(S.Canvas.Tex.Texture, srec, rl.NewRectangle(0, 0, S.Selection.Width, S.Selection.Height), rl.NewVector2(0, 0), 0, rl.White)
		rl.EndTextureMode()
		rl.BeginTextureMode(*S.Canvas.Tex)
		S = S.Change(true)
		rl.DrawRectangleRec(S.Selection.Rec(), S.Brush.BackgroundColor)
		rl.EndTextureMode()
		S.Change(false)
		rl.BeginDrawing()
	}
	S.Mode &= ^M_Selection
	S.Mode |= M_Handles
	return S
}

func DoHandles(S *State) *State {
	/*if S.Mode&M_Shape != 0 {
		S.Mode = S.LastMode
		return CommitShape(S)
	}*/

	selscr := S.Selection.CanvasToScreen(&S.Canvas)

	if S.Mode&M_Shape != 0 {
		selscr.DrawShape(S.Shape, S.Brush.Thickness, S.Brush.ForegroundColor)
	} else {
		selscr.DrawRenderTex(S.SelectionTex)
	}

	rl.DrawRectangleLinesEx(selscr.Offset(-1, -1, 2, 2).Rec(), 1, rl.Black)
	rl.DrawRectangleLinesEx(selscr.Offset(-2, -2, 4, 4).Rec(), 1, rl.White)

	hovering := false
	if IsDragging {
		hovering = true
	}
	for x := 0.0; x <= 1; x += 0.5 {
		for y := 0.0; y <= 1; y += 0.5 {
			if x == 0.5 && y == 0.5 {
				continue
			}
			if DoHandle(rl.NewVector2(float32(x), float32(y)), selscr, !hovering) {
				hovering = true
				if rl.IsMouseButtonDown(rl.MouseLeftButton) {
					IsDragging = true
					Handle = rl.NewVector2(float32(x), float32(y))
				}
			}
		}
	}
	if rl.IsMouseButtonDown(rl.MouseLeftButton) && rl.CheckCollisionPointRec(rl.GetMousePosition(), selscr.Rec()) && !hovering {
		IsDragging = true
		Handle = rl.NewVector2(0.5, 0.5)
	}
	if IsDragging {
		HandleResize(S, &selscr, Handle)
	}

	if rl.IsMouseButtonPressed(rl.MouseLeftButton) && !hovering && !rl.CheckCollisionPointRec(rl.GetMousePosition(), selscr.Rec()) {
		// Apply selection
		S = S.Change(true)
		m := S.Mode
		S.Mode = S.LastMode
		if m&M_Shape != 0 {
			S = CommitShape(S)
		} else {
			S = CommitSelection(S)
			S = S.Change(false)
		}
	}

	return S
}

func DoHandle(h rl.Vector2, recscr dRectangle, interactive bool) bool {
	x := recscr.X + recscr.Width*h.X
	y := recscr.Y + recscr.Height*h.Y

	hrec := NewdRectangle(rl.NewVector2(x-float32(C_HandleSize), y-float32(C_HandleSize)), rl.NewVector2(x+float32(C_HandleSize), y+float32(C_HandleSize)))
	rl.DrawRectangleLinesEx(hrec.Offset(-1, -1, 2, 2).Rec(), 1, rl.Black)
	rl.DrawRectangleLinesEx(hrec.Offset(-2, -2, 4, 4).Rec(), 1, rl.White)

	if interactive && rl.CheckCollisionPointRec(rl.GetMousePosition(), hrec.Rec()) {
		rl.DrawRectangleLinesEx(hrec.Offset(-3, -3, 6, 6).Rec(), 1, rl.Black)
		return true
	}
	return false
}

var (
	DragStart     rl.Vector2
	DragStartDisp rl.Vector2 // Displacement of mouse cursor in relation to handle center
	StartRec      dRectangle
	IsDragging    bool = false
	Handle        rl.Vector2
)

func HandleResize(S *State, r *dRectangle, h rl.Vector2) {
	Handle = h
	hx := r.X + r.Width*h.X
	hy := r.Y + r.Height*h.Y

	if rl.IsMouseButtonUp(rl.MouseLeftButton) {
		IsDragging = false
	}
	if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		IsDragging = true
		DragStart = rl.GetMousePosition()
		DragStartDisp = rl.Vector2Subtract(DragStart, rl.NewVector2(hx, hy))
		StartRec = *r
		S.BaseSize = rl.NewVector2(r.Width, r.Height)
	}

	switch h.X {
	case 0:
		switch h.Y {
		case 0:
			r.Y = float32(rl.GetMouseY()) - DragStart.Y + StartRec.Y
			r.Height = -float32(rl.GetMouseY()) + DragStart.Y + StartRec.Height
			r.X = float32(rl.GetMouseX()) - DragStart.X + StartRec.X
			r.Width = -float32(rl.GetMouseX()) + DragStart.X + StartRec.Width
			if r.Height < 0 {
				r.Y += r.Height
				r.Height = 0
			}
			if r.Width < 0 {
				r.X += r.Width
				r.Width = 0
			}
		case 0.5:
			r.X = float32(rl.GetMouseX()) - DragStart.X + StartRec.X
			r.Width = -float32(rl.GetMouseX()) + DragStart.X + StartRec.Width
			if r.Width < 0 {
				r.X += r.Width
				r.Width = 0
			}
		case 1:
			r.Height = float32(rl.GetMouseY()) - DragStart.Y + StartRec.Height
			r.X = float32(rl.GetMouseX()) - DragStart.X + StartRec.X
			r.Width = -float32(rl.GetMouseX()) + DragStart.X + StartRec.Width
			if r.Width < 0 {
				r.X += r.Width
				r.Width = 0
			}
		}
	case 0.5:
		switch h.Y {
		case 0:
			r.Y = float32(rl.GetMouseY()) - DragStart.Y + StartRec.Y
			r.Height = -float32(rl.GetMouseY()) + DragStart.Y + StartRec.Height
			if r.Height < 0 {
				r.Y += r.Height
				r.Height = 0
			}
		case 0.5:
			r.Y = float32(rl.GetMouseY()) - DragStart.Y + StartRec.Y
			r.X = float32(rl.GetMouseX()) - DragStart.X + StartRec.X
		case 1:
			r.Height = float32(rl.GetMouseY()) - DragStart.Y + StartRec.Height
		}
	case 1:
		switch h.Y {
		case 0:
			r.Y = float32(rl.GetMouseY()) - DragStart.Y + StartRec.Y
			r.Height = -float32(rl.GetMouseY()) + DragStart.Y + StartRec.Height
			r.Width = float32(rl.GetMouseX()) - DragStart.X + StartRec.Width
			if r.Height < 0 {
				r.Y += r.Height
				r.Height = 0
			}
		case 0.5:
			r.Width = float32(rl.GetMouseX()) - DragStart.X + StartRec.Width
		case 1:
			r.Width = float32(rl.GetMouseX()) - DragStart.X + StartRec.Width
			r.Height = float32(rl.GetMouseY()) - DragStart.Y + StartRec.Height
		}
	}

	// TODO: Some handles are buggy when resizing with ratio lock
	if rl.IsKeyDown(rl.KeyLeftShift) {
		switch h.X {
		case 0:
			switch h.Y {
			case 0:
				if float32(rl.GetMouseY()) >= float32(rl.GetMouseY())+r.Height-r.Width*S.BaseSize.Y/S.BaseSize.X {
					my := float32(rl.GetMouseY())
					if my > r.Y+r.Height {
						my = r.Y + r.Height
					}
					nh := r.Width * S.BaseSize.Y / S.BaseSize.X
					dh := r.Height - nh
					r.Height = nh
					r.Y = my + dh
				} else {
					mx := float32(rl.GetMouseX())
					if mx > r.X+r.Width {
						mx = r.X + r.Width
					}
					nw := r.Height * S.BaseSize.X / S.BaseSize.Y
					dw := r.Width - nw
					r.Width = nw
					r.X = mx + dw
				}
			case 0.5:
				mx := float32(rl.GetMouseX())
				if mx > r.X+r.Width {
					mx = r.X + r.Width
				}
				nh := r.Width * S.BaseSize.Y / S.BaseSize.X
				r.X = mx
				r.Height = nh
			case 1:
				mx := float32(rl.GetMouseX())
				if mx > r.X+r.Width {
					mx = r.X + r.Width
				}
				nw := r.Height * S.BaseSize.X / S.BaseSize.Y
				dw := r.Width - nw
				r.Width = nw
				r.X = mx + dw
				r.Height = r.Width * S.BaseSize.Y / S.BaseSize.X
			}
		case 0.5:
			r.Width = r.Height * S.BaseSize.X / S.BaseSize.Y
		case 1:
			switch h.Y {
			case 0:
				if rl.GetMouseX() > int32(r.X)+int32(r.Height*S.BaseSize.X/S.BaseSize.Y) {
					my := float32(rl.GetMouseY())
					if my > r.Y+r.Height {
						my = r.Y + r.Height
					}
					nh := r.Width * S.BaseSize.Y / S.BaseSize.X
					dh := r.Height - nh
					r.Height = nh
					r.Y = my + dh
				} else {
					r.Width = r.Height * S.BaseSize.X / S.BaseSize.Y
				}
			case 0.5:
				r.Height = r.Width * S.BaseSize.Y / S.BaseSize.X
			case 1:
				if rl.GetMouseX() > int32(r.X)+int32(r.Height*S.BaseSize.X/S.BaseSize.Y) {
					r.Height = r.Width * S.BaseSize.Y / S.BaseSize.X
				} else {
					r.Width = r.Height * S.BaseSize.X / S.BaseSize.Y
				}
			}
		}
	}

	if r.Height < 1 {
		r.Height = 1
	}
	if r.Width < 1 {
		r.Width = 1
	}
	S.Selection = r.ScreenToCanvas(&S.Canvas)

}

func CommitShape(S *State) *State {
	rl.EndDrawing()
	rl.BeginTextureMode(*S.Canvas.Tex)
	S.Selection.DrawShape(S.Shape, S.Brush.Thickness, S.Brush.ForegroundColor)
	rl.EndTextureMode()
	rl.BeginDrawing()
	rl.GenTextureMipmaps(&S.Canvas.Tex.Texture)
	return S
}

func CommitSelection(S *State) *State {
	rl.EndDrawing()
	rl.BeginTextureMode(*S.Canvas.Tex)
	S.Selection.DrawRenderTex(S.SelectionTex)
	rl.EndTextureMode()
	rl.BeginDrawing()
	rl.GenTextureMipmaps(&S.Canvas.Tex.Texture)
	return S
}
