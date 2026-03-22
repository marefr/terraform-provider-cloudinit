resource "cloudinit_iso" "seed" {
  name       = "vm-init"
  output_dir = "/var/lib/libvirt/images"
  user_data  = file("user-data.yaml")
  meta_data = yamlencode({
    instance-id    = "vm-01"
    local-hostname = "webserver"
  })
}
