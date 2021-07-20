package main

import (
	"fmt"
	"os"
	"path/filepath"

	dro "github.com/atmatto/dro/src"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func usage() {
	fmt.Println("usage: dro [FILE]")
	os.Exit(1)
}

func main() {
	cmd := os.Args

	file := dro.File{Path: "", IsUnnamed: true, Saved: true}

	if len(cmd) > 2 {
		usage()
	}

	rl.SetTraceLog(rl.LogWarning)

	rl.InitWindow(800, 450, "dro")
	rl.SetWindowState(rl.FlagWindowResizable)
	rl.SetWindowMinSize(500, 300)
	rl.SetTargetFPS(60)
	rl.SetExitKey(0)

	dro.LoadFont()

	var S *dro.State = new(dro.State)
	S.Reset()
	S.File = &file
	makenew := false // If true, the next frame a new State and file will be created

	if len(cmd) == 2 {
		file.Path = cmd[1]
	}
	var err error
	if file.Path != "" {
		_, err = os.Stat(file.Path)
	}
	if file.Path == "" || os.IsNotExist(err) {
		S.Canvas.Load(nil)
	} else {
		t := rl.LoadTexture(file.Path)
		S.Canvas.Load(&t)
		file.IsUnnamed = false
	}

	if dro.C_Filter {
		rl.GenTextureMipmaps(&S.Canvas.Tex.Texture)
		rl.SetTextureFilter(S.Canvas.Tex.Texture, rl.FilterTrilinear)
	}
	S = S.Change(true)
	file.Saved = true

	usesBars := false

	for {

		if rl.IsMouseButtonDown(rl.MouseLeftButton) && rl.CheckCollisionPointRec(rl.GetMousePosition(), rl.NewRectangle(0, 0, float32(rl.GetScreenWidth()), 20)) && rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			usesBars = true
		}
		if rl.IsMouseButtonDown(rl.MouseLeftButton) && rl.CheckCollisionPointRec(rl.GetMousePosition(), rl.NewRectangle(0, float32(rl.GetScreenHeight())-20, float32(rl.GetScreenWidth()), 20)) && rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			usesBars = true
		}
		if usesBars {
			S.MouseOld = rl.NewVector2(-1, -1)
		}

		rl.BeginTextureMode(*S.Canvas.Tex)
		if S.Mode&dro.M_Brush != 0 {
			S = S.Brush.DoDrawing(S)
			if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
				S = S.Change(true)
			}
		}
		rl.EndTextureMode()

		rl.BeginDrawing()
		rl.ClearBackground(rl.Gray)

		S.Canvas.DoScale()
		S.Canvas.Draw()
		S.Canvas.Pan(S)
		S.Brush.DoThickness()
		//dro.NewdRectangle(S.MouseOld, rl.GetMousePosition()).TestDirection()

		if S.Brush.CanDraw(S) {
			S.MouseOld = rl.GetMousePosition()
		}

		if S.Mode&dro.M_Brush != 0 && rl.GetMouseY() > 20 && rl.GetMouseY()+20 < int32(rl.GetScreenHeight()) {
			rl.HideCursor()
			rl.DrawCircleLines(rl.GetMouseX(), rl.GetMouseY(), S.Canvas.Scale*S.Brush.Thickness/2, S.Brush.ForegroundColor)
		} else {
			rl.ShowCursor()
		}

		if S.Mode == dro.M_Dialog {
			switch S.Dialog {
			case dro.D_Color:
				S = dro.DoColorDialog(S)
			case dro.D_Shapes:
				S = dro.DoShapeSelector(S)
			case dro.D_Prompt:
				S.DoPrompt()
			}
		}

		if S.Mode&dro.M_Selection != 0 {
			S = dro.DoSelectFragment(S)
		}
		if S.Mode&dro.M_Handles != 0 {
			S = dro.DoHandles(S)
		}

		S = dro.DrawMenuBar(S)
		dro.DrawStatusBar(S)

		rl.EndDrawing()

		if !(S.Mode == dro.M_Dialog && S.Dialog == dro.D_Prompt) {
			if rl.IsKeyPressed(rl.KeyZ) && rl.IsKeyDown(rl.KeyLeftControl) {
				S = S.Undo()
			}
			if rl.IsKeyPressed(rl.KeyY) && rl.IsKeyDown(rl.KeyLeftControl) {
				S = S.Redo()
			}
			if rl.IsKeyPressed(rl.KeyC) && !rl.IsKeyDown(rl.KeyLeftControl) {
				S = dro.OpenColorDialog(S)
			}
			if rl.IsKeyPressed(rl.KeyS) && !rl.IsKeyDown(rl.KeyLeftControl) {
				S.Mode = dro.M_Selection
			}
			if dro.UiBtnFlag == "save" || (rl.IsKeyPressed(rl.KeyS) && rl.IsKeyDown(rl.KeyLeftControl)) {
				if file.IsUnnamed {
					h, _ := os.UserHomeDir()
					S.ShowPrompt("save", "Save file", h+"/")
				} else {
					S.ShowPrompt("save", "Save file", file.Path)
				}
			}
			if dro.UiBtnFlag == "new" || (rl.IsKeyPressed(rl.KeyN) && rl.IsKeyDown(rl.KeyLeftControl)) {
				if file.Saved {
					makenew = true
				} else {
					S.Alert("Warning: current file is not saved!")
					S.ShowPrompt("new/open", "Create new file. Leave the field empty or enter a path:", "")
				}
			}
			if dro.UiBtnFlag == "open" || (rl.IsKeyPressed(rl.KeyO) && rl.IsKeyDown(rl.KeyLeftControl)) {
				if !file.Saved {
					S.Alert("Warning: current file is not saved!")
				}
				S.ShowPrompt("new/open", "Open file. Enter a path or leave the field empty to create a new file:", "")
			}
			if dro.UiBtnFlag == "quit" || (rl.IsKeyPressed(rl.KeyQ) && rl.IsKeyDown(rl.KeyLeftControl)) {
				if file.Saved || rl.IsKeyDown(rl.KeyLeftShift) {
					break
				} else {
					S.Alert("Save the file before exiting, or hold shift to force exit.")
				}
			}
		}

		dro.UiBtnFlag = ""

		if rl.WindowShouldClose() {
			if file.Saved {
				break
			} else {
				S.Alert("Save the file before exiting, or use ctrl-shift-q to force exit.")
			}
		}

		if S.Prompt.Accepted {
			S.Prompt.Accepted = false
			switch S.Prompt.Action {
			case "save":
				if filepath.Ext(S.Prompt.Input) == "" {
					S.Prompt.Input += ".png"
				}
				if filepath.Ext(S.Prompt.Input) != ".png" {
					S.Alert("Dro can only save .png files.")
				} else {
					i := rl.GetTextureData(S.Canvas.Tex.Texture)
					rl.ImageFlipVertical(i)
					rl.ExportImage(*i, S.Prompt.Input)
					rl.UnloadImage(i)
					S.Alert("File saved: " + S.Prompt.Input)
					file.Saved = true
					file.Path = S.Prompt.Input
					file.IsUnnamed = false
				}
			case "new/open":
				makenew = true
				file.Path = S.Prompt.Input
				if file.Path == "" {
					file.IsUnnamed = true
				} else {
					file.IsUnnamed = false
				}
				file.Saved = false
			}
		}

		if makenew {
			makenew = false
			S.Cleanup(true)
			S = new(dro.State)
			S.Reset()
			S.File = &file

			if file.Path == "" || os.IsNotExist(err) {
				S.Canvas.Load(nil)
			} else {
				t := rl.LoadTexture(file.Path)
				S.Canvas.Load(&t)
				file.IsUnnamed = false
			}

			if dro.C_Filter {
				rl.SetTextureFilter(S.Canvas.Tex.Texture, rl.FilterTrilinear)
			}
			S = S.Change(true)
			file.Saved = true
		}

		if rl.IsMouseButtonUp(rl.MouseLeftButton) {
			usesBars = false
		}
	}

	rl.CloseWindow()
}
