package main

func dedupe(in []string) []string {
	if len(in) <= 0 {
		return in
	} else if len(in) <= 3 {
		// deduping a small slice is probably wasteful
		return in
	}

	j := 0
	for i := 1; i < len(in); i++ {
		if in[j] == in[i] {
			continue
		}
		j++
		// preserve the original data
		// in[i], in[j] = in[j], in[i]
		// only set what is required
		in[j] = in[i]
	}
	result := in[:j+1]
	return result
}

// difference returns the elements in `a` that aren't in `b`.
func difference(a, b []string) []string {
	mb := make(map[string]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}
	var diff = []string{}
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}
