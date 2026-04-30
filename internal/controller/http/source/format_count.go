package source

import "strconv"

// formatCount renders n with thousands separators (e.g. 12481 -> "12,481").
func formatCount(n int) string {
	s := strconv.Itoa(n)
	if len(s) <= 3 {
		return s
	}

	out := make([]byte, 0, len(s)+len(s)/3)

	for i := 0; i < len(s); i++ {
		if i > 0 && (len(s)-i)%3 == 0 {
			out = append(out, ',')
		}

		out = append(out, s[i])
	}

	return string(out)
}
