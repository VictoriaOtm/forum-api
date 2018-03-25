package helpers

func isDigit(r rune) bool {
	return (r - '0') < 10
}

func IsNumber(s string) bool {
	for _, r := range s {
		if !isDigit(r) {
			return false
		}
	}

	return true
}