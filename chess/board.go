package main

import (
	"slices"
	"strings"
)

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
		for move := range s {
			if newboard := b.move(move); newboard.isCheck(v) {
				delete(s, move)
			}
		}
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

		if strings.ContainsRune(string(whitePieces), valueAtShift(0, 0)) && strings.ContainsRune(string(blackPieces), valueAtShift(di, dj)) {
			return true
		}
		if strings.ContainsRune(string(blackPieces), valueAtShift(0, 0)) && strings.ContainsRune(string(whitePieces), valueAtShift(di, dj)) {
			return true
		}
		return false
		// return valueAtShift(0, 0) == '♙' && valueAtShift(di, dj) == '♟' || valueAtShift(0, 0) == '♟' && valueAtShift(di, dj) == '♙'
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

func (b *board) attackers(cell string) set {
	s := Set([]string{})

	appendSuitable := func(moves set, pieces ...rune) {
		cells := []string{}
		for mv := range moves {
			cells = append(cells, mv[2:])
		}

		for _, c := range cells {
			if v, _, _ := b.valueAt(c); slices.Contains(pieces, v) {
				Append(s, c)
			}
		}
	}

	appendSuitable(b.movementsPawn(cell), '♙', '♟')
	appendSuitable(b.movementsKnight(cell), '♘', '♞')
	appendSuitable(b.movementsOrthogonal(cell, 7), '♖', '♜', '♕', '♛')
	appendSuitable(b.movementsOrthogonal(cell, 1), '♔', '♚')
	appendSuitable(b.movementsDiagonals(cell, 7), '♗', '♝', '♕', '♛')
	appendSuitable(b.movementsDiagonals(cell, 1), '♔', '♚')
	appendSuitable(b.movementsOrthogonal(cell, 1), '♔', '♚')

	return s
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

func (b *board) copy() board {
	cp := board(make([][]rune, len(*b)))
	for i, v := range *b {
		cp[i] = make([]rune, len((*b)[i]))
		copy(cp[i], v)
	}
	return cp
}

func (b *board) move(move string) board {
	cp := b.copy()
	v, i, j := b.valueAt(move[:2])
	cp[i][j] = ' '
	_, i, j = b.valueAt(move[2:])
	cp[i][j] = v
	return cp
}

// accepts white '♔' and black '♚'
func (b *board) isCheck(king rune) bool {
	// find given king
	for i, line := range *b {
		for j, v := range line {
			if v == king {
				return len(b.attackers(b.cellAt(i, j))) > 0
			}
		}
	}
	panic("King is not present on board")
}

// accepts white '♔' and black '♚'
func (b *board) isCheckmate(king rune) bool {
	if !b.isCheck(king) {
		return false
	}

	var allyPieces []rune
	if king == '♔' {
		allyPieces = whitePieces
	} else {
		allyPieces = blackPieces
	}
	// spot save moves
	allMoves := Set([]string{})
	for i, line := range *b {
		for j, v := range line {
			if strings.ContainsRune(string(allyPieces), v) {
				allMoves = Union(allMoves, b.movements(b.cellAt(i, j)))
			}
		}
	}
	for move := range allMoves {
		newboard := b.move(move)
		if !newboard.isCheck(king) {
			return false
		}
	}
	return true
}

// accepts white '♔' and black '♚'
func (b *board) isStalemate(king rune) bool {
	if b.isCheck(king) {
		return false
	}

	var allyPieces []rune
	if king == '♔' {
		allyPieces = whitePieces
	} else {
		allyPieces = blackPieces
	}

	allMoves := Set([]string{})
	for i, line := range *b {
		for j, v := range line {
			if strings.ContainsRune(string(allyPieces), v) {
				allMoves = Union(allMoves, b.movements(b.cellAt(i, j)))
			}
		}
	}
	return len(allMoves) == 0
}
