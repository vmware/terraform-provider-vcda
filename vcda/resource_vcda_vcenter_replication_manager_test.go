/* Copyright 2023 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package vcda

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"
	"testing"
)

func (at *AccTests) TestAccVcdaVcenterReplicationManager_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccVcdaAppliancePasswordPreCheck(t)
			testAccVcdaVcenterReplicationManagerPreCheck(t)
		},
		ProviderFactories: testProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccVcdaAppliancePasswordConfigBasic(os.Getenv(ManagerVMName), "manager", os.Getenv(ManagerAddress)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vcda_appliance_password.appliance_password", "root_password_expired", "false"),
					resource.TestCheckResourceAttrSet("vcda_appliance_password.appliance_password", "seconds_until_expiration"),
				),
			},
			{
				Config: testAccVcdaVcenterReplicationManagerConfigBasic(os.Getenv(ManagerVMName), "manager-site1"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vcda_vcenter_replication_manager.manager_site", "is_licensed", "true"),
					resource.TestCheckResourceAttr("vcda_vcenter_replication_manager.manager_site", "expiration_date", "0"),
					resource.TestCheckResourceAttr("vcda_vcenter_replication_manager.manager_site",
						"ls_url", "https://"+os.Getenv(LookupServiceAddress)+":443/lookupservice/sdk"),
					resource.TestCheckResourceAttr("vcda_vcenter_replication_manager.manager_site", "vsphere_plugin_status", "OK"),

					resource.TestCheckResourceAttrSet("vcda_vcenter_replication_manager.manager_site", "ls_thumbprint"),
				),
			},
		},
	})
}

func testAccVcdaVcenterReplicationManagerPreCheck(t *testing.T) {
	if v := os.Getenv(ManagerAddress); v == "" {
		t.Fatal(ManagerAddress + " must be set for vcda_vcenter_replication_manager acceptance tests")
	}
	err := os.Setenv(VcdaIP, os.Getenv(ManagerAddress))
	if err != nil {
		t.Fatal("error setting" + VcdaIP + " to " + ManagerAddress + " for vcda_vcenter_replication_manager acceptance tests")
	}
	if v := os.Getenv(ManagerVMName); v == "" {
		t.Fatal(ManagerVMName + " must be set for vcda_vcenter_replication_manager acceptance tests")
	}
	if os.Getenv(LookupServiceAddress) == "" {
		t.Fatal(LookupServiceAddress + " must be set for vcda_vcenter_replication_manager acceptance tests")
	}
	if os.Getenv(LicenseKey) == "" {
		t.Fatal(LicenseKey + " must be set for vcda_vcenter_replication_manager acceptance tests")
	}
	if os.Getenv(SsoUser) == "" {
		t.Fatal(SsoUser + " must be set for vcda_vcenter_replication_manager acceptance tests")
	}
	if os.Getenv(SsoPassword) == "" {
		t.Fatal(SsoPassword + " must be set for vcda_vcenter_replication_manager acceptance tests")
	}
}

func testAccVcdaVcenterReplicationManagerConfigBasic(managerVMName string, siteName string) string {
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
  service_cert              = data.vcda_service_cert.manager_service_cert.id
  lookup_service_thumbprint = data.vcda_remote_services_thumbprint.ls_thumbprint.id

  license_key        = %q
  site_name          = %q
  lookup_service_url = %q
  sso_user           = %q
  sso_password       = %q
}
`,
		os.Getenv(DatacenterID),
		managerVMName,
		os.Getenv(LookupServiceAddress),
		os.Getenv(LicenseKey),
		siteName,
		"https://"+os.Getenv(LookupServiceAddress)+":443/lookupservice/sdk",
		os.Getenv(SsoUser),
		os.Getenv(SsoPassword),
	)
}
