package main

import rl "github.com/gen2brain/raylib-go/raylib"

const PIXEL_SCALE = 16
const SCREEN_SIZE_X = 64
const SCREEN_SIZE_Y = 32

type Render struct {
	affichage [64][32]bool
}

func New_Render() *Render {
	r := new(Render)
	return r
}

func (r *Render) Set_pixel(x, y uint8) bool {
	r.affichage[x][y] = !r.affichage[x][y]
	return !r.affichage[x][y]
}

func (r *Render) Clear_render() {
	r.affichage = [64][32]bool{}
}

func (r *Render) Draw() {
	for x := 0; x < SCREEN_SIZE_X; x++ {
		for y := 0; y < SCREEN_SIZE_Y; y++ {
			if r.affichage[x][y] {
				rl.DrawRectangle(int32(x*PIXEL_SCALE), int32(y*PIXEL_SCALE), PIXEL_SCALE, PIXEL_SCALE, rl.White)
			}
		}
	}
}
