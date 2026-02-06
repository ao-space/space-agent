package tools

import "testing"

func TestStringToLines(t *testing.T) {
	input := "a\r\nb\rc\n"
	lines := StringToLines(input)
	if len(lines) != 4 {
		t.Fatalf("unexpected line count: %d", len(lines))
	}
	if lines[0] != "a" || lines[1] != "b" || lines[2] != "c" || lines[3] != "" {
		t.Fatalf("unexpected lines: %+v", lines)
	}
}
