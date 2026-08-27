---
page_title: "shorten function - AG5"
description: |-
  Shortens a name and appends a stable MD5 hash suffix.
---

# function: shorten

Returns the input unchanged when it fits within the maximum length. A longer
input is truncated to `max_length - hash_length - 1` Unicode characters, loses a
trailing `-` rather than doubling it, and is then filled up to exactly
`max_length` characters with the lowercase hexadecimal MD5 hash of the complete
original name.

This mirrors the `id` that
[cloudposse/terraform-null-label](https://registry.terraform.io/modules/cloudposse/label/null)
produces when `id_length_limit` is set, so both can be used on the same
resources without changing their names.

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
aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa-36a92
```

Pass a third argument to reserve more hash characters:

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
aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa-36a92cc94a9e
```

## Signature

```text
shorten(name string, max_length number, hash_length number...) string
```

## Arguments

- `name` — name to shorten.
- `max_length` — maximum number of Unicode characters in the result. Must be at
  least `hash_length + 1`. A truncated result always uses the full length.
- `hash_length` — optional number of hash characters to reserve. Must be between
  1 and 32. Defaults to 5 when omitted, matching the module's `id_hash_length`.
  One extra hash character appears when a trailing `-` is trimmed from the
  truncated name.
