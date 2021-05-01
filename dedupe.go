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
