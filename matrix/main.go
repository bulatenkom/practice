package main

import "fmt"

func main() {
	fmt.Println("Test: Matrix initialization and print")
	for _, m := range []*Matrix{
		NewMatrix(3, 3, 0),
		NewMatrix(1, 5, 7),
		NewMatrix(2, 4, 2),
		NewMatrix(4, 1, 9),
	} {
		m.print()
	}

	fmt.Println("\nTest: Vector creation, dimensions, and scalar multiplication")
	v := NewVector(4, 2)
	v.print()
	fmt.Println("Dimensions:", v.dimensions())

	v = v.mul(3)
	v.print()

	g := NewVector(4, 2)
	fmt.Println("Are vectors equal?", v.equals(g))

	fmt.Println("\nTest: Matrix addition")
	{
		m1 := NewMatrix(1, 5, 7)
		m2 := NewMatrix(1, 5, 2)
		res := m1.add(m2)
		res.print()
	}

	fmt.Println("\nTest: Set element after addition")
	{
		m1 := NewMatrix(1, 5, 7)
		m2 := NewMatrix(1, 5, 2)
		res := m1.add(m2)
		res.at(0, 0).set(10)
		res.print()
	}

	fmt.Println("\nTest: Matrix multiplication")
	{
		m1 := NewMatrix(2, 3, 2)
		m1.at(0, 1).set(3)
		m1.at(0, 2).set(1)
		m1.at(1, 0).set(1)
		m1.print()

		m2 := NewMatrix(3, 2, 0)
		m2.at(0, 0).set(2.65)
		m2.at(0, 1).set(2.25)
		m2.at(1, 0).set(1.55)
		m2.at(1, 1).set(1.5)
		m2.at(2, 0).set(2.35)
		m2.at(2, 1).set(2.2)
		m2.print()

		res := m1.mulm(m2)
		res.print()
	}

	fmt.Println("\nTest: Vector × Matrix multiplication")
	{
		v := VectorOf(1, 2, 3)
		m := NewMatrix(3, 1, 1)

		v.print()
		m.print()

		res := v.mulm(m)
		res.print()
	}

	fmt.Println("\nTest: Multiplication of arbitrary matrices")
	{
		m1 := MatrixOf([]float64{1, 0}, []float64{1, 2})
		m2 := MatrixOf([]float64{-1, 1}, []float64{0, 3})

		m1.print()
		m2.print()

		m1.mulm(m2).print()
		m2.mulm(m1).print()
	}

	fmt.Println("\nTest: Multiplication with zero/identity-like matrices")
	{
		m1 := MatrixOf([]float64{1, 0}, []float64{0, 0})
		m2 := MatrixOf([]float64{0, 0}, []float64{0, 1})

		m1.print()
		m2.print()

		m1.mulm(m2).print()
		m2.mulm(m1).print()
	}

	fmt.Println("\nTest: Matrix × Identity and × Ones")
	{
		m1 := MatrixOf([]float64{1, 2}, []float64{3, 4})
		m2 := NewMatrix(2, 2, 1)
		m3 := MatrixOf([]float64{1, 0}, []float64{0, 1})

		m1.print()

		m1.mulm(m2).print()
		m1.mulm(m3).print()
		m3.mulm(m1).print()
	}

	fmt.Println("\nTest: Matrix exponentiation")
	{
		m1 := MatrixOf([]float64{1, 2}, []float64{3, 4})

		m1.mulm(m1).print()                   // m1^2
		m1.mulm(m1).mulm(m1).print()          // m1^3
		m1.mulm(m1).mulm(m1).mulm(m1).print() // m1^4
	}
}
