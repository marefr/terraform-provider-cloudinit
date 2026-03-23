# Terraform Provider: cloud-init ISO

[![License](https://img.shields.io/github/license/marefr/terraform-provider-cloudinit)](LICENSE)
[![GitHub release](https://img.shields.io/github/v/release/marefr/terraform-provider-cloudinit)](https://github.com/marefr/terraform-provider-cloudinit/releases)
[![Terraform Provider Downloads](https://img.shields.io/terraform/provider/dt/1659142)](https://registry.terraform.io/providers/marefr/cloudinit/latest)
[![CI](https://github.com/marefr/terraform-provider-cloudinit/actions/workflows/test.yml/badge.svg)](https://github.com/marefr/terraform-provider-cloudinit/actions/workflows/test.yml)
[![Go](https://img.shields.io/github/go-mod/go-version/marefr/terraform-provider-cloudinit)](go.mod)
[![Go Report Card](https://goreportcard.com/badge/github.com/marefr/terraform-provider-cloudinit)](https://goreportcard.com/report/github.com/marefr/terraform-provider-cloudinit)

A Terraform/Opentofu provider for generating cloud-init ISO images for use with the
[NoCloud data source](https://docs.cloud-init.io/en/latest/reference/datasources/nocloud.html) and the
[drive with labeled filesystem](https://docs.cloud-init.io/en/latest/reference/datasources/nocloud.html#source-2-drive-with-labeled-filesystem)
configuration source.

Use this to provision VMs with tools like libvirt, QEMU, Proxmox, or any hypervisor that supports attaching ISO images.

## Installation

```terraform
terraform {
  required_providers {
    cloudinit = {
      source = "marefr/cloudinit"
    }
  }
}
```

## Usage

```terraform
resource "cloudinit_iso" "seed" {
  name      = "vm-init"
  user_data = file("user-data.yaml")
  meta_data = yamlencode({
    instance-id    = "vm-01"
    local-hostname = "webserver"
  })
}
```

## Documentation

Full documentation is available on the [Opentofu Registry](https://search.opentofu.org/provider/marefr/cloudinit/latest) and the [Terraform Registry](https://registry.terraform.io/providers/marefr/cloudinit/latest/docs).

## Development

### Commands

- `make install` - Install the provider
- `make build` - Build the provider
- `make test` - Run unit tests
- `make testacc` - Run acceptance tests
- `make generate` - Generate documentation

### Use of AI

The author uses AI for this project, but as the maintainer, he owns the outcome and consequences.

## License

Apache 2.0 - See [LICENSE](LICENSE).

Based on work from [terraform-provider-libvirt](https://github.com/dmacvicar/terraform-provider-libvirt) - see [NOTICE.md](NOTICE.md).
