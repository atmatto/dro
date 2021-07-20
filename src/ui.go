package dro

import (
	"path/filepath"
	"strconv"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type MenuBarButton struct {
	Name     string
	Function func(S *State) *State
}

var UiBtnFlag = ""

var MenuBarButtons = []MenuBarButton{
	{Name: "New", Function: NewFile},
	{Name: "Open", Function: OpenFile},
	{Name: "Save", Function: SaveFile},
	{Name: "Quit", Function: Quit},
	{Name: "Color", Function: OpenColorDialog},
	{Name: "Shape", Function: OpenShapeSelector},
}

var usingMenu = false // Don't draw when true

func MenuPlaceholder(S *State) *State {
	//S = S.Change(false)
	return S
}

func NewFile(S *State) *State {
	UiBtnFlag = "new"
	return S
}

func OpenFile(S *State) *State {
	UiBtnFlag = "open"
	return S
}

func SaveFile(S *State) *State {
	UiBtnFlag = "save"
	return S
}

func Quit(S *State) *State {
	UiBtnFlag = "quit"
	return S
}

func DrawMenuBar(S *State) (ret *State) {
	ret = S
	var top int32 = 0
	rl.DrawLine(0, top+int32(C_StatusBarHeight)+1, int32(rl.GetScreenWidth()), top+int32(C_StatusBarHeight)+1, rl.NewColor(187, 187, 187, 255))
	rl.DrawRectangle(0, top, int32(rl.GetScreenWidth()), int32(C_StatusBarHeight), rl.NewColor(248, 248, 248, 255))

	left := 10
	for _, b := range MenuBarButtons {
		v := rl.MeasureTextEx(F, b.Name, 16, 1)
		m := rl.GetMousePosition()
		if m.Y < 20 && int(m.X) > left && m.X < float32(left)+v.X {
			if rl.IsMouseButtonUp(rl.MouseLeftButton) && rl.IsMouseButtonUp(rl.MouseMiddleButton) && rl.IsMouseButtonUp(rl.MouseRightButton) {
				rl.DrawRectangle(int32(left-5), 0, int32(v.X+10), 20, rl.NewColor(187, 187, 187, 255))
			}
			if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
				//fmt.Println(b.Name)
				ret = b.Function(S)
			}
		}
		rl.DrawTextEx(F, b.Name, rl.NewVector2(float32(left), 3), 16, 1, rl.NewColor(80, 80, 80, 255))
		left += int(v.X) + 10
	}
	if rl.GetMouseY() < 20 && rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		usingMenu = true
	}
	if rl.GetMouseY() < 20 && rl.IsMouseButtonDown(rl.MouseLeftButton) && usingMenu {
		S.MouseOld = rl.NewVector2(-1, -1)
	}
	if rl.GetMouseY() > 20 {
		usingMenu = false
	}
	return
}

func DrawStatusBar(S *State) {
	top := int32(rl.GetScreenHeight() - C_StatusBarHeight)
	rl.DrawLine(0, top, int32(rl.GetScreenWidth()), top, rl.NewColor(187, 187, 187, 255))
	rl.DrawRectangle(0, top, int32(rl.GetScreenWidth()), int32(C_StatusBarHeight), rl.NewColor(248, 248, 248, 255))

	st := ""
	if S.File.IsUnnamed {
		st = "unnamed"
	} else {
		st = filepath.Base(S.File.Path)
	}
	if !S.File.Saved {
		st += "*"
	}
	st += "   "
	st += "T " + strconv.FormatFloat(float64(S.Brush.Thickness), 'f', 0, 32) + "   Z " + strconv.FormatFloat(float64(S.Canvas.Scale*100), 'f', 0, 32) + "%   "

	/*p, n := S.DebugSnapshots()
	for i := 0; i < p; i++ {
		st += "o"
	}
	st += "|"
	for i := 0; i < n; i++ {
		st += "o"
	}*/

	if rl.GetTime()-S.AlertTime < float32(C_AlertTime) {
		st += S.AlertMessage
	}

	rl.DrawTextEx(F, st, rl.NewVector2(10, float32(rl.GetScreenHeight()-17)), 16, 1, rl.NewColor(80, 80, 80, 255))
}

func (S *State) Alert(a string) {
	S.AlertTime = rl.GetTime()
	S.AlertMessage = a
}
