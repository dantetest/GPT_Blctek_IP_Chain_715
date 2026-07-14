package manifestspec

import (
	"errors"
	"testing"
)

func TestCanonicalPath(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr error
	}{
		{name: "windows separator", input: `nested\file.txt`, want: "nested/file.txt"},
		{name: "unicode NFC", input: "nested/e\u0301.txt", want: "nested/é.txt"},
		{name: "absolute unix", input: "/etc/passwd", wantErr: ErrInvalidPath},
		{name: "absolute windows", input: `C:\data\file.txt`, wantErr: ErrInvalidPath},
		{name: "traversal", input: "../file.txt", wantErr: ErrInvalidPath},
		{name: "empty segment", input: "a//b.txt", wantErr: ErrInvalidPath},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := CanonicalPath(test.input)
			if test.wantErr != nil {
				if !errors.Is(err, test.wantErr) {
					t.Fatalf("error=%v", err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if got != test.want {
				t.Fatalf("got=%q want=%q", got, test.want)
			}
		})
	}
}
