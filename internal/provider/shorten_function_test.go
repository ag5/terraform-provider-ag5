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
			want:       strings.Repeat("a", 58) + "-36a92",
		},
		"unicode characters": {
			name:       "vaardigheid-🚀-productie",
			maxLength:  20,
			hashLength: 5,
			want:       "vaardigheid-🚀-09f0e9",
		},
		"trailing delimiter is not doubled": {
			name:       "abcde-fghij-klmno",
			maxLength:  12,
			hashLength: 5,
			want:       "abcde-b124b4",
		},
		"no room for the name": {
			name:       "a-name-that-is-too-long",
			maxLength:  6,
			hashLength: 5,
			want:       "0b5f41",
		},
		"longer hash": {
			name:       strings.Repeat("a", 100),
			maxLength:  64,
			hashLength: 12,
			want:       strings.Repeat("a", 51) + "-36a92cc94a9e",
		},
		"minimum hash length": {
			name:       strings.Repeat("a", 100),
			maxLength:  64,
			hashLength: 1,
			want:       strings.Repeat("a", 62) + "-3",
		},
		"maximum hash length": {
			name:       strings.Repeat("a", 100),
			maxLength:  36,
			hashLength: 32,
			want:       "aaa-36a92cc94a9e0fa21f625f8bfb007adf",
		},
		"maximum too small": {
			name:       "anything",
			maxLength:  5,
			hashLength: 5,
			wantErr:    "max_length must be at least 6",
		},
		"maximum too small for longer hash": {
			name:       "anything",
			maxLength:  12,
			hashLength: 12,
			wantErr:    "max_length must be at least 13",
		},
		"hash length zero": {
			name:       strings.Repeat("a", 100),
			maxLength:  64,
			hashLength: 0,
			wantErr:    "hash_length must be between 1 and 32",
		},
		"hash length negative": {
			name:       strings.Repeat("a", 100),
			maxLength:  64,
			hashLength: -1,
			wantErr:    "hash_length must be between 1 and 32",
		},
		"hash length above digest": {
			name:       strings.Repeat("a", 100),
			maxLength:  200,
			hashLength: 33,
			wantErr:    "hash_length must be between 1 and 32",
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

// TestShortenFillsMaxLength asserts the module's property that a truncated
// result always uses the full budget.
func TestShortenFillsMaxLength(t *testing.T) {
	t.Parallel()

	for _, maxLength := range []int64{6, 7, 12, 20, 33, 64, 120} {
		got, err := shorten(strings.Repeat("ag5-", 60), maxLength, 5)
		if err != nil {
			t.Fatalf("shorten(maxLength=%d) unexpected error: %v", maxLength, err)
		}
		if int64(utf8.RuneCountInString(got)) != maxLength {
			t.Errorf("shorten(maxLength=%d) = %q, length %d, want exactly %d",
				maxLength, got, utf8.RuneCountInString(got), maxLength)
		}
	}
}
