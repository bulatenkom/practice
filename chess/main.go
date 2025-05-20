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
			if input[:2] == input[2:] {
				app.ui.err = fmt.Errorf("move '%s' is not possible. Piece cannot move to the same cell", input)
				continue
			}
			if v, _, _ := app.board.valueAt(input[:2]); v == ' ' {
				app.ui.err = fmt.Errorf("move '%s' is not possible. Cell is empty, nothing to move", input)
				continue
			}
			if v, _, _ := app.board.valueAt(input[:2]); v == 'Z' { // TODO
				app.ui.err = fmt.Errorf("move '%s' is not possible. You cannot move enemy piece", input)
				continue
			}
			if !Has(app.board.movements(input[:2]), input) {
				app.ui.err = fmt.Errorf("move '%s' is not possible. Such move is not possible for given piece", input)
				continue
			}

			// TODO handle special case for Pawn transformation
			// TODO handle special case for King and Rook

			// TODO apply move and check that ally King is not under attack

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
		if !avMovesRegex.MatchString(input) {
			app.ui.movements = nil
		}
	}
}

func helpTemplate() string {
	return `HELP
  help, h, ?, /? - toggles help window
  [a-h][1-8] - prints possible moves for selected piece (e.g. 'e2' prints moves for Pawn)
  [a-h][1-8][a-h][1-8] - move piece (e.g. 'e2e4' moves piece from e2 to e4)
  hist, history - prints performed moves during the game
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
			NewBox(RenderBoard(app.board, app.ui.movements)),
			NewBox("move count: " + strconv.Itoa(app.moveCount)),
		},
		Vertical,
	)
	sideview := Layout(
		[]Box{
			NewBox(app.ui.sideview),
			// NewBox(ansiRed(app.ui.sideview)), // FIXME
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

var (
	blackPieces = []rune{'♟', '♞', '♝', '♜', '♛', '♚'}
	whitePieces = []rune{'♙', '♘', '♗', '♖', '♕', '♔'}
)

// get current player and its pieces set
func (app *app) currentPlayer() string {
	if app.moveCount%2 == 0 {
		return "W"
	} else {
		return "B"
	}
}

type uiInfo struct {
	banner    string // элемент для вывода текстовой информации хода игры
	sideview  string // элемент для вывода вспомогательной текстовой информации
	movements set    // элемент содержащий возможные ходы для выбранной фигуры (используется для раскраски целевых клеток)
	history   string // история ходов
	err       error  // элемент для вывода ошибок пользовательского ввода
}

type board [][]rune

// movements() of piece in selected cell
func (b *board) movements(cell string) set {
	s := Set([]string{})

	v, _, _ := b.valueAt(cell)

	switch v {
	case '♙', '♟':
		s = b.movementsPawn(cell)
	case '♘', '♞':
		s = b.movementsKnight(cell)
	case '♗', '♝':
		s = b.movementsDiagonals(cell, 7)
	case '♖', '♜':
		s = b.movementsOrthogonal(cell, 7)
	case '♕', '♛':
		s = Union(
			b.movementsOrthogonal(cell, 7),
			b.movementsDiagonals(cell, 7),
		)
	case '♔', '♚':
		s = Union(
			b.movementsOrthogonal(cell, 1),
			b.movementsDiagonals(cell, 1),
		)
		// TODO specific cases
	}
	return s
}

// util function that traces possible moves by diagonals limited by step (cell must be ally!!!)
func (b *board) movementsDiagonals(cell string, step int) set {
	s := Set([]string{})

	cv, _, _ := b.valueAt(cell)

	if cv == ' ' {
		return s
	}

	trace := b.makeTraceFn(s, cell)

	// directions
	topRight, bottomRight, bottomLeft, topLeft := true, true, true, true

	for i := 1; i <= step; i++ {
		if topRight {
			topRight = trace(-i, i)
		}
		if bottomRight {
			bottomRight = trace(i, i)
		}
		if bottomLeft {
			bottomLeft = trace(i, -i)
		}
		if topLeft {
			topLeft = trace(-i, -i)
		}
	}
	return s
}

// util function that traces possible moves by orthogonal lines limited by step (cell must be ally!!!)
func (b *board) movementsOrthogonal(cell string, step int) set {
	s := Set([]string{})

	cv, _, _ := b.valueAt(cell)

	if cv == ' ' {
		return s
	}

	trace := b.makeTraceFn(s, cell)

	// directions
	top, right, bottom, left := true, true, true, true

	for i := 1; i <= step; i++ {
		if top {
			top = trace(-i, 0)
		}
		if right {
			right = trace(0, i)
		}
		if bottom {
			bottom = trace(i, 0)
		}
		if left {
			left = trace(0, -i)
		}
	}
	return s
}

func (b *board) movementsPawn(cell string) set {
	s := Set([]string{})

	cv, ci, cj := b.valueAt(cell)

	if cv == ' ' {
		return s
	}

	trace := b.makeTraceFn(s, cell)

	isBoardBounds := func(di, dj int) bool {
		return (ci+di) >= 0 && (ci+di) < 8 && (cj+dj) >= 0 && (cj+dj) < 8
	}

	valueAtShift := func(di, dj int) rune {
		return (*b)[ci+di][cj+dj]
	}

	isEnemy := func(di, dj int) bool {
		return valueAtShift(0, 0) == '♙' && valueAtShift(di, dj) == '♟' || valueAtShift(0, 0) == '♟' && valueAtShift(di, dj) == '♙'
	}

	// 1	- calculate for black Pawn
	// -1	- calculate for white Pawn
	movements := func(direction int) set {
		if isBoardBounds(direction*1, 0) && valueAtShift(direction*1, 0) == ' ' {
			trace(direction*1, 0)
			if isBoardBounds(direction*2, 0) && valueAtShift(direction*2, 0) == ' ' && (ci == 1 || ci == 6) {
				trace(direction*2, 0)
			}
		}
		if isBoardBounds(direction*1, 1) && isEnemy(direction*1, 1) {
			trace(direction*1, 1)
		}
		if isBoardBounds(direction*1, -1) && isEnemy(direction*1, -1) {
			trace(direction*1, -1)
		}
		return s
	}

	if strings.ContainsRune(string(whitePieces), cv) {
		return movements(-1)
	} else {
		return movements(1)
	}
}

func (b *board) movementsKnight(cell string) set {
	s := Set([]string{})

	cv, _, _ := b.valueAt(cell)

	if cv == ' ' {
		return s
	}

	trace := b.makeTraceFn(s, cell)

	trace(2, 1)
	trace(2, -1)
	trace(-2, 1)
	trace(-2, -1)
	trace(1, 2)
	trace(-1, 2)
	trace(1, -2)
	trace(-1, -2)

	return s
}

// cell must contain piece
// piece placed within cell is assumed to be allied
func (b *board) makeTraceFn(s set, cell string) func(int, int) bool {
	cv, ci, cj := b.valueAt(cell)

	var allyPieces, enemyPieces []rune
	if strings.ContainsRune(string(whitePieces), cv) {
		allyPieces = whitePieces
		enemyPieces = blackPieces
	} else {
		allyPieces = blackPieces
		enemyPieces = whitePieces
	}
	unused(enemyPieces)

	return func(di, dj int) bool {
		if (ci+di) >= 0 && (ci+di) < 8 && (cj+dj) >= 0 && (cj+dj) < 8 {
			v := (*b)[ci+di][cj+dj]
			if v == ' ' {
				Append(s, cell+b.cellAt(ci+di, cj+dj))
				return true
			} else if strings.ContainsRune(string(allyPieces), v) {
				return false
			} else {
				Append(s, cell+b.cellAt(ci+di, cj+dj))
				return false
			}
		}
		return false
	}
}

func unused(i ...any) {}

// func (b *board) attackers(cell string) set {
// 	s := Set([]string{})

// 	v, i, j := b.valueAt(cell)

// 	return s
// }

func movementsFormatted(mv map[string]struct{}) string {
	pairs := []string{}
	for k := range mv {
		pairs = append(pairs, fmt.Sprintf("(%s)", k))
	}
	return strings.Join(pairs, " ")
}

func (b *board) valueAt(cell string) (rune, int, int) {
	// 49 50 .. 56 - asci code
	//  1  2 ..  8 - chess notation
	//  7  6 ..  0 - board index
	i := -(int(cell[1]) - 56)

	j := int(cell[0]) - 97
	return (*b)[i][j], i, j
}

func (b *board) cellAt(i, j int) string {
	// backward operation to valueAt(...)
	return string(rune(j)+97) + string(rune(-i)+56)
}

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
			[]rune{' ', '♟', ' ', ' ', ' ', ' ', ' ', ' '},
			[]rune{'♟', ' ', ' ', '♟', ' ', '♟', ' ', ' '},
			// []rune{' ', ' ', ' ', '♟', ' ', '♟', ' ', ' '},
			[]rune{'♙', '♙', '♙', '♙', '♙', '♙', '♙', '♙'},
			[]rune{'♖', '♘', '♗', '♕', '♔', '♗', '♘', '♖'},
		},
	}

	if err := app.gameloop(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
