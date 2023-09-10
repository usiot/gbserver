// BEGIN: 8f6a3b4c6d5e
package util

import (
	"bytes"
	"testing"
)

func TestRmBSpace(t *testing.T) {
	tests := []struct {
		name string
		p    []byte
		want []byte
	}{
		{
			name: "no newlines",
			p:    []byte("abc"),
			want: []byte("abc"),
		},
		{
			name: "single newline",
			p:    []byte("a\nb"),
			want: []byte("a\nb"),
		},
		{
			name: "multiple newlines",
			p:    []byte("a \n b \n c \n"),
			want: []byte("a\nb\nc\n"),
		},
		{
			name: "leading/trailing spaces",
			p:    []byte("  a \n b  \n  c  \n  "),
			want: []byte("a\nb\nc\n"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RmBSpace(tt.p)
			if !bytes.Equal(got, tt.want) {
				t.Errorf("RmBSpace() = %q, want %q", got, tt.want)
			}
		})
	}
}

// END: 8f6a3b4c6d5e
