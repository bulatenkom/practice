package main

import "fmt"

func main() {
	arr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	mapped := Map(arr, func(it int) int { return it + 5 })
	filtered := filter(arr, func(it int) bool { return it%2 == 0 })
	reduced := reduce(arr, func(acc, iter int) int { return acc + iter }, 0)

	fmt.Println(arr)
	fmt.Println(mapped)
	fmt.Println(filtered)
	fmt.Println(reduced)
}

// в go зарезервировано ключевое слово 'map', поэтому используем CamelCase нотацию
func Map(arr []int, mapFn func(int) int) []int {
	r := make([]int, len(arr))

	for i := range arr {
		r[i] = mapFn(arr[i])
	}
	return r
}

func filter(arr []int, predicateFn ...func(int) bool) []int {
	r := []int{}
outer:
	for _, v := range arr {
		for _, predicate := range predicateFn {
			if !predicate(v) {
				continue outer
			}
		}
		r = append(r, v)
	}
	return r
}

func reduce(arr []int, reducerFn func(acc, iter int) int, base int) int {
	r := base
	for _, v := range arr {
		r = reducerFn(r, v)
	}
	return r
}
