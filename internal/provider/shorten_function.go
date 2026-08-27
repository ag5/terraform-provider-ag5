package provider

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/function"
)

const (
	// defaultHashLength is the hash length used when hash_length is omitted.
	defaultHashLength = 5

	// maxHashLength is the number of hexadecimal characters in a SHA-256 digest.
	maxHashLength = sha256.Size * 2

	// minPrefixLength is the number of characters of the original name that a
	// shortened result always retains, so that the result never degrades to a
	// bare hash suffix.
	minPrefixLength = 3
)

// minMaxLength returns the smallest max_length that leaves room for
// minPrefixLength name characters, the hyphen, and the hash.
func minMaxLength(hashLength int64) int64 {
	return hashLength + 1 + minPrefixLength
}

var _ function.Function = ShortenFunction{}

// ShortenFunction implements the shorten provider function.
type ShortenFunction struct{}

func NewShortenFunction() function.Function {
	return ShortenFunction{}
}

func (ShortenFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "shorten"
}

func (ShortenFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary: "Shorten a name while retaining a stable hash suffix",
		MarkdownDescription: "Returns the name unchanged when it fits. Otherwise, returns the first " +
			"`max_length - hash_length - 1` characters followed by `-` and the first `hash_length` " +
			"lowercase hexadecimal characters of the name's SHA-256 hash.",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:                "name",
				MarkdownDescription: "Name to shorten.",
			},
			function.Int64Parameter{
				Name: "max_length",
				MarkdownDescription: "Maximum number of Unicode characters in the result. Must be at least " +
					"`hash_length + 4`, so that at least three characters of the name are retained.",
			},
		},
		VariadicParameter: function.Int64Parameter{
			Name: "hash_length",
			MarkdownDescription: fmt.Sprintf(
				"Number of hexadecimal hash characters to append. Must be between 1 and %d. Defaults to %d when omitted.",
				maxHashLength, defaultHashLength,
			),
		},
		Return: function.StringReturn{},
	}
}

func (ShortenFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var name string
	var maxLength int64
	var hashLengths []int64

	resp.Error = function.ConcatFuncErrors(req.Arguments.Get(ctx, &name, &maxLength, &hashLengths))
	if resp.Error != nil {
		return
	}

	if len(hashLengths) > 1 {
		resp.Error = function.NewArgumentFuncError(2, "hash_length accepts at most one value")
		return
	}

	hashLength := int64(defaultHashLength)
	if len(hashLengths) == 1 {
		hashLength = hashLengths[0]
	}

	if err := validateHashLength(hashLength); err != nil {
		resp.Error = function.NewArgumentFuncError(2, err.Error())
		return
	}

	shortened, err := shorten(name, maxLength, hashLength)
	if err != nil {
		resp.Error = function.NewArgumentFuncError(1, err.Error())
		return
	}

	resp.Error = function.ConcatFuncErrors(resp.Result.Set(ctx, shortened))
}

// validateHashLength reports whether hashLength characters can be taken from a
// SHA-256 digest rendered as hexadecimal.
func validateHashLength(hashLength int64) error {
	if hashLength < 1 || hashLength > maxHashLength {
		return fmt.Errorf("hash_length must be between 1 and %d", maxHashLength)
	}

	return nil
}

func shorten(name string, maxLength int64, hashLength int64) (string, error) {
	if err := validateHashLength(hashLength); err != nil {
		return "", err
	}

	if minimum := minMaxLength(hashLength); maxLength < minimum {
		return "", fmt.Errorf("max_length must be at least %d", minimum)
	}

	characters := []rune(name)
	if int64(len(characters)) <= maxLength {
		return name, nil
	}

	digest := sha256.Sum256([]byte(name))
	hashPrefix := hex.EncodeToString(digest[:])[:hashLength]
	prefixLength := maxLength - hashLength - 1

	return string(characters[:prefixLength]) + "-" + hashPrefix, nil
}
