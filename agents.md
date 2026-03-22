# AI Agent Instructions

This document provides guidance for AI agents working with the `terraform-provider-cloudinit` codebase.

## Project Overview

This is a Terraform provider that generates cloud-init ISO images for use with the [NoCloud datasource](https://docs.cloud-init.io/en/latest/reference/datasources/nocloud.html). The ISOs are created with the `cidata` volume label, which cloud-init automatically detects when attached to a VM.

**Primary use case:** Provisioning VMs with tools like libvirt, QEMU, Proxmox, or any hypervisor that supports attaching ISO images.

## Project Structure

```
.
├── internal/provider/           # Provider and resource implementations
│   ├── provider.go              # Main provider definition
│   ├── provider_test.go         # Provider test setup
│   ├── cloudinit_iso_resource.go      # The cloudinit_iso resource
│   └── cloudinit_iso_resource_test.go # Resource tests
├── docs/                        # Auto-generated documentation (do not edit manually)
├── examples/                    # Example Terraform configurations
│   ├── provider/                # Provider configuration example
│   └── resources/               # Resource examples
├── tools/                       # Code generation tools
├── main.go                      # Provider entry point
├── Makefile                     # Build commands
└── go.mod                       # Go module definition
```

## Development Commands

| Command | Description |
|---------|-------------|
| `make build` | Build the provider binary |
| `make test` | Run unit tests |
| `make testacc` | Run acceptance tests (creates real ISO files) |
| `make lint` | Run golangci-lint |
| `make fmt` | Format Go code |
| `make generate` | Generate documentation from schema |
| `make` | Run fmt, lint, install, and generate |

## Key Implementation Details

### Resource: `cloudinit_iso`

The only resource in this provider. Located in `internal/provider/cloudinit_iso_resource.go`.

**Attributes:**

| Attribute | Type | Required | Description |
|-----------|------|----------|-------------|
| `name` | string | Yes | Resource name, included in ID hash |
| `user_data` | string | Yes | Cloud-init user-data (sensitive) |
| `meta_data` | string | Yes | Cloud-init meta-data |
| `network_config` | string | No | Cloud-init network configuration |
| `output_dir` | string | No | Custom output directory (defaults to temp) |
| `id` | string | Computed | SHA256 hash (first 16 chars) of name + content |
| `path` | string | Computed | Full path to generated ISO |
| `size` | int64 | Computed | ISO file size in bytes |

**Behavior:**
- All content attributes use `RequiresReplace()` - any change triggers resource recreation
- `user_data` is marked `Sensitive: true` to prevent exposure in logs/plans
- ISO files are named `cloudinit-{id}.iso` for uniqueness
- The ID is derived from SHA256 of: name + user_data + meta_data + network_config (if set)

### ISO Generation

Uses the `github.com/kdomanski/iso9660` library to create ISO9660 images with:
- Volume label: `cidata` (required for cloud-init NoCloud detection)
- Files: `user-data`, `meta-data`, and optionally `network-config`

## Testing

### Unit Tests
```bash
make test
```

### Acceptance Tests
```bash
make testacc
```

Acceptance tests create real ISO files in a temporary directory. They verify:
- ISO file creation
- Correct file paths
- Attribute values in state

### Test Patterns

Test files follow the pattern `*_test.go`. Key test helpers:
- `testAccProtoV6ProviderFactories` - Provider factory for tests
- `testAccPreCheck(t)` - Pre-test validation
- `testAccCheckCloudInitDiskExists` - Verifies ISO file exists on disk

## Common Tasks

### Adding a New Attribute

1. Add field to `CloudInitISOResourceModel` struct
2. Add schema attribute in `Schema()` method
3. Update `Create()` to use the new attribute
4. Update `Read()` if the attribute affects state
5. Add test coverage
6. Run `make generate` to update docs

### Modifying ISO Content

The `generateCloudInitISO()` method handles ISO creation. To add new files:
1. Add the file using `writer.AddFile()`
2. Ensure the attribute is included in the ID hash if it affects content

### Debugging

Enable Terraform debug logging:
```bash
TF_LOG=DEBUG terraform apply
```

Provider logs use `tflog` package with structured logging.

## Code Style

- Follow standard Go conventions
- Use Terraform Plugin Framework patterns
- Error messages should be user-friendly (shown in Terraform output)
- Use `tflog` for debug/info logging, not `fmt.Print` or `log`
- Schema descriptions are used for documentation generation

## Example Usage

### Basic Configuration

```terraform
resource "cloudinit_iso" "example" {
  name      = "my-vm-init"
  user_data = file("cloud-init/user-data.yaml")
  meta_data = yamlencode({
    instance-id    = "vm-001"
    local-hostname = "webserver"
  })
}
```

### With Network Configuration

```terraform
resource "cloudinit_iso" "example" {
  name      = "my-vm-init"
  user_data = file("cloud-init/user-data.yaml")
  meta_data = file("cloud-init/meta-data.yaml")
  network_config = file("cloud-init/network-config.yaml")
}
```

### With Custom Output Directory

```terraform
resource "cloudinit_iso" "example" {
  name       = "my-vm-init"
  output_dir = "/var/lib/libvirt/images"
  user_data  = file("cloud-init/user-data.yaml")
  meta_data  = file("cloud-init/meta-data.yaml")
}
```

## Dependencies

Key dependencies (see `go.mod`):
- `github.com/hashicorp/terraform-plugin-framework` - Terraform provider SDK
- `github.com/hashicorp/terraform-plugin-testing` - Testing utilities
- `github.com/hashicorp/terraform-plugin-log` - Structured logging
- `github.com/kdomanski/iso9660` - ISO9660 image generation

## CI/CD

GitHub Actions workflows in `.github/workflows/`:
- `test.yml` - Runs on PRs and pushes to main; builds, lints, and runs acceptance tests
- `release.yml` - Handles releases via GoReleaser

## Important Notes

1. **Documentation is auto-generated** - Never edit files in `docs/` directly. Modify schema descriptions and run `make generate`.

2. **ISO files in temp directory** - By default, ISOs are created in `$TMPDIR/terraform-provider-cloudinit/`. These may be cleaned on reboot. Use `output_dir` for persistence.

3. **Immutable resources** - The resource is designed to be immutable. Any content change creates a new ISO with a new ID.

4. **Sensitive data** - `user_data` often contains secrets (passwords, SSH keys). It's marked sensitive to prevent exposure.
