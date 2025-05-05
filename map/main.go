package main

import (
	"fmt"
	"strconv"
)

func main() {

	fmt.Println(hashCode("Шsome"))
	fmt.Println(hashCode("13"))
	fmt.Println(hashCode("soma"))
	fmt.Println(hashCode("sone"))
	fmt.Println(hashCode("sonea"))
	fmt.Println(hashCode("soneaa"))
	fmt.Println(hashCode("cat"))
	fmt.Println(hashCode("dog"))
	fmt.Println(hashCode("sonesoneqwe"))

	fmt.Println(hashCode("zzzz"))
	fmt.Println(hashCode("zyzy"))
	fmt.Println(hashCode("WZFA"))
	fmt.Println(hashCode("zJJo"))
	fmt.Println(hashCode("ZZW9"))
	fmt.Println(hashCode("WWRA"))
	fmt.Println(hashCode("PASR"))
	fmt.Println(hashCode("MARK"))
	fmt.Println(hashCode("mark"))
	fmt.Println(hashCode("ROMA"))
	fmt.Println(hashCode("roma"))
	fmt.Println(hashCode("abcd"))

	m := makeMap()
	set(m, "aa", "value0")
	set(m, "mark", "value1")
	set(m, "mard", "value2")
	set(m, "roma", "value3")
	set(m, "abcd", "value4")
	set(m, "ZZW9", "value5")
	set(m, "ZZW9Z32FF", "value6")
	set(m, "ZZW9Z32FFFZFEEADE123", "value7")

	fmt.Println(m)

	fmt.Println(get(m, "aa"))
	fmt.Println(get(m, "mark"))
	fmt.Println(get(m, "mard"))
	fmt.Println(get(m, "roma"))
	fmt.Println(get(m, "abcd"))
	fmt.Println(get(m, "ZZW9"))
	fmt.Println(get(m, "ZZW9Z32FF"))
	fmt.Println(get(m, "ZZW9Z32FFFZFEEADE123"))
}

func hashCode(s string) int {
	total := 0
	for i, r := range s {
		if (i+1)%2 == 0 {
			total *= int(r % 33)
		} else {
			total += int(r % 33)
		}

	}
	return total
}

func makeMap() []any {
	slc := make([]any, 20) // produce 1-layered map
	return slc
}

func set(m []any, key string, val string) {
	hash := hashCode(key)
	hashStr := strconv.Itoa(hash)
	setr(m, hashStr, val)
}

func setr(l []any, hashpart string, val string) {
	// choose strategy 'set' or walk deeper
	if len(hashpart) <= 1 { // 'set'
		layer, _ := strconv.Atoi(hashpart)
		l[layer*2] = val
		return
	} else { // 'walk' depper
		layerStr, restPart := hashpart[:1], hashpart[1:]
		layer, _ := strconv.Atoi(layerStr)
		if l[layer*2+1] == nil {
			l[layer*2+1] = makeMap() // make layer
		}
		setr(l[layer*2+1].([]any), restPart, val)
	}
}

func get(m []any, key string) string {
	hash := hashCode(key)
	hashStr := strconv.Itoa(hash)

	return getr(m, hashStr)
}

func getr(l []any, hashpart string) string {
	// choose strategy 'get' or walk deeper
	if len(hashpart) <= 1 { // 'get'
		layer, _ := strconv.Atoi(hashpart)
		return l[layer*2].(string)
	} else { // 'walk' depper
		layer, restPart := hashpart[:1], hashpart[1:]
		layerInt, _ := strconv.Atoi(layer)
		return getr(l[layerInt*2+1].([]any), restPart)
	}
}
