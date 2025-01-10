package main

import (
	"image/color"
	"math"

	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
	"github.com/gopxl/pixel/v2/ext/imdraw"
)

const TILE_COLUMNS, TILE_ROWS int = 16 * 2, 9 * 2
const TILE_SIZE float64 = 32
const WIDTH, HEIGHT int = TILE_COLUMNS * int(TILE_SIZE), TILE_ROWS * int(TILE_SIZE)

const COLUMNS int = 64
const COLUMN_WIDTH float64 = float64(WIDTH) / float64(COLUMNS)

var player = struct {
	x           float64
	y           float64
	dir         float64
	renderDis   float64
	speed       float64
	sensitivity float64
}{
	x:           50,
	y:           50,
	dir:         0,
	renderDis:   64,
	speed:       1,
	sensitivity: 1,
}

var world = [TILE_ROWS + 1][TILE_COLUMNS + 1]int{
	{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
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
	win.SetMatrix(pixel.IM.ScaledXY(pixel.ZV, pixel.V(1, -1)).Moved(pixel.V(0, float64(HEIGHT)))) // 0,0 is now at the top left

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

		/* CONTROLS */
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

		if win.Pressed(pixel.KeyLeft) {
			player.dir -= player.sensitivity
		}
		if win.Pressed(pixel.KeyRight) {
			player.dir += player.sensitivity
		}

		/* DIRECTION */

		if player.dir > 360 {
			player.dir = 0 + (player.dir - 360)
		} else if player.dir < 0 {
			player.dir = 360 + player.dir
		}

		rad := degreesToRadians(player.dir)
		x2 := player.x + player.renderDis*math.Cos(rad)
		y2 := player.y + player.renderDis*math.Sin(rad)

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
					imd.Push(pixel.V(float64(x*int(TILE_SIZE)), float64(y*int(TILE_SIZE))+1))
					imd.Push(pixel.V(float64(x*int(TILE_SIZE)+int(TILE_SIZE))-1, float64(y*int(TILE_SIZE)+int(TILE_SIZE))))
					imd.Rectangle(0)
				}
			}
		}

		/* PLAYER */
		imd.Color = color.RGBA{R: 255, G: 0, B: 0, A: 255}
		imd.Push(pixel.V(float64(player.x), float64(player.y)))
		imd.Circle(3, 1)

		// ray
		imd.Push(pixel.V(player.x, player.y), pixel.V(x2, y2))
		imd.Line(1)

		// collisions
		col := rayCollisions(pixel.L(pixel.V(player.x, player.y), pixel.V(x2, y2)))
		for _, p := range col {
			imd.Push(p)
			imd.Circle(3, 1)
		}

		imd.Draw(win)
		win.Update()
	}
}

func main() {
	opengl.Run(run)
}

// returns the end points of each line
func rayCollisions(line pixel.Line) []pixel.Vec {

	var colPoints []pixel.Vec

	//for _, line := range rays {
	// Vertical lines
	for x := 1; x < TILE_COLUMNS; x++ {
		point, hit := line.Intersect(pixel.L(pixel.V(float64(x)*float64(TILE_SIZE), 0), pixel.V(float64(x)*float64(TILE_SIZE), float64(HEIGHT))))

		if hit {
			colPoints = append(colPoints, point)
		}
	}

	// Horizontal lines
	for y := 1; y < TILE_ROWS; y++ {
		point, hit := line.Intersect(pixel.L(pixel.V(0, float64(y)*float64(TILE_SIZE)), pixel.V(float64(WIDTH), float64(y)*float64(TILE_SIZE))))

		if hit {
			colPoints = append(colPoints, point)
		}
	}
	//}

	return colPoints
}

func degreesToRadians(degrees float64) float64 {

	return degrees * math.Pi / 180
}
