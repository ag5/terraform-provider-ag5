terraform {
  required_version = ">= 1.8.0"

  required_providers {
    ag5 = {
      source  = "ag5/ag5"
      version = "~> 0.1"
    }
  }
}

output "short_name" {
  value = provider::ag5::shorten(
    "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
    64,
  )
}

output "short_name_long_hash" {
  value = provider::ag5::shorten(
    "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
    64,
    12,
  )
}
