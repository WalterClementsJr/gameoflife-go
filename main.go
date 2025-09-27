package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	rg "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	ZoomLevel  = 20
	SquareSize = 40
	ScreenW    = 1400
	ScreenH    = 920
)

var CellsColor = rl.Black

type Game struct {
	screenW     int32
	screenH     int32
	fps         int32
	refreshRate int32
	pixelSize   int32
	zoomFactor  int

	refreshChan chan bool

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
	game.refreshChan = make(chan bool)
	game.refreshRate = 5

	timer := time.NewTicker(time.Millisecond * time.Duration(1000/game.refreshRate))

	// emit refresh event
	go func() {
		for {
			<-timer.C
			log.Println("engine refresh now", time.Now())
			game.refreshChan <- true
		}
	}()

	game.fps = 30
	game.pause = true
	game.generation = 0
	game.zoomFactor = 10

	game.grid = initGrid(*game)
	// default pattern
	game.grid[3][4] = 1
	game.grid[3][5] = 1
	game.grid[4][4] = 1
}

func initGrid(game Game) [][]int32 {
	grid := make([][]int32, int(game.screenW))

	for i := range grid {
		w := int(game.screenH)

		grid[i] = make([]int32, w)
	}
	return grid
}

func draw(game *Game) {
	rl.ClearBackground(rl.RayWhite)

	drawCells(*game)
	drawGrid(*game)
	drawUI(game)
	drawCustomCells(game)
}

func drawCells(game Game) {
	for i := range game.grid {
		for j := range game.grid[i] {
			if game.grid[i][j] == 1 {
				rl.DrawRectangle(int32(i)*int32(game.zoomFactor), int32(j)*int32(game.zoomFactor), int32(game.zoomFactor), int32(game.zoomFactor), CellsColor)
			}
		}
	}
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

	scroll := int(rl.GetMouseWheelMove())

	game.zoomFactor -= scroll
	if game.zoomFactor > ZoomLevel {
		game.zoomFactor = ZoomLevel
	}
	if game.zoomFactor < 2 {
		game.zoomFactor = 2
	}

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
	rl.DrawText(fmt.Sprintf("Generation: %d", game.generation), 100, 80+20*2, 30, rl.Red)
	rl.DrawText(fmt.Sprintf("Zoom level: %d", game.zoomFactor), 100, 80+20*4, 30, rl.Red)
}

// mouse draw
func drawCustomCells(game *Game) {
	mouse := rl.IsMouseButtonDown(rl.MouseButtonLeft)

	if mouse {
		pos := rl.GetMousePosition()

		tx := int(pos.X / float32(game.zoomFactor))
		ty := int(pos.Y / float32(game.zoomFactor))

		game.grid[tx][ty] = 1
	}
}

func drawGrid(game Game) {
	for i := int32(1); i <= game.screenW; i += int32(game.zoomFactor) {
		rl.DrawLine(i, 0, i, game.screenH, rl.Black)
		rl.DrawLine(0, i, game.screenW, i, rl.Black)
	}
}

func gameUpdate(game *Game) {
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

	var liveCount uint = 0
	for i := range game.grid {
		for j := range game.grid[i] {
			// populated
			liveCount := countLiveNeighbor(game.grid, i, j)

			if game.grid[i][j] == 1 {
				if liveCount <= 1 || liveCount >= 4 {
					currentGrid[i][j] = 0 // kill
				} else {
					currentGrid[i][j] = 1 // spaw
					liveCount++
				}
			} else {
				if liveCount == 3 {
					currentGrid[i][j] = 1 // spawn
					liveCount++
				}
			}
		}
	}
	game.liveCellCount = liveCount
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

	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, syscall.SIGINT, syscall.SIGTERM)
	stop := make(chan bool)

	update := func() {
		rl.BeginDrawing()
		draw(&game)
		rl.EndDrawing()
	}

	// listen to OS signal
	go func() {
		<-osSignal
		stop <- true
	}()
	// ignore stop if receive stop signal
	go func() {
		<-stop
		log.Println("Stop signal received")
		os.Exit(0)
	}()

	go func() {
		for {
			<-game.refreshChan
			gameUpdate(&game)
		}
	}()
	// draw cells concurrently, signaled by channel
	for !rl.WindowShouldClose() {
		update()
	}
}
