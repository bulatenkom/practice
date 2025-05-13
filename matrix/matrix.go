package main

import (
	"fmt"
	"strconv"
)

type Matrix [][]float64

func NewMatrix(rowcount, colcount int, fill float64) *Matrix {
	if rowcount < 1 || colcount < 1 {
		panic("Matrix rowcount and colcount must be equal or more than 1")
	}

	m := make([][]float64, rowcount)
	for i := range m {
		m[i] = make([]float64, colcount)
		for j := range m[i] {
			m[i][j] = fill
		}
	}
	return (*Matrix)(&m)
}

func MatrixOf(rows ...[]float64) *Matrix {
	m := NewMatrix(len(rows), len(rows[0]), 0)

	for i := range rows {
		for j := range rows[i] {
			(*m)[i][j] = rows[i][j]
		}
	}
	return m
}

func NewVector(size int, fill float64) *Matrix {
	return NewMatrix(size, 1, fill)
}

func VectorOf(items ...float64) *Matrix {
	return MatrixOf(items)
}

type closure func(value float64)

func (c closure) set(value float64) {
	c(value)
}

func (m *Matrix) at(i, j int) closure {
	return func(value float64) {
		(*m)[i][j] = value
	}
}

func (m *Matrix) fill(v float64) {
	for i := range *m {
		for j := range (*m)[i] {
			(*m)[i][j] = v
		}
	}
}

func (m *Matrix) print() {
	output := ""
	for i := range *m {
		for j := range (*m)[0] {
			output += strconv.FormatFloat((*m)[i][j], 'f', -1, 64) + " "
		}
		output += "\n"
	}
	fmt.Println(output)
}

func (m *Matrix) rowCount() int {
	return len((*m))
}

func (m *Matrix) colCount() int {
	return len((*m)[0])
}

func (m *Matrix) dimensions() string {
	return fmt.Sprint(m.rowCount(), "×", m.colCount())
}

func (m *Matrix) colAt(j int) []*float64 {
	res := make([]*float64, m.rowCount())
	for i := range *m {
		res[i] = &(*m)[i][j]
	}
	return res
}

func (m *Matrix) colAtImmut(j int) []float64 {
	res := make([]float64, m.rowCount())
	for i := range *m {
		res[i] = (*m)[i][j]
	}
	return res
}

func (m *Matrix) mul(k float64) *Matrix {
	res := NewMatrix(m.rowCount(), m.colCount(), 0)
	for i := range *m {
		for j := range (*m)[0] {
			(*res)[i][j] = (*m)[i][j] * k
		}
	}
	return res
}

func (m *Matrix) hasEqualDimensions(other *Matrix) bool {
	return m.dimensions() == other.dimensions()
}

func (m *Matrix) equals(other *Matrix) bool {
	if !m.hasEqualDimensions(other) {
		return false
	}
	for i := range *m {
		for j := range (*m)[0] {
			if (*m)[i][j] != (*other)[i][j] {
				return false
			}
		}
	}
	return true
}

func (m *Matrix) add(other *Matrix) *Matrix {
	m.validateDimensionsEquality(other)

	res := NewMatrix(m.rowCount(), m.colCount(), 0)
	for i := range *m {
		for j := range (*m)[0] {
			(*res)[i][j] = (*m)[i][j] + (*other)[i][j]
		}
	}
	return res
}

func (m *Matrix) subtract(other *Matrix) *Matrix {
	m.validateDimensionsEquality(other)

	res := NewMatrix(m.rowCount(), m.colCount(), 0)
	for i := range *m {
		for j := range (*m)[0] {
			(*res)[i][j] = (*m)[i][j] - (*other)[i][j]
		}
	}
	return res
}

func (m *Matrix) validateDimensionsEquality(other *Matrix) {
	if !m.hasEqualDimensions(other) {
		panic("Matrices of different dimensions cannot be added")
	}
}

func (m *Matrix) __mulRowByColumn(row, col []float64) float64 {
	if len(row) != len(col) {
		panic(fmt.Sprintf("Cannot multiply row %v by column len=%v of different sizes.", len(row), len(col)))
	}

	res := 0.0
	for i := range row {
		res += row[i] * col[i]
	}
	return res
}

func (m *Matrix) mulm(other *Matrix) *Matrix {
	m.validateMulm(other)

	res := NewMatrix(m.rowCount(), other.colCount(), 0)
	for i := range *m {
		for j := range (*other)[0] {
			(*res)[i][j] = m.__mulRowByColumn((*m)[i], other.colAtImmut(j))
		}
	}
	return res
}

func (m *Matrix) validateMulm(other *Matrix) {
	if m.colCount() != other.rowCount() {
		panic(fmt.Sprintf("Incompatible dimensions. Columns-dimension of left matrix %v should be equal to rows-dimensions of second one %v", m.dimensions(), other.dimensions()))
	}
}
