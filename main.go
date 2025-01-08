package main

import (
	"image/color"

	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
	"github.com/gopxl/pixel/v2/ext/imdraw"
)

const TILE_COLUMNS, TILE_ROWS int = 16 * 2, 9 * 2
const TILE_SIZE int = 32
const WIDTH, HEIGHT int = TILE_COLUMNS * TILE_SIZE, TILE_ROWS * TILE_SIZE

const COLUMNS int = 64
const COLUMN_WIDTH float64 = float64(WIDTH) / float64(COLUMNS)

var player = struct {
	x     float32
	y     float32
	speed float32
}{
	x:     50,
	y:     50,
	speed: 1,
}

var world = [TILE_ROWS + 1][TILE_COLUMNS + 1]int{
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
}

func run() {

	cfg := opengl.WindowConfig{
		Title:     "Raycaster",
		Bounds:    pixel.R(0, 0, float64(WIDTH), float64(HEIGHT)),
		Resizable: false,
		VSync:     true,
	}
	win, err := opengl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	win.SetMatrix(pixel.IM.ScaledXY(pixel.ZV, pixel.V(1, -1)).Moved(pixel.V(0, float64(HEIGHT)))) // 0,0 now starts at the top left

	imd := imdraw.New(nil)

	for !win.Closed() {

		// @todo: Draw the colums for the raycaster instead of the tiles
		/*for i := 0; i < COLUMNS; i++ {
			imd.Color = color.RGBA{R: uint8(i * 5), G: 0, B: 0, A: 255}
			imd.Push(pixel.Vec{X: float64(i) * COLUMN_WIDTH, Y: float64(HEIGHT)}, pixel.Vec{X: (float64(i) * COLUMN_WIDTH) + COLUMN_WIDTH, Y: 0})
			imd.Rectangle(0)
		}*/

		////////////
		//// UPDATE
		////////////

		/* MOVMENT */
		// @todo: account for speed up when going diagnally, divide by 1.44
		if win.Pressed(pixel.KeyW) {
			player.y -= player.speed
		}
		if win.Pressed(pixel.KeyS) {
			player.y += player.speed
		}
		if win.Pressed(pixel.KeyA) {
			player.x -= player.speed
		}
		if win.Pressed(pixel.KeyD) {
			player.x += player.speed
		}

		////////////
		//// DRAW
		////////////

		/* LINES */
		imd.Clear()
		win.Clear(color.Black)

		imd.Color = color.RGBA{R: 255, G: 255, B: 255, A: 255}
		for x := 1; x < TILE_COLUMNS; x++ {
			imd.Push(pixel.V(float64(x)*float64(TILE_SIZE), 0), pixel.V(float64(x)*float64(TILE_SIZE), float64(HEIGHT)))
			imd.Line(1)
		}
		for y := 1; y < TILE_ROWS; y++ {
			imd.Push(pixel.V(0, float64(y)*float64(TILE_SIZE)), pixel.V(float64(WIDTH), float64(y)*float64(TILE_SIZE)))
			imd.Line(1)
		}

		/* WALLS */
		imd.Color = color.RGBA{R: 0, G: 0, B: 255, A: 255}
		for x := 0; x < TILE_COLUMNS; x++ {
			for y := TILE_ROWS; y > -1; y-- {
				if world[y][x] == 1 {
					imd.Push(pixel.V(float64(x*TILE_SIZE)-1, float64(y*TILE_SIZE)))
					imd.Push(pixel.V(float64(x*TILE_SIZE+TILE_SIZE), float64(y*TILE_SIZE+TILE_SIZE)+1))
					imd.Rectangle(0)
				}
			}
		}

		/* PLAYER */
		imd.Color = color.RGBA{R: 255, G: 0, B: 0, A: 255}
		imd.Push(pixel.V(float64(player.x), float64(player.y)))
		imd.Circle(3, 1)

		imd.Draw(win)

		win.Update()
	}
}

func main() {
	opengl.Run(run)
}
