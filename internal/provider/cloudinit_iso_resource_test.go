package provider

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCloudInitDiskResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccCloudInitDiskResourceConfigBasic("test-cloudinit"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cloudinit_iso.test", "name", "test-cloudinit"),
					resource.TestCheckResourceAttrSet("cloudinit_iso.test", "id"),
					resource.TestCheckResourceAttrSet("cloudinit_iso.test", "path"),
					resource.TestCheckResourceAttrSet("cloudinit_iso.test", "size"),
					testAccCheckCloudInitDiskExists("cloudinit_iso.test"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccCloudInitDiskResource_withNetworkConfig(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudInitDiskResourceConfigWithNetwork("test-cloudinit-net"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cloudinit_iso.test", "name", "test-cloudinit-net"),
					resource.TestCheckResourceAttrSet("cloudinit_iso.test", "network_config"),
					resource.TestCheckResourceAttrSet("cloudinit_iso.test", "path"),
					testAccCheckCloudInitDiskExists("cloudinit_iso.test"),
				),
			},
		},
	})
}

// testAccCheckCloudInitDiskExists verifies that the ISO file exists on disk.
func testAccCheckCloudInitDiskExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set for %s", resourceName)
		}

		path := rs.Primary.Attributes["path"]
		if path == "" {
			return fmt.Errorf("no path is set for %s", resourceName)
		}

		// Check if the ISO file exists
		if _, err := os.Stat(path); err != nil {
			return fmt.Errorf("ISO file does not exist at %s: %w", path, err)
		}

		// Verify it's in the expected temp directory
		expectedDir := filepath.Join(os.TempDir(), "terraform-provider-cloudinit")
		if !strings.HasPrefix(path, expectedDir) {
			return fmt.Errorf("ISO file is not in expected directory: got %s, expected prefix %s", path, expectedDir)
		}

		return nil
	}
}

func testAccCloudInitDiskResourceConfigBasic(name string) string {
	return fmt.Sprintf(`
resource "cloudinit_iso" "test" {
  name      = %[1]q
  user_data = <<-EOF
    #cloud-config
    users:
      - name: root
        ssh_authorized_keys:
          - ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC0/Ho1w+1D4vJccMzEQBzREzCY4NjkrJYh8+9rQJgDrYPrWLe1PJYvDG6r1uDlLrZJhwwq1PcJQw test@example.com
  EOF

  meta_data = <<-EOF
    instance-id: %[1]s-001
    local-hostname: %[1]s
  EOF
}
`, name)
}

func testAccCloudInitDiskResourceConfigWithNetwork(name string) string {
	return fmt.Sprintf(`
resource "cloudinit_iso" "test" {
  name      = %[1]q
  user_data = <<-EOF
    #cloud-config
    users:
      - name: root
        ssh_authorized_keys:
          - ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC0/Ho1w+1D4vJccMzEQBzREzCY4NjkrJYh8+9rQJgDrYPrWLe1PJYvDG6r1uDlLrZJhwwq1PcJQw test@example.com
  EOF

  meta_data = <<-EOF
    instance-id: %[1]s-001
    local-hostname: %[1]s
  EOF

  network_config = <<-EOF
    version: 2
    ethernets:
      eth0:
        dhcp4: true
  EOF
}
`, name)
}
