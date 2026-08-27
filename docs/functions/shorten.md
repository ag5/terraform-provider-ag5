---
page_title: "shorten function - AG5"
description: |-
  Shortens a name and appends a stable SHA-256 hash suffix.
---

# function: shorten

Returns the input unchanged when it fits within the maximum length. A longer
input is shortened to `max_length - hash_length - 1` Unicode characters and
suffixed with a hyphen and the first `hash_length` lowercase hexadecimal
characters of its SHA-256 hash. `hash_length` is optional and defaults to 5.

At least three characters of the original name are always retained, so
`max_length` must leave room for them alongside the hyphen and the hash.

## Example Usage

```terraform
output "short_name" {
  value = provider::ag5::shorten(
    "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
    64,
  )
}
```

The output is:

```text
aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa-28165
```

Pass a third argument to use a longer hash:

```terraform
output "short_name_long_hash" {
  value = provider::ag5::shorten(
    "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
    64,
    12,
  )
}
```

The output is:

```text
aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa-2816597888e4
```

## Signature

```text
shorten(name string, max_length number, hash_length number...) string
```

## Arguments

- `name` — name to shorten.
- `max_length` — maximum number of Unicode characters in the result. Must be at
  least `hash_length + 4`, so that at least three characters of the name are
  retained.
- `hash_length` — optional number of hexadecimal hash characters to append. Must
  be between 1 and 64. Defaults to 5 when omitted.
