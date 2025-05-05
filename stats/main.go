package main

import "fmt"

func main() {
	series := []int{2, 3, 1, 1, 0, 0, 4, 3, 0, 1, 2, 3, 2, 1, 4}
	fmt.Println("series", series)
	fmt.Println("statistics:")
	fmt.Println("mode", mode(series))
	fmt.Println("median", median(series))
	fmt.Println("mean", mean(series))
	fmt.Println("std", std(series))
	fmt.Println("var", variance(series))
}
