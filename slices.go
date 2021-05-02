package main

func dedupe(in []string) []string {
	if len(in) <= 1 {
		return in
	}

	seen := make(map[string]struct{}, len(in))
	j := 0
	for _, v := range in {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		in[j] = v
		j++
	}
	return in[:j]
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
