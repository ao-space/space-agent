package tools

import "testing"

func TestArrayEqual(t *testing.T) {
	a := []string{"b", "a"}
	b := []string{"a", "b"}
	if !ArrayEqual(a, b) {
		t.Fatalf("expected arrays to be equal")
	}
	// ensure inputs are not mutated
	if a[0] != "b" || a[1] != "a" {
		t.Fatalf("input slice mutated: %+v", a)
	}
}

func TestArrayContains(t *testing.T) {
	arr := []string{"x", "y"}
	if !ArrayContains(arr, "x") {
		t.Fatalf("expected to find element")
	}
	if ArrayContains(arr, "z") {
		t.Fatalf("did not expect to find element")
	}
}
