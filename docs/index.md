---
page_title: "Provider: AG5"
description: |-
  Deterministic helper functions for AG5 Terraform configurations.
---

# AG5 Provider

The AG5 provider exposes deterministic, local-only helper functions. It has no
configuration arguments and does not communicate with remote services.

Provider-defined functions require Terraform 1.8 or newer.

## Example Usage

```terraform
terraform {
  required_version = ">= 1.8.0"

  required_providers {
    ag5 = {
      source  = "ag5/ag5"
      version = "~> 0.1"
    }
  }
}
```
