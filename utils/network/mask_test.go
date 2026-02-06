package network

import "testing"

func TestSubNetMaskToLen(t *testing.T) {
	cases := map[string]int{
		"255.255.255.0":   24,
		"255.255.0.0":     16,
		"255.255.255.255": 32,
		"0.0.0.0":         0,
	}
	for mask, want := range cases {
		got, err := SubNetMaskToLen(mask)
		if err != nil {
			t.Fatalf("unexpected error for %s: %v", mask, err)
		}
		if got != want {
			t.Fatalf("mask %s: expected %d, got %d", mask, want, got)
		}
	}
}

func TestSubNetMaskToLenInvalid(t *testing.T) {
	if _, err := SubNetMaskToLen("255.255.255"); err == nil {
		t.Fatalf("expected error for invalid mask")
	}
	if _, err := SubNetMaskToLen("256.0.0.0"); err == nil {
		t.Fatalf("expected error for out-of-range mask")
	}
}

func TestLenToSubNetMask(t *testing.T) {
	if got := LenToSubNetMask(24); got != "255.255.255.0" {
		t.Fatalf("expected 255.255.255.0, got %s", got)
	}
	if got := LenToSubNetMask(0); got != "0.0.0.0" {
		t.Fatalf("expected 0.0.0.0, got %s", got)
	}
	if got := LenToSubNetMask(32); got != "255.255.255.255" {
		t.Fatalf("expected 255.255.255.255, got %s", got)
	}
}
