package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
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

// paintUI()
func (app *app) paintUI() error {
	if _, err := clearScreen(); err != nil {
		return err
	}
	mainview := renderView(
		ansiYellow(app.ui.banner),
		renderBoard(app.board),
		"move count: "+strconv.Itoa(app.moveCount),
	)
	layout := mainview.content
	if _, err := fmt.Println(layout); err != nil {
		return err
	}
	if app.ui.err != nil {
		if _, err := fmt.Println(ansiRed(app.ui.err.Error())); err != nil {
			return err
		}
	}
	return nil
}

func renderView(components ...string) Box {
	mainview := ""
	for _, v := range components {
		mainview += v
		if v[len(v)-1] != '\n' {
			mainview += "\n"
		}
	}
	return NewBox(mainview)
}

func clearScreen() (int, error) {
	return fmt.Print("\x1b[H\x1b[2J")
}

// Box may contain single(multi)line content with all lines expanded to the width of the longest line
type Box struct {
	// width, height
	w, h    int
	content string
}

func NewBox(s string) Box {
	lines := strings.Split(s, "\n")
	if lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	width := 0
	for _, line := range lines {
		width = max(width, len(line))
	}
	content := ""
	for _, line := range lines {
		buf := make([]rune, width+1)
		copy(buf, []rune(line))
		for i := len(line); i < width; i++ {
			buf[i] = '-'
		}
		buf[width] = '\n'

		content += string(buf)
	}

	return Box{
		w:       width,
		h:       len(lines),
		content: content,
	}
}

func renderBoard(board board) string {
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
	return output
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
	mainview string // элемент для вывода игрового интерфейса
	sideview string // элемент для вывода вспомогательной текстовой информации
	// история ходов
	err error // элемент для вывода ошибок пользовательского ввода
}

func main() {
	// in := "abc\nxyzxyz\nsome\n"
	// box := NewBox(in)
	// fmt.Println(in)
	// fmt.Println(box)

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
