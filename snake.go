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

type Cardinal byte

type Block struct {
	x    int
	y    int
	next *Block
}

func (b1 *Block) positionEquals(b2 *Block) bool {
	return b1.x == b2.x && b1.y == b2.y
}

const (
	NONE  Cardinal = 0x0
	UP    Cardinal = 0x1
	DOWN  Cardinal = 0x2
	LEFT  Cardinal = 0x4
	RIGHT Cardinal = 0x8
)

var (
	window_size, rows, snake_width int
	occupied                       []bool
)

func do_draw_block(block *Block, picture *pixel.PictureData, color color.RGBA) {
	for y := block.y * snake_width; y < (block.y+1)*snake_width; y++ {
		for x := block.x * snake_width; x < (block.x+1)*snake_width; x++ {
			picture.Pix[y*picture.Stride+x] = color
		}
	}
}

func draw_block(block *Block, picture *pixel.PictureData) {
	do_draw_block(block, picture, colornames.White)
}

func erase_block(block *Block, picture *pixel.PictureData) {
	do_draw_block(block, picture, colornames.Black)
}

func get_random_block(existing *Block) *Block {
	var new_block *Block
	do := true
	for do {
		new_block = &Block{
			x: rand.Intn(rows),
			y: rand.Intn(rows),
		}
		do = occupied[new_block.y*rows+new_block.x]
	}
	return new_block
}

func main() {
	rand.Seed(time.Now().UnixNano())
	pixelgl.Run(func() {
		window_size = 1000
		rows = 100
		snake_width = window_size / rows

		window, err := pixelgl.NewWindow(pixelgl.WindowConfig{
			Title:  "Snake",
			Bounds: pixel.R(0, 0, float64(window_size), float64(window_size)),
			// Resizable: true,
			// VSync: true,
		})
		if err != nil {
			panic(err)
		}
		window.SetPos(pixel.Vec{
			X: 1,
			Y: 27,
		})

		picture := pixel.MakePictureData(window.Bounds())
		sprite := pixel.NewSprite(picture, picture.Rect)

		occupied = make([]bool, rows*rows)

		var direction, new_direction Cardinal
		head := get_random_block(nil)
		occupied[head.y*rows+head.x] = true
		tail := head
		draw_block(head, picture)
		food := get_random_block(tail)
		do_draw_block(food, picture, colornames.Chartreuse)

		var accumulator int64 = 0
		var dt int64 = time.Second.Nanoseconds() / int64(30)
		last := time.Now()
	MAIN_LOOP:
		for !window.Closed() {
			now := time.Now()
			frame_time := now.Sub(last).Nanoseconds()
			accumulator += frame_time
			last = now

			switch {
			case window.JustPressed(pixelgl.KeyEscape):
				window.Destroy()
				break MAIN_LOOP
			case direction != DOWN && window.JustPressed(pixelgl.KeyUp):
				new_direction = UP
			case direction != UP && window.JustPressed(pixelgl.KeyDown):
				new_direction = DOWN
			case direction != RIGHT && window.JustPressed(pixelgl.KeyLeft):
				new_direction = LEFT
			case direction != LEFT && window.JustPressed(pixelgl.KeyRight):
				new_direction = RIGHT
			}

			if accumulator >= dt {
				direction = new_direction
				if direction != NONE {
					var new_x, new_y int
					switch direction {
					case UP:
						new_x = head.x
						new_y = head.y + 1
					case DOWN:
						new_x = head.x
						new_y = head.y - 1
					case LEFT:
						new_x = head.x - 1
						new_y = head.y
					case RIGHT:
						new_x = head.x + 1
						new_y = head.y
					}

					if new_x < 0 || new_x >= rows || new_y < 0 || new_y >= rows {
						fmt.Printf("%d, %d out of bounds\n", new_x, new_y)
						break MAIN_LOOP
					}

					head.next = &Block{
						x: new_x,
						y: new_y,
					}
					head = head.next
					draw_block(head, picture)

					if head.positionEquals(food) {
						fmt.Printf("Yum!\n")
						food = get_random_block(tail)
						do_draw_block(food, picture, colornames.Chartreuse)

					} else {
						occupied[tail.y*rows+tail.x] = false
						erase_block(tail, picture)
						tail_next := tail.next
						tail.next = nil
						tail = tail_next
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
	})
}
