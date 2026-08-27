package provider

import (
	"strings"
	"testing"
	"unicode/utf8"
)

func TestShorten(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		name       string
		maxLength  int64
		hashLength int64
		want       string
		wantErr    string
	}{
		"shorter than maximum": {
			name:       "ag5-production",
			maxLength:  64,
			hashLength: 5,
			want:       "ag5-production",
		},
		"exactly maximum": {
			name:       strings.Repeat("a", 64),
			maxLength:  64,
			hashLength: 5,
			want:       strings.Repeat("a", 64),
		},
		"longer than maximum": {
			name:       strings.Repeat("a", 100),
			maxLength:  64,
			hashLength: 5,
			want:       strings.Repeat("a", 58) + "-28165",
		},
		"unicode characters": {
			name:       "vaardigheid-🚀-productie",
			maxLength:  20,
			hashLength: 5,
			want:       "vaardigheid-🚀--3e0eb",
		},
		"minimum maximum": {
			name:       "a-name-that-is-too-long",
			maxLength:  9,
			hashLength: 5,
			want:       "a-n-c8400",
		},
		"longer hash": {
			name:       strings.Repeat("a", 100),
			maxLength:  64,
			hashLength: 12,
			want:       strings.Repeat("a", 51) + "-2816597888e4",
		},
		"minimum hash length": {
			name:       strings.Repeat("a", 100),
			maxLength:  64,
			hashLength: 1,
			want:       strings.Repeat("a", 62) + "-2",
		},
		"maximum hash length": {
			name:       strings.Repeat("a", 100),
			maxLength:  68,
			hashLength: 64,
			want:       "aaa-2816597888e4a0d3a36b82b83316ab32680eb8f00f8cd3b904d681246d285a0e",
		},
		"maximum too small": {
			name:       "anything",
			maxLength:  8,
			hashLength: 5,
			wantErr:    "max_length must be at least 9",
		},
		"maximum too small for longer hash": {
			name:       "anything",
			maxLength:  15,
			hashLength: 12,
			wantErr:    "max_length must be at least 16",
		},
		"hash length zero": {
			name:       strings.Repeat("a", 100),
			maxLength:  64,
			hashLength: 0,
			wantErr:    "hash_length must be between 1 and 64",
		},
		"hash length negative": {
			name:       strings.Repeat("a", 100),
			maxLength:  64,
			hashLength: -1,
			wantErr:    "hash_length must be between 1 and 64",
		},
		"hash length above digest": {
			name:       strings.Repeat("a", 100),
			maxLength:  200,
			hashLength: 65,
			wantErr:    "hash_length must be between 1 and 64",
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			t.Parallel()

			got, err := shorten(test.name, test.maxLength, test.hashLength)
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
