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

### `shorten(name, max_length)`

Returns `name` unchanged when its length is at most `max_length`. Otherwise it
returns:

```text
{name[0:max_length-9]}-{sha256(name)[0:8]}
```

The function counts Unicode characters, hashes the UTF-8 representation of the
complete original name with SHA-256, and renders the digest as lowercase
hexadecimal. `max_length` must be at least 9.

For example, shortening 100 `a` characters to 64 characters produces:

```text
aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa-28165978
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
