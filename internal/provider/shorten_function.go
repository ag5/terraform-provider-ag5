package provider

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/function"
)

const hashSuffixLength = 9

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
			"`max_length - 9` characters followed by `-` and the first eight lowercase hexadecimal " +
			"characters of the name's SHA-256 hash.",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:                "name",
				MarkdownDescription: "Name to shorten.",
			},
			function.Int64Parameter{
				Name:                "max_length",
				MarkdownDescription: "Maximum number of Unicode characters in the result. Must be at least 9.",
			},
		},
		Return: function.StringReturn{},
	}
}

func (ShortenFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var name string
	var maxLength int64

	resp.Error = function.ConcatFuncErrors(req.Arguments.Get(ctx, &name, &maxLength))
	if resp.Error != nil {
		return
	}

	shortened, err := shorten(name, maxLength)
	if err != nil {
		resp.Error = function.NewArgumentFuncError(1, err.Error())
		return
	}

	resp.Error = function.ConcatFuncErrors(resp.Result.Set(ctx, shortened))
}

func shorten(name string, maxLength int64) (string, error) {
	if maxLength < hashSuffixLength {
		return "", fmt.Errorf("max_length must be at least %d", hashSuffixLength)
	}

	characters := []rune(name)
	if int64(len(characters)) <= maxLength {
		return name, nil
	}

	digest := sha256.Sum256([]byte(name))
	hashPrefix := hex.EncodeToString(digest[:])[:hashSuffixLength-1]
	prefixLength := int(maxLength - hashSuffixLength)

	return string(characters[:prefixLength]) + "-" + hashPrefix, nil
}
