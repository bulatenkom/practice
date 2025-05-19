package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
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
	histRegex := regexp.MustCompile(`^(hist|history)$`)

	for {
		// paint UI
		if err := app.paintUI(); err != nil {
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
			app.ui.sideview = helpTemplate()
		case histRegex.MatchString(input):
			// toggle help
			app.ui.banner = "toggle history"
			app.ui.sideview = app.ui.history
		case avMovesRegex.MatchString(input):
			// togle movements for piece
			app.ui.banner = "toggle movements"
			app.ui.sideview = "(moves goes here...)"
		case moveRegex.MatchString(input):
			// handle move
			if input[:2] == input[2:] {
				app.ui.err = fmt.Errorf("move '%s' is not possible. Piece cannot move to the same cell", input)
				continue
			}

			app.ui.banner = "handle move"
			app.moveCount++
			app.ui.history += fmt.Sprintf(" %v. %v \n", app.moveCount, input)
			app.ui.sideview = app.ui.history
		default:
			app.ui.err = fmt.Errorf("command '%s' is unknown. Type 'h' to get full list of commands", input)
			continue
		}
		// reset temp states
		app.ui.err = nil
	}
}

func helpTemplate() string {
	return `HELP
  help, h, ?, /? - toggles help window
  [a-h][1-8] - prints possible moves for selected piece (e.g. 'e2' prints moves for Pawn)
  [a-h][1-8][a-h][1-8] - move piece (e.g. 'e2e4' moves piece from e2 to e4)
`
}

// paintUI()
func (app *app) paintUI() error {
	if _, err := ClearScreen(); err != nil {
		return err
	}
	mainview := Layout(
		[]Box{
			NewBox(AnsiYellow(app.ui.banner)),
			NewBox(RenderBoard(app.board)),
			NewBox("move count: " + strconv.Itoa(app.moveCount)),
		},
		Vertical,
	)
	sideview := Layout(
		[]Box{
			NewBox(app.ui.sideview),
			// NewBox(ansiRed(app.ui.sideview)),
		},
		Vertical,
	)
	layout := Layout([]Box{mainview, sideview}, Horizontal)
	// layout := Layout([]Box{sideview}, Horizontal)
	if _, err := fmt.Println(layout.String()); err != nil {
		return err
	}
	if app.ui.err != nil {
		if _, err := fmt.Println(AnsiRed(app.ui.err.Error())); err != nil {
			return err
		}
	}
	return nil
}

type uiInfo struct {
	banner   string // элемент для вывода текстовой информации хода игры
	sideview string // элемент для вывода вспомогательной текстовой информации
	history  string // история ходов
	err      error  // элемент для вывода ошибок пользовательского ввода
}

type board [][]rune

// movements() of piece in selected cell
// move() piece
// checkMove()
// updateBoard()
// checkRules() - check, checkmate, ...
// updateGameInfo()

func main() {
	app := app{
		board: board{
			[]rune{'♜', '♞', '♝', '♛', '♚', '♝', '♞', '♜'},
			[]rune{'♟', '♟', '♟', '♟', '♟', '♟', '♟', '♟'},
			[]rune{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '},
			[]rune{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '},
			[]rune{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '},
			[]rune{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '},
			[]rune{'♙', '♙', '♙', '♙', '♙', '♙', '♙', '♙'},
			[]rune{'♖', '♘', '♗', '♕', '♔', '♗', '♘', '♖'},
		},
	}

	if err := app.gameloop(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
