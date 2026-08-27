---
page_title: "shorten function - AG5"
description: |-
  Shortens a name and appends a stable SHA-256 hash suffix.
---

# function: shorten

Returns the input unchanged when it fits within the maximum length. A longer
input is shortened to `max_length - 9` Unicode characters and suffixed with a
hyphen and the first eight lowercase hexadecimal characters of its SHA-256
hash.

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
aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa-28165978
```

## Signature

```text
shorten(name string, max_length number) string
```

## Arguments

- `name` — name to shorten.
- `max_length` — maximum number of Unicode characters in the result. Must be at
  least 9.
