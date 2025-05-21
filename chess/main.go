package main

import (
	"fmt"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

type app struct {
	board     board // игровое поле с фигурами
	moveCount int   // счётчик ходов
	flags     struct {
		check     bool // шах
		checkmate bool // мат
		stalemate bool // пат
	}
	// текущее время
	// общее время партии
	ui uiInfo // элементы интерфейса
}

type uiInfo struct {
	banner    string // элемент для вывода текстовой информации хода игры
	sideview  string // элемент для вывода вспомогательной текстовой информации
	movements set    // элемент содержащий возможные ходы для выбранной фигуры (используется для раскраски целевых клеток)
	history   string // история ходов
	err       error  // элемент для вывода ошибок пользовательского ввода
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
			app.ui.movements = app.board.movements(input)
			app.ui.sideview = movementsFormatted(app.ui.movements)
		case moveRegex.MatchString(input):
			// handle move
			if app.flags.checkmate || app.flags.stalemate {
				app.ui.err = fmt.Errorf("move '%s' is not possible. Game has finished, restart a program to start a new one", input)
				continue
			}
			if input[:2] == input[2:] {
				app.ui.err = fmt.Errorf("move '%s' is not possible. Piece cannot move to the same cell", input)
				continue
			}
			if v, _, _ := app.board.valueAt(input[:2]); v == ' ' {
				app.ui.err = fmt.Errorf("move '%s' is not possible. Cell is empty, nothing to move", input)
				continue
			}
			var enemyPieces []rune
			if app.currentPlayer() == "W" {
				enemyPieces = blackPieces
			} else {
				enemyPieces = whitePieces
			}
			if v, _, _ := app.board.valueAt(input[:2]); slices.Contains(enemyPieces, v) {
				app.ui.err = fmt.Errorf("move '%s' is not possible. You cannot move enemy piece", input)
				continue
			}
			if !Has(app.board.movements(input[:2]), input) {
				app.ui.err = fmt.Errorf("move '%s' is not possible. Such move is not possible for given piece", input)
				continue
			}

			var king rune
			if app.currentPlayer() == "W" {
				king = '♔'
			} else {
				king = '♚'
			}
			newboard := app.board.move(input)
			if newboard.isCheck(king) {
				app.ui.err = fmt.Errorf("move '%s' is not possible. Your '%s' is under check", input, string(king))
				continue
			}
			app.flags.check = false
			// TODO handle special case for Pawn transformation
			// TODO handle special case for King and Rook
			app.ui.banner = ""
			app.board = newboard
			app.updateStateFlags()
			app.moveCount++
			app.ui.history += fmt.Sprintf(" %v. %v \n", app.moveCount, input)
			app.ui.sideview = app.ui.history
		default:
			app.ui.err = fmt.Errorf("command '%s' is unknown. Type 'h' to get full list of commands", input)
			continue
		}
		// reset temp states
		app.ui.err = nil
		if !avMovesRegex.MatchString(input) {
			app.ui.movements = nil
		}
	}
}

func (app *app) updateStateFlags() {
	var enemyKing rune
	if app.currentPlayer() == "W" {
		enemyKing = '♚'
	} else {
		enemyKing = '♔'
	}

	if app.board.isCheck(enemyKing) {
		app.flags.check = true
	}
	if app.board.isCheckmate(enemyKing) {
		app.flags.checkmate = true
	}
	if app.board.isStalemate(enemyKing) {
		app.flags.stalemate = true
	}
}

// paintUI()
func (app *app) paintUI() error {
	if _, err := ClearScreen(); err != nil {
		return err
	}
	mainview := Layout(
		[]Box{
			NewBox(renderBanner(app) + " " + AnsiYellow(app.currentPlayerText())),
			NewBox(RenderBoard(app.board, app.ui.movements)),
			NewBox("move count: " + strconv.Itoa(app.moveCount)),
		},
		Vertical,
	)
	sideview := Layout(
		[]Box{
			NewBox(app.ui.sideview),
		},
		Vertical,
	)
	layout := Layout([]Box{mainview, sideview}, Horizontal)
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

// get current player and its pieces set
func (app *app) currentPlayer() string {
	if app.moveCount%2 == 0 {
		return "W"
	} else {
		return "B"
	}
}

func (app *app) currentPlayerText() string {
	switch app.currentPlayer() {
	case "W":
		return "White"
	case "B":
		return "Black"
	default:
		return ""
	}
}

func renderBanner(app *app) string {
	if app.flags.checkmate {
		return AnsiRed("checkmate")
	}
	if app.flags.stalemate {
		return AnsiRed("stalemate")
	}
	if app.flags.check {
		return AnsiRed("check")
	}
	return AnsiYellow(app.ui.banner)
}

func helpTemplate() string {
	return `HELP
  help, h, ?, /? - toggles help window
  [a-h][1-8] - prints possible moves for selected piece (e.g. 'e2' prints moves for Pawn)
  [a-h][1-8][a-h][1-8] - move piece (e.g. 'e2e4' moves piece from e2 to e4)
  hist, history - prints performed moves during the game
`
}

func movementsFormatted(mv map[string]struct{}) string {
	pairs := []string{}
	wc := 0
	for k := range mv {
		pairs = append(pairs, fmt.Sprintf("(%s)", k))
		wc++
		if wc%6 == 0 {
			pairs = append(pairs, "\n")
		}
	}
	return strings.Join(pairs, " ")
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
