package main

import (
	"fmt"
	"image/color"
	"math"
	"time"

	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
	"github.com/gopxl/pixel/v2/ext/imdraw"
)

const TILE_COLUMNS, TILE_ROWS int = 16 * 2, 9 * 2
const TILE_SIZE float64 = 32
const WIDTH, HEIGHT int = TILE_COLUMNS * int(TILE_SIZE), TILE_ROWS * int(TILE_SIZE)

const COLUMNS int = 256
const COLUMN_WIDTH float64 = float64(WIDTH) / float64(COLUMNS)

const FPS = 60

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
	renderDis:   1064,
	speed:       1.5,
	sensitivity: 1.5,
}

var world = [TILE_ROWS + 1][TILE_COLUMNS + 1]int{
	{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1},
	{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
}

func run() {

	cfg := opengl.WindowConfig{
		Title:     "Raycaster",
		Bounds:    pixel.R(0, 0, float64(WIDTH), float64(HEIGHT)),
		Resizable: false,
		VSync:     false,
	}
	win, err := opengl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	win.SetMatrix(pixel.IM.ScaledXY(pixel.ZV, pixel.V(1, -1)).Moved(pixel.V(0, float64(HEIGHT)))) // 0,0 is now at the top left

	imd := imdraw.New(nil)

	lastTime := 0
	frames := 0
	ticks := 0
	seconds := time.Tick(time.Second / FPS)

	var rays []pixel.Line

	debug := false // will draw the game top down instead of in person if true

	for !win.Closed() {

		select {
		case <-seconds:

			////////////
			//// UPDATE
			////////////

			/* CONTROLS */
			// @todo: account for speed up when going diagnally, divide by 1.44
			// @todo: use delta time

			playerRad := degreesToRadians(player.dir)
			if win.Pressed(pixel.KeyW) {
				player.x += math.Cos(playerRad) * player.speed
				player.y += math.Sin(playerRad) * player.speed
			}
			if win.Pressed(pixel.KeyS) {
				player.x -= math.Cos(playerRad) * player.speed
				player.y -= math.Sin(playerRad) * player.speed
			}
			if win.Pressed(pixel.KeyA) {
				player.x += math.Cos(playerRad-math.Pi/2) * player.speed
				player.y += math.Sin(playerRad-math.Pi/2) * player.speed
			}
			if win.Pressed(pixel.KeyD) {
				player.x += math.Cos(playerRad+math.Pi/2) * player.speed
				player.y += math.Sin(playerRad+math.Pi/2) * player.speed
			}

			if win.Pressed(pixel.KeyLeft) {
				player.dir -= player.sensitivity
			}
			if win.Pressed(pixel.KeyRight) {
				player.dir += player.sensitivity
			}

			// @debug
			if win.Pressed(pixel.KeySpace) {
				debug = debug == false
			}

			/* DIRECTION */
			rays = rays[:0]
			gap := .25
			for i := -COLUMNS / 2; i < COLUMNS/2; i++ {
				rad := degreesToRadians(player.dir + float64(i)*gap)
				x := player.x + player.renderDis*math.Cos(rad)
				y := player.y + player.renderDis*math.Sin(rad)

				hit, ray := rayCollisions(pixel.L(pixel.V(player.x, player.y), pixel.V(x, y)))
				if hit {
					rays = append(rays, pixel.L(pixel.V(player.x, player.y), ray))
				} else {
					// @cleanup This's accounting for the vertical line issue, it should panic
					rays = append(rays, pixel.L(pixel.V(player.x, player.y), pixel.V(player.renderDis, player.renderDis)))
				}
			}

			////////////
			//// DRAW
			////////////

			imd.Clear()
			win.Clear(color.RGBA{R: 0, G: 0, B: 0, A: 255})

			if debug {

				/* LINES */
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
				// rays
				for _, ray := range rays {
					imd.Push(ray.A, ray.B)
					imd.Line(1)
				}

			} else {

				/* FIRST PERSON */
				for i := 0; i < COLUMNS; i++ {
					distance := rays[i].Len() / (float64(TILE_COLUMNS)) // in tiles
					distance *= 25

					if distance >= float64(HEIGHT)/2-10 {
						distance = float64(HEIGHT)/2 - 10
					}

					imd.Color = color.White
					imd.Push(pixel.V(float64(i)*COLUMN_WIDTH, 0))
					imd.Push(pixel.V(float64(i)*COLUMN_WIDTH+COLUMN_WIDTH, float64(HEIGHT)))
					imd.Rectangle(0)

					imd.Color = color.Black
					// top
					imd.Push(pixel.V(float64(i)*COLUMN_WIDTH, 0))
					imd.Push(pixel.V(float64(i)*COLUMN_WIDTH+COLUMN_WIDTH, distance))
					imd.Rectangle(0)
					// bottom
					imd.Push(pixel.V(float64(i)*COLUMN_WIDTH, float64(HEIGHT)-distance))
					imd.Push(pixel.V(float64(i)*COLUMN_WIDTH+COLUMN_WIDTH, float64(HEIGHT)))
					imd.Rectangle(0)
				}
			}

			imd.Draw(win)
			frames++
		default:
		}

		win.Update()

		// @todo FPS doesn't print out after awhile
		ticks++
		if time.Now().Second() >= lastTime+1 {
			fmt.Printf("FPS: %d | Ticks: %d\n", frames, ticks)
			lastTime = time.Now().Second()
			frames = 0
			ticks = 0
		}
	}
}

// returns true and the end point of the line if a wall is hit, if not it'll return false and a zero vector
func rayCollisions(line pixel.Line) (bool, pixel.Vec) {
	// @optimize: should check each peice of the grid sequentially based on where the ray is going
	// For now it goes through each tile to see if it's solid, then checks to see if a ray is hitting it

	var colPoints []pixel.Vec

	for y := 0; y < TILE_ROWS; y++ {
		for x := 0; x < TILE_COLUMNS; x++ {
			tile := pixel.V(float64(x), float64(y))

			if isSolidTile(tile) {
				tileLeft := pixel.L(pixel.V(tile.X*TILE_SIZE, tile.Y*TILE_SIZE), pixel.V(tile.X*TILE_SIZE, (tile.Y+1)*TILE_SIZE))
				tileRight := pixel.L(pixel.V((tile.X+1)*TILE_SIZE, tile.Y*TILE_SIZE), pixel.V((tile.X+1)*TILE_SIZE, (tile.Y+1)*TILE_SIZE))
				tileTop := pixel.L(pixel.V(tile.X*TILE_SIZE, (tile.Y+1)*TILE_SIZE), pixel.V((tile.X+1)*TILE_SIZE, (tile.Y+1)*TILE_SIZE))
				tileBottom := pixel.L(pixel.V(tile.X*TILE_SIZE, tile.Y*TILE_SIZE), pixel.V((tile.X+1)*TILE_SIZE, tile.Y*TILE_SIZE))

				// @todo Fix the vertical line rendering issue, line.Intersect doesn't register vertical lines
				for _, edge := range []pixel.Line{tileLeft, tileRight, tileTop, tileBottom} {
					point, hit := line.Intersect(edge)
					if hit {
						colPoints = append(colPoints, point)
					}
				}
			}
		}
	}

	shortestDistance := player.renderDis
	var hitPoint pixel.Vec
	hit := false

	for _, point := range colPoints {
		pDistance := pixel.L(pixel.V(player.x, player.y), point).Len()

		if pDistance < shortestDistance {
			shortestDistance = pDistance
			hitPoint = point
			hit = true
		}
	}
	if hit {
		return true, hitPoint
	}

	return false, pixel.ZV
}

func CoordToTile(v pixel.Vec) pixel.Vec {
	return pixel.Vec{
		X: math.Floor(v.X / TILE_SIZE),
		Y: math.Floor(v.Y / TILE_SIZE),
	}
}

func isSolidTile(tile pixel.Vec) bool {
	return world[int(tile.Y)][int(tile.X)] == 1
}

func degreesToRadians(degrees float64) float64 {

	return degrees * math.Pi / 180
}

func main() {
	opengl.Run(run)
}
