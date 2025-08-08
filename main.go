package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	WIN_X                = 1000
	WIN_Y                = 600
	GRID_TILE_SIZE       = 25
	GRID_TILE_LINE_WIDTH = 1

	COLS = WIN_X / GRID_TILE_SIZE
	ROWS = WIN_Y / GRID_TILE_SIZE
)

const (
	DEAD  = 0
	ALIVE = 1
)

const (
	STOPPED = 0
	RUNNING = 1
)

type CellState int

type Cell struct {
	state CellState
}

type Game struct {
	zoom       float32
	cells      [ROWS][COLS]Cell
	generation int
	state      int
}

var game = Game{}

func initalize() {
	game.zoom = 1.0
	game.state = STOPPED

	for y := 0; y < ROWS; y++ {
		for x := 0; x < COLS; x++ {
			game.cells[y][x] = Cell{state: DEAD}
		}
	}
}

func spawn() {
	if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
		mousePos := rl.GetMousePosition()
		mousePosX := mousePos.X
		mousePosY := mousePos.Y

		x := int(mousePosX / GRID_TILE_SIZE)
		y := int(mousePosY / GRID_TILE_SIZE)

		game.cells[y][x].state = ALIVE
	}
}

func processGeneration() {
	nextGeneration()

	// if rl.IsKeyPressed(rl.KeyN) {
	//	nextGeneration()
	//	fmt.Println("Generation:", game.generation)
	// }
}

func nextGeneration() {
	new := game.cells

	for y := 0; y < ROWS; y++ {
		for x := 0; x < COLS; x++ {
			// rule 1: any live cell with fewer than two live neighbors dies (underpopulation)
			// rule 2: any live cell with two or three live neighbors lives on to the next generation
			// rule 3: any live cell with more than three neighbors dies (overpopulation)
			// rule 4: any dead cell with exactly three live neighbors becomes a live cell (reproduction)

			if x-1 < 0 || y-1 < 0 || x+1 >= COLS || y+1 >= ROWS {
				continue
			}

			neighboursAlive := 0
			for yy := y - 1; yy <= y+1; yy++ {
				for xx := x - 1; xx <= x+1; xx++ {
					if xx == x && yy == y {
						continue
					}

					if game.cells[yy][xx].state == ALIVE {
						neighboursAlive++
					}
				}
			}

			if neighboursAlive < 2 {
				new[y][x].state = DEAD
			} else if neighboursAlive == 3 {
				new[y][x].state = ALIVE
			} else if neighboursAlive < 4 {
				continue
			} else if neighboursAlive > 3 {
				new[y][x].state = DEAD
			}

		}
	}

	game.generation++
	game.cells = new
}

func processMousewheelInput() {
	zoom := game.zoom + rl.GetMouseWheelMove()*0.1

	if zoom > 3 || zoom < 0.5 {
		return
	} else {
		game.zoom = zoom
	}

}

func drawGrid() {
	scaledSize := float32(GRID_TILE_SIZE) * float32(game.zoom)

	for x := float32(0); x < float32(WIN_X); x += scaledSize {
		rl.DrawLineEx(rl.Vector2{X: float32(x), Y: 0}, rl.Vector2{X: float32(x), Y: WIN_Y}, GRID_TILE_LINE_WIDTH, rl.Black)
	}

	for y := float32(0); y < float32(WIN_Y); y += scaledSize {
		rl.DrawLineEx(rl.Vector2{X: 0, Y: float32(y)}, rl.Vector2{X: WIN_X, Y: float32(y)}, GRID_TILE_LINE_WIDTH, rl.Black)

	}
}

func drawCells() {
	for y := 0; y < ROWS; y++ {
		for x := 0; x < COLS; x++ {
			if game.cells[y][x].state == ALIVE {
				rl.DrawRectangle(int32(x*GRID_TILE_SIZE), int32(y*GRID_TILE_SIZE), GRID_TILE_SIZE, GRID_TILE_SIZE, rl.Black)
			}
		}
	}
}

func drawUI() {
	state := ""
	switch game.state {
	case STOPPED:
		state = "STOPPED"
	case RUNNING:
		state = "RUNNING"
	}

	genStr := fmt.Sprintf("Generation: %d", game.generation)
	stateStr := fmt.Sprintf("State: %s", state)
	fontSize := 24

	rl.DrawText(stateStr, WIN_X-rl.MeasureText(stateStr, int32(fontSize))-5, 5, int32(fontSize), rl.Black)
	rl.DrawText(genStr, WIN_X-rl.MeasureText(genStr, int32(fontSize))-5, 30, int32(fontSize), rl.Black)
}

func main() {
	initalize()
	rl.InitWindow(WIN_X, WIN_Y, "game of life")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)
	fc := 1

	for !rl.WindowShouldClose() {
		if rl.IsKeyPressed(rl.KeyF1) && game.state == STOPPED {
			game.state = RUNNING
		} else if rl.IsKeyPressed(rl.KeyF1) && game.state == RUNNING {
			game.state = STOPPED
		}

		if game.state == RUNNING {
			fc++
		}

		// processMousewheelInput()
		spawn()

		if fc%10 == 0 {
			processGeneration()
		}

		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		drawCells()
		// drawGrid()
		drawUI()

		rl.EndDrawing()
	}
}
