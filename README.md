# Terraform Provider for generating cloud-init iso

[![License](https://img.shields.io/github/license/marefr/terraform-provider-cloudinit)](LICENSE)
[![Build and Test](https://github.com/marefr/terraform-provider-cloudinit/actions/workflows/test.yml/badge.svg)](https://github.com/marefr/terraform-provider-cloudinit/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/marefr/terraform-provider-cloudinit)](https://goreportcard.com/report/github.com/marefr/terraform-provider-cloudinit)

The provider allows generating cloud-init ISO images, suitable for use with the
[NoCloud data source](https://docs.cloud-init.io/en/latest/reference/datasources/nocloud.html) and the
[drive with labeled filesystem](https://docs.cloud-init.io/en/latest/reference/datasources/nocloud.html#source-2-drive-with-labeled-filesystem)
configuration source.

See the [cloud-init documentation](https://cloudinit.readthedocs.io/) for more details.

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.24

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the `make build` command:

```shell
make build
```

## Using the provider

```terraform
terraform {
  required_providers {
    cloudinit = {
      source = "marefr/cloudinit"
    }
  }
}

resource "cloudinit_iso" "seed" {
  name      = "vm-init"
  user_data = file("user-data.yaml")
  meta_data = yamlencode({
    instance-id    = "vm-01"
    local-hostname = "webserver"
  })
}
```

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `make generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

_Note:_ Acceptance tests create real resources, and often cost money to run.

```shell
make testacc
```
