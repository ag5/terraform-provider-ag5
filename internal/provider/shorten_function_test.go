package provider

import (
	"strings"
	"testing"
	"unicode/utf8"
)

func TestShorten(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		name      string
		maxLength int64
		want      string
		wantErr   string
	}{
		"shorter than maximum": {
			name:      "ag5-production",
			maxLength: 64,
			want:      "ag5-production",
		},
		"exactly maximum": {
			name:      strings.Repeat("a", 64),
			maxLength: 64,
			want:      strings.Repeat("a", 64),
		},
		"longer than maximum": {
			name:      strings.Repeat("a", 100),
			maxLength: 64,
			want:      strings.Repeat("a", 55) + "-28165978",
		},
		"unicode characters": {
			name:      "vaardigheid-🚀-productie",
			maxLength: 20,
			want:      "vaardigheid-3e0eb185",
		},
		"minimum maximum": {
			name:      "a-name-that-is-too-long",
			maxLength: 9,
			want:      "-c8400ded",
		},
		"maximum too small": {
			name:      "anything",
			maxLength: 8,
			wantErr:   "max_length must be at least 9",
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			t.Parallel()

			got, err := shorten(test.name, test.maxLength)
			if test.wantErr != "" {
				if err == nil || err.Error() != test.wantErr {
					t.Fatalf("shorten() error = %v, want %q", err, test.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("shorten() unexpected error: %v", err)
			}
			if got != test.want {
				t.Errorf("shorten() = %q, want %q", got, test.want)
			}
			if utf8.RuneCountInString(got) > int(test.maxLength) {
				t.Errorf("shorten() length = %d, exceeds %d", utf8.RuneCountInString(got), test.maxLength)
			}
		})
	}
}
