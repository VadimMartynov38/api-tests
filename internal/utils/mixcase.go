package utils

import "strings"

func MixCase(s string) string {
	var b strings.Builder
	upper := false
	for _, r := range s {
		if upper {
			b.WriteRune(r)
		} else {
			b.WriteString(strings.ToLower(string(r)))
		}
		upper = !upper
	}
	return b.String()
}
