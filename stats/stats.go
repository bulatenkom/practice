package main

import "math"

// slice 's' must have only whole numbers from 0 to 99
func mode(s []int) int {
	freqTable := make([]int, 100)
	for i := 0; i < len(s); i++ {
		freqTable[s[i]]++
	}
	max := s[0]
	memi := 0
	for i := 1; i < len(freqTable); i++ {
		if freqTable[i] > max {
			max = freqTable[i]
			memi = i
		}
	}
	return memi
}

func mean(s []int) float64 {
	total := 0
	for i := 0; i < len(s); i++ {
		total += s[i]
	}
	return float64(total) / float64(len(s))
}

// median calculated approx. (can be slitghtly less than actual)
func median(s []int) int {
	// sort
	for k := 0; k < len(s); k++ {
		for i := 0; i < len(s)-1; i++ {
			if s[i] > s[i+1] {
				s[i], s[i+1] = s[i+1], s[i]
			}
		}
	}

	return s[(len(s)+1)/2]
}

func std(s []int) float64 {
	squareDiff := 0.0
	mean := mean(s)
	for i := 0; i < len(s); i++ {
		squareDiff += (float64(s[i]) - mean) * (float64(s[i]) - mean)
	}
	return math.Sqrt(squareDiff / float64(len(s)))
}

func variance(s []int) float64 {
	std := std(s)
	return std * std
}
