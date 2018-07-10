package cmd

import "testing"

func TestNegativePort(t *testing.T) {
	_, err := parsePortRange("-100")
	if err == nil {
		t.Fatal("Negative port numbers are not allowed")
	}
	_, err = parsePortRange("-1-10")
	if err == nil {
		t.Fatal("Negative starting ports are not allowed.")
	}
}

func TestOutOfOrderRange(t *testing.T) {
	_, err := parsePortRange("10-1")
	if err == nil {
		t.Fatal("Out of order ranges are not allowed.")
	}
}

func TestValidRange(t *testing.T) {
	r, err := parsePortRange("1-10")
	if err != nil {
		t.Fatalf("Failed to parse valid range: %s", err)
	}
	if len(r) != 10 {
		t.Fatalf("Unexpected length of port range: %d", len(r))
	}
}
