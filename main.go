package main

import rl "github.com/gen2brain/raylib-go/raylib"

func main() {
	rend := New_Render()
	cpu := New_CPU(10, rend, "program.ch8")

	rl.InitWindow(SCREEN_SIZE_X*PIXEL_SCALE, SCREEN_SIZE_Y*PIXEL_SCALE, "emutest")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		cpu.Cycle()

		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		rend.Draw()
		rl.EndDrawing()
	}
}
