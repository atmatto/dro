package dro

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// ModeMask
const (
	M_Brush     = 1 << iota // Brush enabled
	M_Dialog                // Display a dialog
	M_Selection             // Select a region on canvas
	M_Handles               // Manipulate a region on canvas

	// Selection is (or is going to be, in case of M_Selection) stored separately from main canvas.
	// Used for things like clipboard and moving fragments of canvas. Not set e.g. when the region is used to draw a shape.
	M_ExternalSelection
	M_Shape
)

// Dialog
const (
	D_Shapes = iota
	D_Color
	D_Prompt
)

type State struct {
	// Brush
	Brush Brush

	// Editor mode

	Mode     int
	LastMode int

	// M_Brush

	// M_Dialog

	Dialog int

	// D_Prompt
	Prompt *Prompt

	// M_Selection | M_Handles

	BaseSize  rl.Vector2 // Used for aspect ratio calculations
	Selection dRectangle // Canvas space

	// M_ExternalSelection
	SelectionTex rl.RenderTexture2D

	// M_Shape
	Shape int

	// Canvas
	Canvas Canvas

	// Input

	MouseOld rl.Vector2 // Don't draw when MouseOld==(Vector2){-1,-1}
	TimeoutT float32    // Don't draw when GetTime() < TimeoutT + C_Timeout

	// History

	Time      float32  // When was this state created
	OldCanvas rl.Image // Only for history
	Next      *State
	Previous  *State

	// File
	File *File

	// Alert

	AlertTime    float32
	AlertMessage string
}

func (S *State) Reset() {
	// Commented values are not set. They should be set manually before use.
	S.Brush.ForegroundColor = rl.Black
	S.Brush.BackgroundColor = rl.White
	S.Brush.Thickness = float32(C_DefaultThickness)
	S.Mode = M_Brush
	S.LastMode = M_Brush
	//S.Dialog
	S.Prompt = &Prompt{}
	//S.BaseSize
	//S.Selection1
	//S.Selection2
	//S.Selection
	S.Shape = -1
	//S.Canvas
	S.Canvas.Position = rl.Vector2{X: 0, Y: 20}
	S.Canvas.Size = C_DefaultCanvasSize
	S.Canvas.Scale = 1
	S.MouseOld = rl.Vector2{X: -1, Y: -1}
	S.TimeoutT = rl.GetTime()
	S.Time = rl.GetTime()
	S.Next = nil
	S.Previous = nil
	// S.File
	S.AlertTime = -100
	S.AlertMessage = ""
}

func (S *State) CanGetMouseMovement() bool {
	if S.MouseOld.X == -1 && S.MouseOld.Y == -1 {
		S.MouseOld = rl.GetMousePosition()
		return false
	}
	return rl.GetTime() > S.TimeoutT+float32(C_Timeout)
}

func (S *State) Change(snapshot bool) *State {
	if !snapshot { //&& rl.GetTime() < S.Time+float32(C_ChangeSnapshotT) {
		if C_Filter {
			rl.GenTextureMipmaps(&S.Canvas.Tex.Texture)
		}
		return S
	}

	S.AlertTime = -100

	N := *S

	N.Next = S.Next
	if N.Next != nil {
		N.Next.Previous = &N
	}

	S.Next = &N
	N.Previous = S

	N.Time = rl.GetTime()

	S.OldCanvas = *rl.GetTextureData(S.Canvas.Tex.Texture)
	N.Canvas = S.Canvas

	N.SelectionTex = S.SelectionTex

	N.Prompt = S.Prompt

	N.File = S.File
	N.File.Saved = false

	N.Cleanup(false)
	if C_Filter {
		rl.GenTextureMipmaps(&N.Canvas.Tex.Texture)
	}
	return &N
}

func (S *State) DebugSnapshots() (int, int) {
	prev := 0
	next := 0

	ptr := S
	for {
		if ptr.Previous != nil {
			ptr = ptr.Previous
			prev++
		} else {
			break
		}
	}
	ptr = S
	for {
		if ptr.Next != nil {
			ptr = ptr.Next
			next++
		} else {
			break
		}
	}

	return prev, next
}

func (S *State) DebugMemory() {
	p, n := S.DebugSnapshots()

	fmt.Println("----------")

	ptr := S
	for i := 1; i <= p; i++ {
		ptr = ptr.Previous
		fmt.Print("<")
	}
	fmt.Print("|")
	for i := 1; i <= n; i++ {
		fmt.Print(">")
	}
	fmt.Println("\n----------")
	for i := 0; i <= p+n; i++ {
		if ptr == S {
			fmt.Printf("%p S\n", ptr)
		} else {
			fmt.Printf("%p\n", ptr)
		}
		ptr = ptr.Next
	}

	fmt.Println("----------")
}

func (S *State) Undo() *State {
	if S.Previous != nil {
		if S.OldCanvas.Height != 0 {
			rl.UnloadImage(&S.OldCanvas)
		}
		S.OldCanvas = *rl.GetTextureData(S.Canvas.Tex.Texture)
		S.Previous.Canvas = S.Canvas
		//S.Previous.Canvas.Tex.Texture = rl.LoadTextureFromImage(&S.Previous.OldCanvas)
		rl.UpdateTexture(S.Previous.Canvas.Tex.Texture, rl.GetImageData(&S.Previous.OldCanvas))
		//S.Previous.DebugMemory()
		S.Previous.Time = rl.GetTime()
		if C_Filter {
			rl.GenTextureMipmaps(&S.Previous.Canvas.Tex.Texture)
		}
		return S.Previous
	} else {
		S.Alert("Nothing to undo.")
		return S
	}
}

func (S *State) Redo() *State {
	if S.Next != nil {
		S.Next.Canvas = S.Canvas
		//S.Next.Canvas.Tex.Texture = rl.LoadTextureFromImage(&S.Next.OldCanvas)
		rl.UpdateTexture(S.Next.Canvas.Tex.Texture, rl.GetImageData(&S.Next.OldCanvas))
		//S.Next.DebugMemory()
		S.Next.Time = rl.GetTime()
		if C_Filter {
			rl.GenTextureMipmaps(&S.Next.Canvas.Tex.Texture)
		}
		return S.Next
	} else {
		S.Alert("Nothing to redo.")
		return S
	}
}

// Remove unneded States from memory
func (S *State) Cleanup(all bool) {
	// Forward
	//S.DebugMemory()
	ptr := S
	for {
		if ptr.Next != nil {
			//fmt.Printf("ptr.Next != nil, %p, %p \n", ptr, ptr.Next)
			ptr = ptr.Next
		} else {
			//fmt.Printf("ptr.Next == nil, %p, %p \n", ptr, ptr.Next)
			break
		}
	}
	for {
		if ptr != S {
			//fmt.Printf("ptr != S,        %p %p\n", ptr, S)
			ptr = ptr.Previous
			ptr.Next.Unload()
		} else {
			//fmt.Printf("ptr == S,        %p %p\n", ptr, S)
			break
		}
	}
	// Backward
	ptr = S
	prevCount := 0
	for {
		if ptr.Previous != nil {
			prevCount++
			ptr = ptr.Previous
		} else {
			break
		}
	}
	for {
		if all || prevCount > C_MaxSnapshots {
			prevCount--
			if ptr.Next == nil {
				if all {
					ptr.Unload()
				}
				break
			}
			ptr = ptr.Next
			ptr.Previous.Unload()
		} else {
			break
		}
	}
}

// Unload won't check if S.Previous or S.Next should be unloaded.
func (S *State) Unload() {
	if S.Previous != nil {
		S.Previous.Next = nil
		S.Previous = nil
	}
	if S.Next != nil {
		S.Next.Previous = nil
		S.Next = nil
	}
	if S.OldCanvas.Height != 0 {
		rl.UnloadImage(&S.OldCanvas)
	}
	if S.SelectionTex.Texture.Width != 0 {
		rl.UnloadRenderTexture(S.SelectionTex)
	}
}

type File struct {
	Path      string
	Saved     bool
	IsUnnamed bool
}
