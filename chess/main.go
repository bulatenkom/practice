package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
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
	if _, err := clearScreen(); err != nil {
		return err
	}
	mainview := Layout(
		[]Box{
			NewBox(ansiYellow(app.ui.banner)),
			NewBox(renderBoard(app.board)),
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
		if _, err := fmt.Println(ansiRed(app.ui.err.Error())); err != nil {
			return err
		}
	}
	return nil
}

func clearScreen() (int, error) {
	return fmt.Print("\x1b[H\x1b[2J")
}

// Box may contain single(multi)line content with all lines expanded to the width of the longest line
type Box struct {
	// width, height
	w, h    int
	content []string
}

func (b *Box) String() string {
	return strings.Join(b.content, "")
}

func NewBox(s string) Box {
	lines := strings.Split(s, "\n")
	if lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	width := 0
	for _, line := range lines {
		width = max(width, CountRunesIgnoringAnsiEscapes(line))
	}
	content := []string{}
	for _, line := range lines {
		// buf := make([]rune, width+1) // reserve place for '\n'
		// copy(buf, []rune(line))
		// for i := utf8.RuneCountInString(line); i < width; i++ {
		// 	buf[i] = '\x00'
		// }
		// buf[width] = '\n'
		// content = append(content, string(buf))

		fline := line + strings.Repeat(" ", width-CountRunesIgnoringAnsiEscapes(line)) + "\n"
		content = append(content, fline)
	}

	return Box{
		w:       width,
		h:       len(lines),
		content: content,
	}
}

var escapeSeqMatcher *regexp.Regexp = regexp.MustCompile("(\033\\[0m|\033\\[31m|\033\\[33m|\033\\[1m|\033\\[38;2;0;0;0m|\033\\[48;2;209;139;71m|\033\\[48;2;255;206;158m)")

func CountRunesIgnoringAnsiEscapes(s string) int {
	res := escapeSeqMatcher.ReplaceAllLiteralString(s, "")
	return utf8.RuneCountInString(res)
}

type Orientation byte

const (
	Horizontal Orientation = iota
	Vertical
)

func Layout(boxes []Box, orientation Orientation) Box {
	if orientation == Vertical {
		content := ""
		for _, b := range boxes {
			content += b.String()
		}
		return NewBox(content)
	}
	if orientation == Horizontal {
		w := 0
		h := 0

		for _, b := range boxes {
			w += b.w
			h = max(h, b.h)
		}

		content := ""
		for r := 0; r < h; r++ {
			line := ""
			for _, b := range boxes {
				if r+1 > b.h {
					line += strings.Repeat(" ", b.w)
					line += "|"
					continue
				}
				s := b.content[r]
				line += s[:len(s)-1]
				line += "|"
			}
			line += "\n"
			content += line
		}

		return NewBox(content)
	}
	return Box{}
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
	history  string // история ходов
	err      error  // элемент для вывода ошибок пользовательского ввода
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
