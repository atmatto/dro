package dro

import rl "github.com/gen2brain/raylib-go/raylib"

func Vector2Divide(v rl.Vector2, i float32) rl.Vector2 {
	return rl.NewVector2(v.X/i, v.Y/i)
}
