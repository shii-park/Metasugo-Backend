package sugoroku

import (
	"testing"
)

func TestRollDice(t *testing.T) {
	for i := 0; i < 100; i++ {
		result := RollDice()
		if result < 1 || result > 6 {
			t.Errorf("RollDice() out of range: got %d, want value between 1 and 6", result)
		}
	}
}
