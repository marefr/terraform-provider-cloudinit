# Terraform Provider: cloud-init ISO

[![License](https://img.shields.io/github/license/marefr/terraform-provider-cloudinit)](LICENSE)
[![Build and Test](https://github.com/marefr/terraform-provider-cloudinit/actions/workflows/test.yml/badge.svg)](https://github.com/marefr/terraform-provider-cloudinit/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/marefr/terraform-provider-cloudinit)](https://goreportcard.com/report/github.com/marefr/terraform-provider-cloudinit)

Generate cloud-init ISO images for use with the
[NoCloud data source](https://docs.cloud-init.io/en/latest/reference/datasources/nocloud.html) and the
[drive with labeled filesystem](https://docs.cloud-init.io/en/latest/reference/datasources/nocloud.html#source-2-drive-with-labeled-filesystem)
configuration source.

Use this to provision VMs with tools like libvirt, QEMU, Proxmox, or any hypervisor that supports attaching ISO images.

See the [cloud-init documentation](https://cloudinit.readthedocs.io/) for more information.

## Usage

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

## Development

### Commands

- `make install` - Install the provider
- `make build` - Build the provider
- `make test` - Run unit tests
- `make testacc` - Run acceptance tests
- `make generate` - Generate documentation

## License

Apache 2.0 - See [LICENSE](LICENSE).

Based on work from [terraform-provider-libvirt](https://github.com/dmacvicar/terraform-provider-libvirt) - see [NOTICE.md](NOTICE.md).
