package main

import (
	"fmt"
	"image/color"
	"math/rand"
	"time"

	"golang.org/x/image/colornames"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type cardinal byte

type block struct {
	x    int
	y    int
	next *block
}

const (
	none  cardinal = 0x0
	up    cardinal = 0x1
	down  cardinal = 0x2
	left  cardinal = 0x4
	right cardinal = 0x8
)

var (
	picture                      *pixel.PictureData
	windowSize, rows, snakeWidth int
	occupied                     []bool
)

func drawBlockWithColor(b *block, c color.RGBA) {
	for y := b.y * snakeWidth; y < (b.y+1)*snakeWidth; y++ {
		for x := b.x * snakeWidth; x < (b.x+1)*snakeWidth; x++ {
			picture.Pix[y*picture.Stride+x] = c
		}
	}
}

func drawBlock(b *block) {
	drawBlockWithColor(b, colornames.White)
}

func eraseBlock(b *block) {
	drawBlockWithColor(b, colornames.Black)
}

func getRandomBlock() *block {
	var b *block
	do := true
	for do {
		b = &block{
			x: rand.Intn(rows),
			y: rand.Intn(rows),
		}
		do = occupied[b.y*rows+b.x]
	}
	return b
}
func run() {
	windowSize = 1000
	rows = 100
	snakeWidth = windowSize / rows

	window, err := pixelgl.NewWindow(pixelgl.WindowConfig{
		Title:  "Snake",
		Bounds: pixel.R(0, 0, float64(windowSize), float64(windowSize)),
		// Resizable: true,
		VSync: true,
	})
	if err != nil {
		panic(err)
	}
	window.SetPos(pixel.Vec{
		X: 1,
		Y: 27,
	})

	picture = pixel.MakePictureData(window.Bounds())
	sprite := pixel.NewSprite(picture, picture.Rect)

	occupied = make([]bool, rows*rows)

	var direction, newDirection cardinal
	head := getRandomBlock()
	occupied[head.y*rows+head.x] = true
	tail := head
	drawBlock(head)
	food := getRandomBlock()
	drawBlockWithColor(food, colornames.Chartreuse)

	var accumulator int64 = 0
	var dt int64 = time.Second.Nanoseconds() / int64(30)
	last := time.Now()
MAIN_LOOP:
	for !window.Closed() {
		now := time.Now()
		frameTime := now.Sub(last).Nanoseconds()
		accumulator += frameTime
		last = now

		switch {
		case window.JustPressed(pixelgl.KeyEscape):
			window.Destroy()
			break MAIN_LOOP
		case direction != down && window.JustPressed(pixelgl.KeyUp):
			newDirection = up
		case direction != up && window.JustPressed(pixelgl.KeyDown):
			newDirection = down
		case direction != right && window.JustPressed(pixelgl.KeyLeft):
			newDirection = left
		case direction != left && window.JustPressed(pixelgl.KeyRight):
			newDirection = right
		}

		if accumulator >= dt {
			direction = newDirection
			if direction != none {
				var newX, newY int
				switch direction {
				case up:
					newX = head.x
					newY = head.y + 1
				case down:
					newX = head.x
					newY = head.y - 1
				case left:
					newX = head.x - 1
					newY = head.y
				case right:
					newX = head.x + 1
					newY = head.y
				}

				if newX < 0 || newX >= rows || newY < 0 || newY >= rows {
					fmt.Printf("%d, %d out of bounds\n", newX, newY)
					break MAIN_LOOP
				}

				head.next = &block{
					x: newX,
					y: newY,
				}
				head = head.next
				drawBlock(head)

				if head.x == food.x && head.y == food.y {
					fmt.Printf("Yum!\n")
					food = getRandomBlock()
					drawBlockWithColor(food, colornames.Chartreuse)

				} else {
					occupied[tail.y*rows+tail.x] = false
					eraseBlock(tail)
					tailNext := tail.next
					tail.next = nil
					tail = tailNext
				}

				if occupied[head.y*rows+head.x] {
					fmt.Printf("%d, %d you ate yourself\n", head.x, head.y)
					break MAIN_LOOP
				}
				occupied[head.y*rows+head.x] = true
			}

			// for block := tail; block != nil; block = block.next {
			// 	fmt.Printf("(%d, %d) -> ", block.x, block.y)
			// }
			// fmt.Printf("nil\n\n")

			sprite = pixel.NewSprite(picture, picture.Rect)
			accumulator -= dt
		}

		sprite.Draw(window, pixel.IM.Moved(window.Bounds().Center()))
		window.Update()
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	pixelgl.Run(run)
}
