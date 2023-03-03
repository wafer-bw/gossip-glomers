package stack

func Push[T any](slice *[]T, i T) {
	*slice = append(*slice, i)
}

func Pop[T any](slice *[]T) (T, bool) {
	s := *slice
	if len(s) == 0 {
		var res T
		return res, false
	}

	res := s[len(s)-1]
	*slice = s[:len(s)-1]

	return res, true
}
