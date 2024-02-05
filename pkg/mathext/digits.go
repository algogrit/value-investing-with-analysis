package mathext

func DigitCount(n int) int {
	if n/10 == 0 {
		return 1
	}
	return 1 + DigitCount(n/10)
}
