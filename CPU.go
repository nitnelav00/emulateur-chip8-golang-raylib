package main

import (
	"bufio"
	"io"
	"math/rand"
	"os"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type CPU struct {
	memory      [4096]uint8
	pc          uint16
	i           uint16
	stack       []uint16
	delay_timer uint8
	sound_timer uint8
	registres   [16]uint8
	speed       int
	pause       bool
	render      *Render
}

func New_CPU(speed int, rend *Render, rom string) *CPU {
	a := new(CPU)
	a.pc = 0x200
	a.load_fonts()
	a.laod_rom(rom)
	a.speed = speed
	a.render = rend
	return a
}

func (c *CPU) load_fonts() {
	fonts := []uint8{
		0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
		0x20, 0x60, 0x20, 0x20, 0x70, // 1
		0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
		0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
		0x90, 0x90, 0xF0, 0x10, 0x10, // 4
		0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
		0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
		0xF0, 0x10, 0x20, 0x40, 0x40, // 7
		0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
		0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
		0xF0, 0x90, 0xF0, 0x90, 0x90, // A
		0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
		0xF0, 0x80, 0x80, 0x80, 0xF0, // C
		0xE0, 0x90, 0x90, 0x90, 0xE0, // D
		0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
		0xF0, 0x80, 0xF0, 0x80, 0x80, // F
	}

	for i := 0; i < len(fonts); i++ {
		c.memory[i] = fonts[i]
	}
}

func (c *CPU) laod_rom(rom string) {
	f, err := os.Open(rom)
	defer f.Close()
	if err != nil {
		panic(err)
	}

	stat, err := f.Stat()
	if err != nil {
		panic(err)
	}

	size := stat.Size()
	tmp := make([]uint8, size)
	_, err = bufio.NewReader(f).Read(tmp)
	if err != nil && err != io.EOF {
		panic(err)
	}

	c.load_program(tmp, int(size))
}

func (c *CPU) load_program(prog []uint8, size int) {
	for i := 0; i < size; i++ {
		c.memory[i+0x200] = prog[i]
	}
}

func (c *CPU) Cycle() {
	for i := 0; i < c.speed; i++ {
		if c.pause {
			break
		}
		var opcode uint16 = uint16(c.memory[c.pc])<<8 + uint16(c.memory[c.pc+1])
		c.execute(opcode)
	}

	if !c.pause {
		c.update_timers()
	}
}

func (c *CPU) update_timers() {
	if c.delay_timer > 0 {
		c.delay_timer--
	}
	if c.sound_timer > 0 {
		c.sound_timer--
	}
}

func (c *CPU) execute(opcode uint16) {
	c.pc += 2
	x := (opcode & 0x0F00) >> 8
	y := (opcode & 0x00F0) >> 4
	var kk uint8 = uint8(opcode & 0x00FF)
	n := opcode & 0x000F
	nnn := opcode & 0x0FFF

	switch opcode & 0xF000 {
	case 0x0000:
		switch opcode {
		case 0x00E0:
			c.render.Clear_render()
		case 0x0EE:
			c.pc = c.stack[len(c.stack)-1]
			c.stack = c.stack[:len(c.stack)-1]
		}
	case 0x1000:
		c.pc = nnn
	case 0x2000:
		c.stack = append(c.stack, c.pc)
		c.pc = nnn
	case 0x3000:
		if c.registres[x] == kk {
			c.pc += 2
		}
	case 0x4000:
		if c.registres[x] != kk {
			c.pc += 2
		}
	case 0x5000:
		if c.registres[x] == c.registres[y] {
			c.pc += 2
		}
	case 0x6000:
		c.registres[x] = kk
	case 0x7000:
		c.registres[x] += kk
	case 0x8000:
		switch opcode & 0xF {
		case 0x0:
			c.registres[x] = c.registres[y]
		case 0x1:
			c.registres[x] |= c.registres[y]
		case 0x2:
			c.registres[x] &= c.registres[y]
		case 0x3:
			c.registres[x] ^= c.registres[y]
		case 0x4:
			var sum uint16 = uint16(c.registres[x] + c.registres[y])
			c.registres[x] = uint8(sum & 0xFF)
			c.registres[0xF] = 0
			if sum > 0xFF {
				c.registres[0xF] = 1
			}
		case 0x5:
			var sub int16 = int16(c.registres[x] - c.registres[y])
			c.registres[x] = uint8(sub & 0xFF)
			c.registres[0xF] = 0
			if sub < 0 {
				c.registres[0xF] = 1
			}
		case 0x6:
			c.registres[0xF] = c.registres[x] & 1
			c.registres[x] >>= 1
		case 0x7:
			var sub int16 = int16(c.registres[y] - c.registres[x])
			c.registres[y] = uint8(sub & 0xFF)
			c.registres[0xF] = 0
			if sub < 0 {
				c.registres[0xF] = 1
			}
		case 0xE:
			c.registres[0xF] = c.registres[x] >> 7
			c.registres[x] <<= 1
		}
	case 0x9000:
		if c.registres[x] != c.registres[y] {
			c.pc += 2
		}
	case 0xA000:
		c.i = nnn
	case 0xB000:
		c.pc = uint16(c.registres[0]) + nnn
	case 0xC000:
		var rd uint8 = uint8(rand.Int())
		c.registres[x] = rd & kk
	case 0xD000:
		width := 8
		height := n
		c.registres[0xF] = 0
		for row := 0; row < int(height); row++ {
			sprite := c.memory[c.i+uint16(row)]
			for col := 0; col < width; col++ {
				if (sprite & 0x80) > 0 {
					xpos := c.registres[x] + uint8(col)
					ypos := c.registres[y] + uint8(row)
					if xpos >= 0 && xpos < SCREEN_SIZE_X && ypos >= 0 && ypos < SCREEN_SIZE_Y {
						pixel_res := c.render.Set_pixel(xpos, ypos)
						if pixel_res {
							c.registres[0xF] = 1
						}
					}
				}
				sprite <<= 1
			}
		}
	case 0xE000:
		switch kk {
		case 0x9E:
			key := translatekey(c.registres[x])
			if rl.IsKeyDown(key) {
				c.pc += 2
			}
		case 0xA1:
			key := translatekey(c.registres[x])
			if !rl.IsKeyDown(key) {
				c.pc += 2
			}
		}
	case 0xF000:
		switch kk {
		case 0x07:
			c.registres[x] = c.delay_timer
		case 0x0A:
		case 0x15:
			c.delay_timer = c.registres[x]
		case 0x18:
			c.sound_timer = c.registres[x]
		case 0x1E:
			c.i += uint16(c.registres[x])
		case 0x29:
			c.i = uint16(c.registres[x]) * 5
		case 0x33:
			c.memory[c.i] = c.registres[x] / 100
			c.memory[c.i+1] = (c.registres[x] % 100) / 10
			c.memory[c.i+2] = c.registres[x] % 10
		case 0x55:
			for index := 0; index <= int(x); index++ {
				c.memory[int(c.i)+index] = c.registres[index]
			}
		case 0x65:
			for index := 0; index <= int(x); index++ {
				c.registres[index] = c.memory[int(c.i)+index]
			}
		}
	}
}

func translatekey(key uint8) int32 {
	switch key {
	case 0x1:
		return rl.KeyOne
	case 0x2:
		return rl.KeyTwo
	case 0x3:
		return rl.KeyThree
	case 0xC:
		return rl.KeyFour
	case 0x4:
		return rl.KeyQ
	case 0x5:
		return rl.KeyW
	case 0x6:
		return rl.KeyE
	case 0xD:
		return rl.KeyR
	case 0x7:
		return rl.KeyA
	case 0x8:
		return rl.KeyS
	case 0x9:
		return rl.KeyD
	case 0xE:
		return rl.KeyF
	case 0xA:
		return rl.KeyZ
	case 0x0:
		return rl.KeyX
	case 0xB:
		return rl.KeyC
	case 0xF:
		return rl.KeyV
	}
	return 0
}
