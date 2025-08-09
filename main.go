package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	WIN_X                = 1600 // game area
	WIN_Y                = 1000 // game area
	UI_Y                 = 200
	GRID_TILE_SIZE       = 20
	GRID_TILE_LINE_WIDTH = 1

	COLS = WIN_X / GRID_TILE_SIZE
	ROWS = WIN_Y / GRID_TILE_SIZE
)

const (
	DEAD  = 0
	ALIVE = 1
	NEVER = 3 // state for never touched cells
)

const (
	STOPPED = 0
	RUNNING = 1
)

type CellState int

type Cell struct {
	state    CellState
	gensDead int
}

type Theme struct {
	bg        rl.Color
	cell      rl.Color
	grid      rl.Color
	lastAlive [4]rl.Color
}

type Game struct {
	zoom                  float32
	font                  rl.Font
	speed                 int
	cells                 [ROWS][COLS]Cell
	generation            int
	state                 int
	buttons               []Button
	fc                    int // frame counter
	theme                 Theme
	lastAliveColorEnabled bool
	gridEnabled           bool
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
	game.fc = 1
	game.lastAliveColorEnabled = false
	game.gridEnabled = true

	game.theme = Theme{
		bg:   rl.White,
		cell: rl.Black,
		grid: rl.Black,
		lastAlive: [4]rl.Color{
			{
				R: 102,
				G: 255,
				B: 51,
				A: 150,
			},
			{
				R: 133,
				G: 214,
				B: 41,
				A: 150,
			},
			{
				R: 224,
				G: 92,
				B: 10,
				A: 150,
			},
			{
				R: 255,
				G: 51,
				B: 0,
				A: 150,
			},
		},
	}

	for y := 0; y < ROWS; y++ {
		for x := 0; x < COLS; x++ {
			game.cells[y][x] = Cell{state: NEVER}
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
		{
			pos:        rl.Vector2{X: 500, Y: WIN_Y + 10},
			size:       rl.Vector2{X: 250, Y: 50},
			color:      rl.Gray,
			hoverColor: rl.LightGray,
			fontColor:  rl.Black,
			text:       "Last Alive Color",
			action:     toggleLastAliveColor,
		},
		{
			pos:        rl.Vector2{X: 500, Y: WIN_Y + 10 + 60},
			size:       rl.Vector2{X: 250, Y: 50},
			color:      rl.Gray,
			hoverColor: rl.LightGray,
			fontColor:  rl.Black,
			text:       "Grid",
			action:     toggleGrid,
		},
	}
}

func toggleGrid() {
	game.gridEnabled = !game.gridEnabled
}

func toggleLastAliveColor() {
	game.lastAliveColorEnabled = !game.lastAliveColorEnabled
}

func increaseSpeed() {
	if game.speed < 10 {
		game.speed++
		return
	}

	new := game.speed + 10

	if new > 120 {
		return
	}

	game.speed = new
}

func decreaseSpeed() {
	new := game.speed - 10

	if new < 2 {
		new = game.speed - 1

		if new < 2 {
			return
		}
	}

	game.speed = new
}

func spawn() {
	if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
		mousePos := rl.GetMousePosition()
		mousePosX := mousePos.X
		mousePosY := mousePos.Y

		x := int(mousePosX / GRID_TILE_SIZE * game.zoom)
		y := int(mousePosY / GRID_TILE_SIZE * game.zoom)

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

			aboveX := x - 1
			sameX := x
			belowX := x + 1

			aboveY := y - 1
			sameY := y
			belowY := y + 1

			if x-1 < 0 || y-1 < 0 || x+1 >= COLS || y+1 >= ROWS {
				continue
			}

			neighboursAlive := 0
			for yy := aboveY; yy <= belowY; yy++ {
				for xx := aboveX; xx <= belowX; xx++ {
					if xx == sameX && yy == sameY {
						continue
					}

					if game.cells[yy][xx].state == ALIVE {
						neighboursAlive++
					}
				}
			}

			if neighboursAlive < 2 {
				if new[y][x].state == NEVER {
					new[y][x].state = NEVER
				} else {
					new[y][x].state = DEAD
				}
			} else if neighboursAlive == 3 {
				new[y][x].state = ALIVE
				new[y][x].gensDead = 0
			} else if neighboursAlive < 4 {
				continue
			} else if neighboursAlive > 3 {
				if new[y][x].state == NEVER {
					new[y][x].state = NEVER
				} else {
					new[y][x].state = DEAD
				}
			}

			if new[y][x].state == DEAD {
				new[y][x].gensDead++
			}

		}
	}

	game.generation++
	game.cells = new
}

// func processMousewheelInput() {
//	zoom := game.zoom + rl.GetMouseWheelMove()*0.1

//	if zoom > 3 || zoom < 0.5 {
//		return
//	} else {
//		game.zoom = zoom
//	}

// }

func drawGrid() {
	if !game.gridEnabled {
		return
	}

	scaledSize := float32(GRID_TILE_SIZE) * float32(game.zoom)

	for x := float32(0); x < float32(WIN_X); x += scaledSize {
		rl.DrawLineEx(rl.Vector2{X: float32(x), Y: 0}, rl.Vector2{X: float32(x), Y: WIN_Y}, GRID_TILE_LINE_WIDTH, game.theme.grid)
	}

	for y := float32(0); y < float32(WIN_Y); y += scaledSize {
		rl.DrawLineEx(rl.Vector2{X: 0, Y: float32(y)}, rl.Vector2{X: WIN_X, Y: float32(y)}, GRID_TILE_LINE_WIDTH, game.theme.grid)

	}
}

func drawCells() {
	for y := 0; y < ROWS; y++ {
		for x := 0; x < COLS; x++ {
			if game.cells[y][x].state == DEAD {
				if game.lastAliveColorEnabled {

					if game.cells[y][x].gensDead < 25 {
						rl.DrawRectangleV(
							rl.Vector2{
								X: float32(float32(x) * GRID_TILE_SIZE * game.zoom),
								Y: float32(float32(y) * GRID_TILE_SIZE * game.zoom),
							},
							rl.Vector2{
								X: GRID_TILE_SIZE * game.zoom,
								Y: GRID_TILE_SIZE * game.zoom,
							},
							game.theme.lastAlive[0],
						)
					} else if game.cells[y][x].gensDead >= 25 && game.cells[y][x].gensDead < 50 {
						rl.DrawRectangleV(
							rl.Vector2{
								X: float32(float32(x) * GRID_TILE_SIZE * game.zoom),
								Y: float32(float32(y) * GRID_TILE_SIZE * game.zoom),
							},
							rl.Vector2{
								X: GRID_TILE_SIZE * game.zoom,
								Y: GRID_TILE_SIZE * game.zoom,
							},
							game.theme.lastAlive[1],
						)

					} else if game.cells[y][x].gensDead >= 50 && game.cells[y][x].gensDead < 75 {
						rl.DrawRectangleV(
							rl.Vector2{
								X: float32(float32(x) * GRID_TILE_SIZE * game.zoom),
								Y: float32(float32(y) * GRID_TILE_SIZE * game.zoom),
							},
							rl.Vector2{
								X: GRID_TILE_SIZE * game.zoom,
								Y: GRID_TILE_SIZE * game.zoom,
							},
							game.theme.lastAlive[2],
						)

					} else if game.cells[y][x].gensDead >= 75 {
						rl.DrawRectangleV(
							rl.Vector2{
								X: float32(float32(x) * GRID_TILE_SIZE * game.zoom),
								Y: float32(float32(y) * GRID_TILE_SIZE * game.zoom),
							},
							rl.Vector2{
								X: GRID_TILE_SIZE * game.zoom,
								Y: GRID_TILE_SIZE * game.zoom,
							},
							game.theme.lastAlive[3],
						)

					}
				}
			}

			if game.cells[y][x].state == ALIVE {
				rl.DrawRectangleV(
					rl.Vector2{
						X: float32(float32(x) * GRID_TILE_SIZE * game.zoom),
						Y: float32(float32(y) * GRID_TILE_SIZE * game.zoom),
					},
					rl.Vector2{
						X: GRID_TILE_SIZE * game.zoom,
						Y: GRID_TILE_SIZE * game.zoom,
					},
					game.theme.cell,
				)
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

	game.fc = 1
}

func main() {
	rl.InitWindow(WIN_X, WIN_Y+UI_Y, "game of life")
	defer rl.CloseWindow()

	initalize()

	rl.SetTargetFPS(120)

	for !rl.WindowShouldClose() {

		if game.state == RUNNING {
			game.fc++
		}

		// processMousewheelInput()
		spawn()

		if game.fc%game.speed == 0 {
			processGeneration()
		}

		rl.BeginDrawing()

		rl.ClearBackground(game.theme.bg)

		drawCells()
		drawGrid()
		drawUI()

		rl.EndDrawing()
	}
}
