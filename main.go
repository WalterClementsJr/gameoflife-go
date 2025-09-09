package main

import (
	"fmt"

	rg "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	SquareSize = 40
	ScreenW    = 1400
	ScreenH    = 920
)

var CellsColor = rl.Black

type Game struct {
	screenW   int32
	screenH   int32
	fps       int32
	pixelSize int32

	pause    bool
	stepOver bool

	liveCellCount uint
	generation    uint

	grid [][]int32 // grid of cells data
}

func (game *Game) init() {
	// init landGrid
	game.screenW = ScreenW
	game.screenH = ScreenH
	game.pixelSize = SquareSize

	game.fps = 10
	game.pause = true
	game.generation = 0

	game.grid = initGrid(*game)
	game.grid[3][4] = 1
	game.grid[3][5] = 1
}

func initGrid(game Game) [][]int32 {
	grid := make([][]int32, int(game.screenW/game.pixelSize))

	for i := range grid {
		w := int(game.screenH / game.pixelSize)

		grid[i] = make([]int32, w)
	}
	return grid
}

func draw(game *Game) {
	rl.ClearBackground(rl.RayWhite)

	drawCells(game)
	drawGrid(*game)
	drawUI(game)
	drawCustomCells(game)
}

func drawCells(game *Game) {
	var cellCount uint
	for i := range game.grid {
		for j := range game.grid[i] {
			if game.grid[i][j] == 1 {
				rl.DrawRectangle(int32(i)*game.pixelSize, int32(j)*game.pixelSize, game.pixelSize, game.pixelSize, CellsColor)
				cellCount++
			}
		}
	}
	game.liveCellCount = cellCount
}

func drawUI(game *Game) {
	// draw pause button
	pauseX, pauseY, pauseW, pauseH := 400, game.screenH-40, 80, 30
	stepX, stepY, stepW, stepH := 500, game.screenH-40, 80, 30

	pauseBtn := rl.Rectangle{X: float32(pauseX), Y: float32(pauseY), Width: float32(pauseW), Height: float32(pauseH)}
	stepBtn := rl.Rectangle{X: float32(stepX), Y: float32(stepY), Width: float32(stepW), Height: float32(stepH)}
	restartRect := rl.Rectangle{X: 600.0, Y: float32(game.screenH - 40), Width: 80.0, Height: 30.0}

	var pauseState string
	if game.pause {
		pauseState = "Start"
	} else {
		pauseState = "Pause"
	}

	pause := rg.Button(pauseBtn, pauseState)
	stepOver := rg.Button(stepBtn, "Step Over")
	restart := rg.Button(restartRect, "Restart")

	if pause {
		game.pause = !game.pause
	}
	if stepOver {
		game.stepOver = true
		game.pause = false
	}
	if restart {
		game.init()
	}
	rl.DrawText(fmt.Sprintf("Live cells: %d", game.liveCellCount), 100, 80, 30, rl.Red)
	rl.DrawText(fmt.Sprintf("Generation: %d", game.generation), 100, 120, 30, rl.Red)
}

func drawCustomCells(game *Game) {
	mouse := rl.IsMouseButtonDown(rl.MouseButtonLeft)

	if mouse {
		pos := rl.GetMousePosition()

		tx := int(pos.X / float32(game.pixelSize))
		ty := int(pos.Y / float32(game.pixelSize))

		game.grid[tx][ty] = 1

	}
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
	game.generation++
	// For a space that is populated:
	//     Each cell with one or no neighbors dies, as if by solitude.
	//     Each cell with four or more neighbors dies, as if by overpopulation.
	//     Each cell with two or three neighbors survives.
	//
	// For a space that is empty or unpopulated:
	//     Each cell with three neighbors becomes populated.

	currentGrid := initGrid(*game)

	for i := range game.grid {
		for j := range game.grid[i] {
			// populated
			liveCount := countLiveNeighbor(game.grid, i, j)

			if game.grid[i][j] == 1 {
				if liveCount <= 1 || liveCount >= 4 {
					currentGrid[i][j] = 0 // kill
				} else {
					currentGrid[i][j] = 1 // kill
				}
			} else {
				if liveCount == 3 {
					currentGrid[i][j] = 1 // spawn
				}
			}
		}
	}
	game.grid = currentGrid

	if game.stepOver {
		game.stepOver = false
		game.pause = true
	}
}

func countLiveNeighbor(grid [][]int32, cellx, celly int) int {
	gridW := len(grid)
	gridH := len(grid[0])
	count := 0

	for x := -1; x <= 1; x++ {
		for y := -1; y <= 1; y++ {
			if x == y && x == 0 {
				continue
			}

			tx := (x + cellx + gridW) % gridW
			ty := (y + celly + gridH) % gridH
			count += int(grid[tx][ty])
		}
	}

	return count
}

func main() {
	game := Game{}
	game.init()

	rl.InitWindow(game.screenW, game.screenH, "Game of Life")
	rl.SetTargetFPS(game.fps)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		draw(&game)
		update(&game)

		rl.EndDrawing()
	}
}
