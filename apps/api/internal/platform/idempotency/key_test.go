package idempotency

import "testing"

func TestParse(t *testing.T) {
	if _, err := Parse("order:create:0123456789"); err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	for _, value := range []string{"short", "unsafe key with spaces", "line\nbreak"} {
		if _, err := Parse(value); err == nil {
			t.Fatalf("Parse(%q) expected an error", value)
		}
	}
}
