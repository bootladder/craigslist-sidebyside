package main

import "testing"

func TestSum(t *testing.T) {
	t.Errorf("Sum was incorrect, got: %d, want: %d.", 39, 10)
}
