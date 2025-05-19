package main

type set = map[string]struct{}

func Set(s []string) map[string]struct{} {
	m := map[string]struct{}{}
	for _, v := range s {
		m[v] = struct{}{}
	}
	return m
}

func Has(set map[string]struct{}, elem string) bool {
	_, ok := set[elem]
	return ok
}

func Get(set map[string]struct{}, elem string) (string, bool) {
	if _, ok := set[elem]; ok {
		return elem, true
	}
	return "", false
}

func Append(set map[string]struct{}, elem string) {
	set[elem] = struct{}{}
}

// S + T
func Union(s, t map[string]struct{}) map[string]struct{} {
	r := make(map[string]struct{}, max(len(s), len(t)))

	for elem, _ := range s {
		r[elem] = struct{}{}
	}
	for elem, _ := range t {
		r[elem] = struct{}{}
	}

	return r
}

// S * T
func Intersection(s, t map[string]struct{}) map[string]struct{} {
	r := make(map[string]struct{})

	for elem, _ := range s {
		if Has(t, elem) {
			r[elem] = struct{}{}
		}
	}
	for elem, _ := range t {
		if Has(s, elem) {
			r[elem] = struct{}{}
		}
	}
	return r
}

// S - T
func Difference(s, t map[string]struct{}) map[string]struct{} {
	r := make(map[string]struct{})

	for elem, _ := range s {
		if !Has(t, elem) {
			r[elem] = struct{}{}
		}
	}
	return r
}

// (S - T) + (T - S)
func Symmetric_difference(s, t map[string]struct{}) map[string]struct{} {
	return Union(Difference(s, t), Difference(t, s))
}
