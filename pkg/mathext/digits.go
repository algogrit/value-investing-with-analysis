package mathext

func DigitCount(n int) int {
	count := 1
	for n/10 != 0 {
		count++
		n = n / 10
	}
	return count
}
