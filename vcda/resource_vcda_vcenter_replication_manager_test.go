package vcda

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"
	"testing"
)

func TestAccVcdaVcenterReplicationManager_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccVcdaVcenterReplicationManagerPreCheck(t)
		},
		ProviderFactories: testProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccVcdaVcenterReplicationManagerConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vcda_vcenter_replication_manager.manager_site", "is_licensed", "true"),
					resource.TestCheckResourceAttr("vcda_vcenter_replication_manager.manager_site", "expiration_date", "0"),
					resource.TestCheckResourceAttr("vcda_vcenter_replication_manager.manager_site", "ls_url", os.Getenv(LookupServiceURL)),

					resource.TestCheckResourceAttrSet("vcda_vcenter_replication_manager.manager_site", "ls_thumbprint"),
					resource.TestCheckResourceAttrSet("vcda_vcenter_replication_manager.manager_site", "tunnel_url"),
					resource.TestCheckResourceAttrSet("vcda_vcenter_replication_manager.manager_site", "tunnel_certificate"),
				),
			},
		},
	})
}

func testAccVcdaVcenterReplicationManagerPreCheck(t *testing.T) {
	if v := os.Getenv(ManagerVmName); v == "" {
		t.Fatal(ManagerVmName + " must be set for vcda_vcenter_replication_manager acceptance tests")
	}
	if os.Getenv(LookupServiceAddress) == "" {
		t.Fatal(LookupServiceAddress + " must be set for vcda_vcenter_replication_manager acceptance tests")
	}
	if os.Getenv(LicenseKey) == "" {
		t.Fatal(LicenseKey + " must be set for vcda_vcenter_replication_manager acceptance tests")
	}
	if os.Getenv(LookupServiceURL) == "" {
		t.Fatal(LookupServiceURL + " must be set for vcda_vcenter_replication_manager acceptance tests")
	}
	if os.Getenv(SsoUser) == "" {
		t.Fatal(SsoUser + " must be set for vcda_vcenter_replication_manager acceptance tests")
	}
	if os.Getenv(SsoPassword) == "" {
		t.Fatal(SsoPassword + " must be set for vcda_vcenter_replication_manager acceptance tests")
	}
}

func testAccVcdaVcenterReplicationManagerConfigBasic() string {
	return fmt.Sprintf(`
data "vcda_service_cert" "manager_service_cert" {
  datacenter_id = %q
  name          = %q
  type          = "manager"
}

data "vcda_remote_services_thumbprint" "ls_thumbprint" {
  address      = %q
  port         = "443"
}

resource "vcda_vcenter_replication_manager" "manager_site" {
  service_cert              = data.vcda_service_cert.manager_service_cert.service_cert
  lookup_service_thumbprint = data.vcda_remote_services_thumbprint.ls_thumbprint.id

  license_key        = %q
  site_name          = "psvet-manager-site1"
  lookup_service_url = %q
  sso_user           = %q
  sso_password       = %q
}
`,
		os.Getenv(DatacenterID),
		os.Getenv(ManagerVmName),
		os.Getenv(LookupServiceAddress),
		os.Getenv(LicenseKey),
		os.Getenv(LookupServiceURL),
		os.Getenv(SsoUser),
		os.Getenv(SsoPassword),
	)
}
