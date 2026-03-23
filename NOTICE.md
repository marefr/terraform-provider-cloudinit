This is based on work in https://github.com/dmacvicar/terraform-provider-libvirt by @dmacvicar, especially the `libvirt_cloudinit_disk` resource and implementation.

This provider/resource was implemented because even though the resource `libvirt_cloudinit_disk` states it doesn't require any provider configuration, it still does and fails if not provided.
