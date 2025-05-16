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
	if _, err := fmt.Println(ansiYellow(app.ui.banner)); err != nil {
		return err
	}
	if err := renderBoard(app.board); err != nil {
		return err
	}
	if _, err := fmt.Println("move count:", app.moveCount); err != nil {
		return err
	}
	if app.ui.err != nil {
		if _, err := fmt.Println(ansiRed(app.ui.err.Error())); err != nil {
			return err
		}
	}
	return nil
}

func clearScreen() (int, error) {
	return fmt.Print("\x1b[H\x1b[2J")
}

// ANSI escape codes
const (
	c_reset = "\033[0m"
	// 3-4 bit
	c_red    = "\033[31m"
	c_yellow = "\033[33m"
	c_bold   = "\033[1m"
	// 24 bit TrueColor
	fg_black       = "\033[38;2;0;0;0m"
	bg_med_brown   = "\033[48;2;209;139;71m"
	bg_light_peach = "\033[48;2;255;206;158m"
)

func fgTrueColor(r, g, b int) string {
	return fmt.Sprintf("\033[38;2;%v;%v;%vm", r, g, b)
}

func bgTrueColor(r, g, b int) string {
	return fmt.Sprintf("\033[48;2;%v;%v;%vm", r, g, b)
}

func ansiYellow(s string) string {
	return c_yellow + s + c_reset
}

func ansiRed(s string) string {
	return c_red + s + c_reset
}

func renderBoard(board board) error {
	vMarkers := []rune{'8', '7', '6', '5', '4', '3', '2', '1'}
	hMarkers := []rune{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h'}

	output := joinHMarkers(hMarkers)

	for i := range board {
		output += string(vMarkers[i]) + " "
		for j := range board[0] {
			cell := " " + string(board[i][j]) + " "
			if (i+j)%2 == 0 {
				output += bg_light_peach + fg_black + cell + c_reset
			} else {
				output += bg_med_brown + fg_black + cell + c_reset
			}
		}
		output += " " + string(vMarkers[i]) + "\n"
	}
	output += joinHMarkers(hMarkers)
	if _, err := fmt.Print(output); err != nil {
		return err
	}
	return nil
}

func joinHMarkers(markers []rune) string {
	output := "  "
	for _, v := range markers {
		output += " " + string(v) + " "
	}
	output += "\n"
	return output
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
