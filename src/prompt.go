package dro

import (
	"math"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Prompt struct {
	Action      string // short action identifier
	Prompt      string // prompt displayed to user
	Input       string
	InputCursor int // cursor position
	Accepted    bool
}

func (S *State) ShowPrompt(action string, prompt string, input string) {
	if S.Mode != M_Dialog {
		S.LastMode = S.Mode
	}
	S.Mode = M_Dialog
	S.Dialog = D_Prompt
	S.Prompt.Action = action
	S.Prompt.Prompt = prompt
	S.Prompt.Input = input
	S.Prompt.InputCursor = len(input)
	S.Prompt.Accepted = false
}

func (S *State) DoPrompt() {
	x := C_PromptMargin
	y := 30
	w := rl.GetScreenWidth() - 2*C_PromptMargin
	h := 62

	if w > C_PromptMaxWidth {
		dw := w - C_PromptMaxWidth
		w -= dw
		x += dw / 2
	}

	rl.DrawRectangle(int32(x), int32(y), int32(w), int32(h), rl.NewColor(248, 248, 248, 255))
	rl.DrawRectangleLines(int32(x), int32(y), int32(w), int32(h), rl.NewColor(187, 187, 187, 255))
	rl.DrawRectangle(int32(x)+10, int32(y)+28, int32(w)-20, 26, rl.White)
	rl.DrawRectangleLines(int32(x)+10, int32(y)+28, int32(w)-20, 26, rl.NewColor(187, 187, 187, 255))

	rl.DrawTextEx(F, S.Prompt.Prompt, rl.NewVector2(float32(x)+10, float32(y+10)), 16, 1, rl.NewColor(80, 80, 80, 255))

	if rl.IsKeyPressed(rl.KeyEscape) ||
		(rl.IsMouseButtonPressed(rl.MouseLeftButton) &&
			!rl.CheckCollisionPointRec(rl.GetMousePosition(), rl.NewRectangle(float32(x), float32(y), float32(w), float32(h)))) {
		S.Mode = S.LastMode
	}

	DoPromptText(S, rl.NewVector2(float32(x)+20, float32(y+34)))
}

var lastCursorChange = rl.GetTime()

func DoPromptText(S *State, position rl.Vector2) {
	rl.DrawTextEx(F, S.Prompt.Input, position, 16, 1, rl.NewColor(80, 80, 80, 255))

	cx := rl.MeasureTextEx(F, S.Prompt.Input[:S.Prompt.InputCursor], 16, 1)
	if math.Mod(float64(rl.GetTime()-lastCursorChange), 1) <= 0.5 {
		rl.DrawTextEx(F, "|", rl.NewVector2(position.X+cx.X, position.Y), 16, 1, rl.Black)
	}

	// TODO: Change when rl.GetCharPressed() will be available
	k := rl.GetKeyPressed()
	c := ""
	if k >= 32 && k < 127 {
		c = string(rune(k))
		if !rl.IsKeyDown(rl.KeyLeftShift) && !rl.IsKeyDown(rl.KeyRightShift) {
			c = strings.ToLower(c)
		}
	}

	if rl.IsKeyDown(rl.KeyLeftControl) {
		c = ""
	}

	if c != "" {
		if S.Prompt.InputCursor == len(S.Prompt.Input) {
			S.Prompt.Input += c
			S.Prompt.InputCursor += 1
		} else if S.Prompt.InputCursor == 0 {
			S.Prompt.Input = c + S.Prompt.Input
			S.Prompt.InputCursor += 1
		} else {
			S.Prompt.Input = S.Prompt.Input[0:S.Prompt.InputCursor] + c + S.Prompt.Input[S.Prompt.InputCursor:]
			S.Prompt.InputCursor += 1
		}
	}

	if rl.IsKeyPressed(rl.KeyBackspace) && S.Prompt.InputCursor != 0 {
		S.Prompt.InputCursor -= 1
		S.Prompt.Input = S.Prompt.Input[:S.Prompt.InputCursor] + S.Prompt.Input[S.Prompt.InputCursor+1:]
	}

	if rl.IsKeyPressed(rl.KeyLeft) {
		S.Prompt.InputCursor -= 1
		lastCursorChange = rl.GetTime()
	}
	if rl.IsKeyPressed(rl.KeyRight) {
		S.Prompt.InputCursor += 1
		lastCursorChange = rl.GetTime()
	}
	if S.Prompt.InputCursor < 0 || rl.IsKeyPressed(rl.KeyUp) {
		S.Prompt.InputCursor = 0
		lastCursorChange = rl.GetTime()
	}
	if S.Prompt.InputCursor > len(S.Prompt.Input) || rl.IsKeyPressed(rl.KeyDown) {
		S.Prompt.InputCursor = len(S.Prompt.Input)
		lastCursorChange = rl.GetTime()
	}

	if rl.IsKeyPressed(rl.KeyEnter) {
		S.Mode = S.LastMode
		S.Prompt.Accepted = true
	}
}
