package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	WIN_X                = 1000 // game area
	WIN_Y                = 600  // game area
	UI_Y                 = 200
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
	font       rl.Font
	speed      int
	cells      [ROWS][COLS]Cell
	generation int
	state      int
	buttons    []Button
}

type Button struct {
	pos        rl.Vector2
	size       rl.Vector2
	color      rl.Color
	hoverColor rl.Color
	fontColor  rl.Color
	text       string
	action     func()
}

var game = Game{}

func initalize() {
	game.zoom = 1.0
	game.state = STOPPED
	game.font = rl.LoadFont("./font/static/OpenSans-Bold.ttf")
	game.speed = 60

	for y := 0; y < ROWS; y++ {
		for x := 0; x < COLS; x++ {
			game.cells[y][x] = Cell{state: DEAD}
		}
	}

	game.buttons = []Button{
		{
			pos:        rl.Vector2{X: 5, Y: WIN_Y + 10},
			size:       rl.Vector2{X: 200, Y: 50},
			color:      rl.Gray,
			hoverColor: rl.LightGray,
			fontColor:  rl.Black,
			text:       "Run/Stop",
			action:     changeGameState,
		},
		{
			pos:        rl.Vector2{X: 5, Y: WIN_Y + 10 + 60},
			size:       rl.Vector2{X: 90, Y: 50},
			color:      rl.Gray,
			hoverColor: rl.LightGray,
			fontColor:  rl.Black,
			text:       "+",
			action:     increaseSpeed,
		},
		{
			pos:        rl.Vector2{X: 115, Y: WIN_Y + 10 + 60},
			size:       rl.Vector2{X: 90, Y: 50},
			color:      rl.Gray,
			hoverColor: rl.LightGray,
			fontColor:  rl.Black,
			text:       "-",
			action:     decreaseSpeed,
		},
		{
			pos:        rl.Vector2{X: 5, Y: WIN_Y + 10 + 60*2},
			size:       rl.Vector2{X: 200, Y: 50},
			color:      rl.Gray,
			hoverColor: rl.LightGray,
			fontColor:  rl.Black,
			text:       "Next Gen.",
			action:     nextGeneration,
		},
	}
}

func increaseSpeed() {
	new := game.speed + 10

	if new > 120 {
		return
	}

	game.speed = new
}

func decreaseSpeed() {
	new := game.speed - 10

	if new < 1 {
		return
	}

	game.speed = new
}

func spawn() {
	if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
		mousePos := rl.GetMousePosition()
		mousePosX := mousePos.X
		mousePosY := mousePos.Y

		x := int(mousePosX / GRID_TILE_SIZE)
		y := int(mousePosY / GRID_TILE_SIZE)

		if x >= COLS || y >= ROWS {
			return
		}

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

func drawButtons() {
	mousePos := rl.GetMousePosition()
	mouseX := mousePos.X
	mouseY := mousePos.Y
	var fontSize float32 = 32.0
	var spacing float32 = 2.0

	for _, button := range game.buttons {
		if mouseX > button.pos.X && mouseX < button.pos.X+button.size.X &&
			mouseY > button.pos.Y && mouseY < button.pos.Y+button.size.Y {
			rl.DrawRectangleV(button.pos, button.size, button.hoverColor)

			if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
				button.action()
			}
		} else {
			rl.DrawRectangleV(button.pos, button.size, button.color)
		}

		rl.DrawTextEx(
			game.font,
			button.text,
			rl.Vector2{
				X: button.pos.X + button.size.X/2 - rl.MeasureTextEx(game.font, button.text, fontSize, spacing).X/2,
				Y: button.pos.Y + button.size.Y/2 - rl.MeasureTextEx(game.font, button.text, fontSize, spacing).Y/2,
			},
			fontSize,
			spacing,
			button.fontColor,
		)
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
	speedStr := fmt.Sprintf("Speed: %d", game.speed)
	var fontSize float32 = 32.0
	var spacing float32 = 2.0

	// top line in ui area
	rl.DrawLineEx(rl.Vector2{X: 0, Y: WIN_Y}, rl.Vector2{X: WIN_X, Y: WIN_Y}, 3.0, rl.Black)

	rl.DrawTextEx(
		game.font,
		stateStr,
		rl.Vector2{
			X: 215,
			Y: WIN_Y + 5,
		},
		fontSize,
		spacing,
		rl.Black)

	rl.DrawTextEx(
		game.font,
		speedStr,
		rl.Vector2{
			X: 215,
			Y: WIN_Y + 5 + 60,
		},
		fontSize,
		spacing,
		rl.Black)

	rl.DrawTextEx(
		game.font,
		genStr,
		rl.Vector2{
			X: 215,
			Y: WIN_Y + 5 + 60*2,
		},
		fontSize,
		spacing,
		rl.Black)

	drawButtons()
}

func changeGameState() {
	switch game.state {
	case STOPPED:
		game.state = RUNNING
	case RUNNING:
		game.state = STOPPED
	}
}

func main() {
	rl.InitWindow(WIN_X, WIN_Y+UI_Y, "game of life")
	defer rl.CloseWindow()

	initalize()

	rl.SetTargetFPS(120)
	fc := 1

	for !rl.WindowShouldClose() {

		if game.state == RUNNING {
			fc++
		}

		// processMousewheelInput()
		spawn()

		if fc%game.speed == 0 {
			processGeneration()
		}

		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		drawCells()
		drawGrid()
		drawUI()

		rl.EndDrawing()
	}
}
