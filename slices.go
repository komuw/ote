package main

func dedupe[T comparable](in []T) []T {
	if len(in) <= 1 {
		return in
	}

	seen := make(map[T]struct{}, len(in))
	out := make([]T, len(in))
	_ = copy(out, in)
	j := 0
	for _, v := range out {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		out[j] = v
		j++
	}
	return out[:j]
}

// difference returns the elements in `a` that aren't in `b`.
func difference[T comparable](a, b []T) []T {
	mb := make(map[T]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}
	diff := []T{}
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}
