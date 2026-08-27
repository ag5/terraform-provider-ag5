package provider

import (
	"context"
	"crypto/md5" // #nosec G501 -- identity hash, chosen for parity with terraform-null-label
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/function"
)

// The suffix mirrors cloudposse/terraform-null-label's id_length_limit
// behavior, so that a shortened name matches the id that module produces for
// the same input:
//
//	id_truncated_length_limit = id_length_limit - (id_hash_length + delimiter_length)
//	id_truncated              = trimsuffix(substr(id_full, 0, ...), delimiter) + delimiter
//	id_hash                   = lower(md5(id_full) + "qrstuvwxyz")
//	id_short                  = substr(id_truncated + id_hash, 0, id_length_limit)
const (
	// defaultHashLength mirrors the module's defaults.id_hash_length.
	defaultHashLength = 5

	// maxHashLength is the number of hexadecimal characters in an MD5 digest.
	maxHashLength = md5.Size * 2

	// delimiter separates the truncated name from the hash.
	delimiter = "-"

	// hashPadding mirrors the module's id_hash_plus. It guarantees the hash is
	// long enough to fill the result up to max_length, even when a trailing
	// delimiter was trimmed from the truncated name.
	hashPadding = "qrstuvwxyz"
)

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
		MarkdownDescription: "Returns the name unchanged when it fits. Otherwise, truncates the name to " +
			"`max_length - hash_length - 1` characters, drops a trailing `-`, and fills the result up to " +
			"`max_length` characters with the name's lowercase hexadecimal MD5 hash. This matches the `id` " +
			"that [cloudposse/terraform-null-label](https://registry.terraform.io/modules/cloudposse/label/null) " +
			"produces when `id_length_limit` is set.",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:                "name",
				MarkdownDescription: "Name to shorten.",
			},
			function.Int64Parameter{
				Name:                "max_length",
				MarkdownDescription: "Maximum number of Unicode characters in the result. Must be at least `hash_length + 1`.",
			},
		},
		VariadicParameter: function.Int64Parameter{
			Name: "hash_length",
			MarkdownDescription: fmt.Sprintf(
				"Number of hash characters to reserve. Must be between 1 and %d. Defaults to %d when omitted, "+
					"matching the module's `id_hash_length`.",
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

// validateHashLength reports whether hashLength characters can be reserved from
// an MD5 digest rendered as hexadecimal.
func validateHashLength(hashLength int64) error {
	if hashLength < 1 || hashLength > maxHashLength {
		return fmt.Errorf("hash_length must be between 1 and %d", maxHashLength)
	}

	return nil
}

// minMaxLength returns the smallest max_length that leaves room for the
// delimiter and the hash. The module enforces the same minimum, expressed as a
// literal 6 for its fixed hash length of 5.
func minMaxLength(hashLength int64) int64 {
	return hashLength + int64(len(delimiter))
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

	digest := md5.Sum([]byte(name)) // #nosec G401 -- identity hash, not a security control
	hash := hex.EncodeToString(digest[:]) + hashPadding

	// Truncate to leave room for the delimiter and the hash. A trailing
	// delimiter is dropped rather than doubled, which leaves one more character
	// for the hash to fill.
	var truncated string
	if truncatedLength := maxLength - minMaxLength(hashLength); truncatedLength > 0 {
		truncated = strings.TrimSuffix(string(characters[:truncatedLength]), delimiter) + delimiter
	}

	result := []rune(truncated + hash)
	if int64(len(result)) > maxLength {
		result = result[:maxLength]
	}

	return string(result), nil
}
