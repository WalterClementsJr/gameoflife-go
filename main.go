package main

import (
	"fmt"

	rg "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	SquareSize = 20
	ScreenW    = 1400
	ScreenH    = 920
)

type Game struct {
	screenW       int32
	screenH       int32
	fps           int32
	pause         bool
	liveCellCount int
	pixelSize     int32

	grid [][]int32
}

func (game *Game) init() {
	// init landGrid
	game.screenW = ScreenW
	game.screenH = ScreenH
	game.pixelSize = SquareSize
	game.fps = 10

	game.grid = make([][]int32, int(game.screenW/game.pixelSize))

	for i := range game.grid {
		game.grid[i] = make([]int32, int(game.screenH/2))
	}

	rl.InitWindow(game.screenW, game.screenH, "Game of Life")
	rl.SetTargetFPS(game.fps)
}

func draw(game *Game) {
	rl.ClearBackground(rl.RayWhite)

	drawGrid(*game)
	drawUI(game)
}

func drawUI(game *Game) {
	// draw pause button
	pauseX, pauseY, pauseW, pauseH := 400, 500, 80, 50

	pauseBtn := rl.Rectangle{X: float32(pauseX), Y: float32(pauseY), Width: float32(pauseW), Height: float32(pauseH)}

	var pauseState string
	fmt.Println("game pause:", game.pause)
	if game.pause {
		pauseState = "Start"
	} else {
		pauseState = "Pause"
	}

	pause := rg.Button(pauseBtn, pauseState)
	if pause {
		game.pause = !game.pause
	}

	rl.DrawText(
		fmt.Sprintf("Live cells: %d", game.liveCellCount),
		200, 80, 20, rl.Red,
	)
}

func drawGrid(game Game) {
	for i := int32(1); i <= game.screenW; i += game.pixelSize {
		rl.DrawLine(i, 0, i, game.screenH, rl.Black)
		rl.DrawLine(0, i, game.screenW, i, rl.Black)
	}
}

func update(game *Game) {
	if game.pause {
		return
	}
}

func main() {
	g := Game{}
	g.init()

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		draw(&g)
		update(&g)

		rl.EndDrawing()
	}
}
