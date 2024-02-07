package mathext_test

import (
	"testing"

	"codermana.com/go/pkg/value_analysis/pkg/mathext"
)

func TestDigitCount(t *testing.T) {
	testCases := []struct {
		input    int
		expected int
	}{
		{1123512, 7},
		{12, 2},
		{0, 1},
		{1, 1},
		{100, 3},
	}

	for _, testCase := range testCases {
		actual := mathext.DigitCount(testCase.input)

		if testCase.expected != actual {
			t.Log("Expected:", testCase.expected)
			t.Log("Actual:", actual)
			t.Fail()
		}
	}
}
