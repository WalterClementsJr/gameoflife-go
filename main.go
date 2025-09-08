package main

import rl "github.com/gen2brain/raylib-go/raylib"

const (
	SquareSize = 20
	ScreenW    = 1400
	ScreenH    = 920
)

type Game struct {
	screenW       int32
	screenH       int32
	fps           int32
	liveCellCount uint
	pixelSize     int32

	grid [][]int32
}

func (game *Game) init() {
	// init landGrid
	game.screenW = ScreenW
	game.screenH = ScreenH
	game.pixelSize = SquareSize
	game.fps = 5

	game.grid = make([][]int32, int(game.screenW/game.pixelSize))

	for i := range game.grid {
		game.grid[i] = make([]int32, int(game.screenH/2))
	}

	rl.InitWindow(game.screenW, game.screenH, "Game of Life")
	rl.SetTargetFPS(game.fps)
}

func draw(game Game) {
	rl.ClearBackground(rl.Gray)
	drawGrid(game)
}

func drawGrid(game Game) {
	for i := int32(1); i <= game.screenW; i += game.pixelSize {
		rl.DrawLine(i, 0, i, game.screenH, rl.Red)
		rl.DrawLine(0, i, game.screenW, i, rl.Red)
	}
}

func update(game *Game) {
}

func main() {
	g := Game{}
	g.init()

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		draw(g)
		update(&g)

		rl.EndDrawing()
	}
}
