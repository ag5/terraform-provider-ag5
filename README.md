# AG5 Terraform Provider

The AG5 provider contains deterministic helper functions for AG5 Terraform
configurations. It does not require configuration and does not communicate with
remote services.

Provider-defined functions require Terraform 1.8 or newer.

## Usage

```hcl
terraform {
  required_version = ">= 1.8.0"

  required_providers {
    ag5 = {
      source  = "ag5/ag5"
      version = "~> 0.1"
    }
  }
}

output "bucket_name" {
  value = provider::ag5::shorten(
    "a-name-that-may-be-longer-than-the-target-platform-allows",
    64,
  )
}
```

## Functions

### `shorten(name, max_length, [hash_length])`

Returns `name` unchanged when its length is at most `max_length`. Otherwise it
truncates the name to `max_length - hash_length - 1` characters, drops a
trailing `-` instead of doubling it, and fills the result up to exactly
`max_length` characters with the lowercase hexadecimal MD5 hash of the complete
original name.

This matches the `id` that
[cloudposse/terraform-null-label](https://registry.terraform.io/modules/cloudposse/label/null)
produces when `id_length_limit` is set.

The function counts Unicode characters and hashes the UTF-8 representation of
the original name. `hash_length` is optional, defaults to 5, and must be between
1 and 32. `max_length` must be at least `hash_length + 1`.

For example, shortening 100 `a` characters to 64 characters produces:

```text
aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa-36a92
```

## Development

```shell
go test ./...
go build -o terraform-provider-ag5 .
```

Run the Terraform-level function tests with:

```shell
make testacc
```
