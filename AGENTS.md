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

See [docs/resources/iso.md](docs/resources/iso.md) for full schema documentation and examples.

**Behavior:**
- All content attributes use `RequiresReplace()` - any change triggers resource recreation
- `user_data` is marked `Sensitive: true` to prevent exposure in logs/plans
- ISO files are named `cloudinit-{id}.iso` for uniqueness
- The ID is derived from SHA256 of: name + user_data + meta_data + network_config (if set)

Uses the `github.com/kdomanski/iso9660` library to create ISO9660 images with volume label `cidata`.

## Testing

```bash
make test      # Unit tests
make testacc   # Acceptance tests (creates real ISO files)
```

Key test helpers in `*_test.go` files:
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

## Important Notes

1. **Documentation is auto-generated** - Never edit files in `docs/` directly. Modify schema descriptions and run `make generate`.

2. **ISO files in temp directory** - By default, ISOs are created in `$TMPDIR/terraform-provider-cloudinit/`. These may be cleaned on reboot. Use `output_dir` for persistence.

3. **Immutable resources** - The resource is designed to be immutable. Any content change creates a new ISO with a new ID.

4. **Sensitive data** - `user_data` often contains secrets (passwords, SSH keys). It's marked sensitive to prevent exposure.
