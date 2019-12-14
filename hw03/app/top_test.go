package app

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func generateText(words map[string]int) string {
	var text []string
	for k, v := range words {
		for i := 0; i < v; i++ {
			text = append(text, k)
		}
	}
	return strings.Join(text, " ")
}

func Test_Top10_With_Empty_Text(t *testing.T) {
	text := ""
	res := Top10(text)
	assert.Equal(t, 0, len(res))
}

func Test_Top10_With_Small_Text(t *testing.T) {
	text := "It is simple text with simple phrase"
	res := Top10(text)
	assert.Equal(t, 6, len(res))
	assert.Equal(t, "simple", res[0])
}

func Test_Top10_With_Normal_Text(t *testing.T) {
	text := generateText(map[string]int{
		"one":    12,
		"two":    11,
		"three":  10,
		"four":   9,
		"five":   8,
		"six":    7,
		"seven":  6,
		"eight":  5,
		"nine":   4,
		"ten":    3,
		"eleven": 2,
		"twelve": 1,
	})
	res := Top10(text)
	assert.Equal(t, 10, len(res))
	assert.Equal(t, "one", res[0])
	assert.Equal(t, "two", res[1])
	assert.Equal(t, "three", res[2])
	assert.Equal(t, "four", res[3])
	assert.Equal(t, "five", res[4])
	assert.Equal(t, "six", res[5])
	assert.Equal(t, "seven", res[6])
	assert.Equal(t, "eight", res[7])
	assert.Equal(t, "nine", res[8])
	assert.Equal(t, "ten", res[9])
}

func Test_Top10_With_Normal_Text2(t *testing.T) {
	text := generateText(map[string]int{
		"one":    12,
		"two":    11,
		"three":  10,
		"four":   9,
		"five":   8,
		"six":    8,
		"seven":  8,
		"eight":  5,
		"nine":   4,
		"ten":    3,
		"eleven": 2,
		"twelve": 1,
	})
	res := Top10(text)
	assert.Equal(t, 10, len(res))
	assert.Equal(t, "one", res[0])
	assert.Equal(t, "two", res[1])
	assert.Equal(t, "three", res[2])
	assert.Equal(t, "four", res[3])
	assert.Contains(t, []string{"five", "six", "seven"}, res[4])
	assert.Contains(t, []string{"five", "six", "seven"}, res[5])
	assert.Contains(t, []string{"five", "six", "seven"}, res[6])
	assert.Equal(t, "eight", res[7])
	assert.Equal(t, "nine", res[8])
	assert.Equal(t, "ten", res[9])
}

func Test_Top10_With_Register(t *testing.T) {
	text := generateText(map[string]int{
		"one":    5,
		"ONe":    7,
		"two":    11,
		"three":  10,
		"four":   9,
		"five":   8,
		"siX":    3,
		"Six":    3,
		"sIx":    2,
		"seven":  8,
		"eight":  5,
		"nine":   4,
		"ten":    3,
		"eleven": 1,
		"Eleven": 1,
		"twelve": 1,
	})
	res := Top10(text)
	assert.Equal(t, 10, len(res))
	assert.Equal(t, "one", res[0])
	assert.Equal(t, "two", res[1])
	assert.Equal(t, "three", res[2])
	assert.Equal(t, "four", res[3])
	assert.Contains(t, []string{"five", "six", "seven"}, res[4])
	assert.Contains(t, []string{"five", "six", "seven"}, res[5])
	assert.Contains(t, []string{"five", "six", "seven"}, res[6])
	assert.Equal(t, "eight", res[7])
	assert.Equal(t, "nine", res[8])
	assert.Equal(t, "ten", res[9])
}

func Test_Top10_With_Marks(t *testing.T) {
	text := generateText(map[string]int{
		"one,":    5,
		"ONe?":    7,
		"two.":    11,
		"three!":  10,
		"four:":   9,
		"five;":   8,
		"siX-":    3,
		"Six":     3,
		"sIx":     2,
		"seven\n": 8,
		"eight\t": 5,
		"nine":    4,
		"ten":     3,
		"eleven":  1,
		"Eleven":  1,
		"twelve":  1,
		"-":       23,
	})
	res := Top10(text)
	assert.Equal(t, 10, len(res))
	assert.Equal(t, "one", res[0])
	assert.Equal(t, "two", res[1])
	assert.Equal(t, "three", res[2])
	assert.Equal(t, "four", res[3])
	assert.Contains(t, []string{"five", "six", "seven"}, res[4])
	assert.Contains(t, []string{"five", "six", "seven"}, res[5])
	assert.Contains(t, []string{"five", "six", "seven"}, res[6])
	assert.Equal(t, "eight", res[7])
	assert.Equal(t, "nine", res[8])
	assert.Equal(t, "ten", res[9])
}
