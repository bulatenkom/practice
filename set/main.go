package main

import "fmt"

func main() {
	s := set([]int{1, 2, 2, 2, 3, 4, 5, 1, 0, 0, 1, 4}) // s = [0,1,2,3,4,5]
	fmt.Println(s)
	fmt.Println(get(s, 5))
	has(s, 0)    // true
	append(s, 5) // s = [0,1,2,3,4,5]
	fmt.Println(s)
	append(s, 6) // s = [0,1,2,3,4,5,6]

	fmt.Println(has(s, 0))
	fmt.Println(s)
	delete(s, 6) // s = [0,1,2,3,4,5]
	fmt.Println(s)
	clear(s) // s = []
	fmt.Println(s)

	s = set([]int{1, 2, 3, 4, 5})
	t := set([]int{3, 4, 5, 6, 7})
	fmt.Println(union(s, t))                // 1,2,3,4,5,6,7
	fmt.Println(intersection(s, t))         // 3,4,5
	fmt.Println(difference(s, t))           // 1,2
	fmt.Println(symmetric_difference(s, t)) // 1,2,6,7
}

func set(s []int) map[int]struct{} {
	m := map[int]struct{}{}
	for _, v := range s {
		m[v] = struct{}{}
	}
	return m
}

func has(set map[int]struct{}, elem int) bool {
	_, ok := set[elem]
	return ok
}

func get(set map[int]struct{}, elem int) (int, bool) {
	if _, ok := set[elem]; ok {
		return elem, true
	}
	return 0, false
}

func append(set map[int]struct{}, elem int) {
	set[elem] = struct{}{}
}

func union(s, t map[int]struct{}) map[int]struct{} {
	r := make(map[int]struct{}, max(len(s), len(t)))

	for elem, _ := range s {
		r[elem] = struct{}{}
	}
	for elem, _ := range t {
		r[elem] = struct{}{}
	}

	return r
}

func intersection(s, t map[int]struct{}) map[int]struct{} {
	r := make(map[int]struct{})

	for elem, _ := range s {
		if has(t, elem) {
			r[elem] = struct{}{}
		}
	}
	for elem, _ := range t {
		if has(s, elem) {
			r[elem] = struct{}{}
		}
	}
	return r
}

func difference(s, t map[int]struct{}) map[int]struct{} {
	r := make(map[int]struct{})

	for elem, _ := range s {
		if !has(t, elem) {
			r[elem] = struct{}{}
		}
	}
	return r
}

func symmetric_difference(s, t map[int]struct{}) map[int]struct{} {
	return union(difference(s, t), difference(t, s))
}
