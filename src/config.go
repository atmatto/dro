package dro

import rl "github.com/gen2brain/raylib-go/raylib"

var (
	C_Timeout = 0.1 // Don't draw when rl.GetTime() < State.TimeoutT + C_Timeout

	// Show alert for ... seconds
	C_AlertTime         = 4
	C_DefaultThickness  = 6
	C_DefaultCanvasSize = rl.NewVector2(1920, 1080)
	C_StatusBarHeight   = 20

	C_ChangeSnapshotT = 1 // How often to take state snapshots for history [s]
	C_MaxSnapshots    = 16

	C_Filter = true // If true, texture filtering with mipmaps will be used.

	C_HandleSize = 8 // Half of handle size (~radius)

	C_PromptMargin   = 30
	C_PromptMinWidth = 200
	C_PromptMaxWidth = 800
)
