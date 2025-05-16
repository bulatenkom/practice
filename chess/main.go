package main

import (
	"fmt"
	"os"
	"regexp"
)

type app struct {
	board     board // игровое поле с фигурами
	moveCount int   // счётчик ходов
	// текущее время
	// общее время партии
	ui uiInfo // элементы интерфейса
}

// gameloop()
func (app *app) gameloop() error {
	// prepare input matchers
	helpRegex := regexp.MustCompile(`^(?:\?|\/\?|help|h)$`)
	avMovesRegex := regexp.MustCompile(`^[a-h][1-8]$`)
	moveRegex := regexp.MustCompile(`^([a-h][1-8]){2}$`)

	for {
		// render UI
		if err := app.renderUI(); err != nil {
			return err
		}

		// user input
		var input string
		if _, err := fmt.Print("Input: "); err != nil {
			return err
		}
		if _, err := fmt.Scanf("%s", &input); err != nil {
			app.ui.err = err
			continue
		}

		switch {
		case helpRegex.MatchString(input):
			// toggle help
			app.ui.banner = "toggle help"
		case avMovesRegex.MatchString(input):
			// togle movements for piece
			app.ui.banner = "toggle movements"
		case moveRegex.MatchString(input):
			// handle move
			app.ui.banner = "handle move"
			app.moveCount++
		default:
			app.ui.err = fmt.Errorf("command '%s' is unknown. Type 'h' to get full list of commands", input)
			continue
		}
		// reset temp states
		app.ui.err = nil
	}
}

// renderUI()
func (app *app) renderUI() error {
	if _, err := clearScreen(); err != nil {
		return err
	}
	if _, err := fmt.Println(app.ui.banner); err != nil {
		return err
	}
	if err := renderBoard(app.board); err != nil {
		return err
	}
	if _, err := fmt.Println("move count:", app.moveCount); err != nil {
		return err
	}
	if app.ui.err != nil {
		if _, err := fmt.Println(app.ui.err); err != nil {
			return err
		}
	}
	return nil
}

func clearScreen() (int, error) {
	return fmt.Print("\x1b[H\x1b[2J")
}

func renderBoard(board board) error {
	output := ""
	for i := range board {
		for j := range board[0] {
			output += string(board[i][j])
		}
		output += "\n"
	}
	if _, err := fmt.Print(output); err != nil {
		return err
	}
	return nil
}

type board [][]rune

// movements() of piece in selected cell
// move() piece
// checkMove()
// updateBoard()
// checkRules() - check, checkmate, ...
// updateGameInfo()

type uiInfo struct {
	banner   string // элемент для вывода текстовой информации хода игры
	sideview string // элемент для вывода вспомогательной текстовой информации
	// история ходов
	err error // элемент для вывода ошибок пользовательского ввода
}

func main() {
	app := app{
		board: board{
			[]rune{'R', 'K', 'B', 'W', 'Q', 'B', 'K', 'R'},
			[]rune{'P', 'P', 'P', 'P', 'P', 'P', 'P', 'P'},
			[]rune{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '},
			[]rune{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '},
			[]rune{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '},
			[]rune{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '},
			[]rune{'P', 'P', 'P', 'P', 'P', 'P', 'P', 'P'},
			[]rune{'R', 'K', 'B', 'W', 'Q', 'B', 'K', 'R'},
		},
	}

	if err := app.gameloop(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
