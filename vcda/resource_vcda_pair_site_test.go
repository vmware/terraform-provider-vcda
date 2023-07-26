/* Copyright 2023 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package vcda

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"
	"testing"
)

func (at *AccTests) TestAccVcdaPairSite_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccVcdaPairSitePreCheck(t)
		},
		ProviderFactories: testProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccVcdaAppliancePasswordConfigBasic(os.Getenv(RemoteManagerVMName), "manager", os.Getenv(RemoteManagerAddress)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vcda_appliance_password.appliance_password", "root_password_expired", "false"),
					resource.TestCheckResourceAttrSet("vcda_appliance_password.appliance_password", "seconds_until_expiration"),
				),
			},
			{
				// initial config second site
				Config: testAccVcdaVcenterReplicationManagerConfigBasic(os.Getenv(RemoteManagerVMName), "manager-site2"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vcda_vcenter_replication_manager.manager_site", "is_licensed", "true"),
					resource.TestCheckResourceAttr("vcda_vcenter_replication_manager.manager_site", "expiration_date", "0"),
					resource.TestCheckResourceAttr("vcda_vcenter_replication_manager.manager_site",
						"ls_url", "https://"+os.Getenv(LookupServiceAddress)+":443/lookupservice/sdk"),
					resource.TestCheckResourceAttr("vcda_vcenter_replication_manager.manager_site", "vsphere_plugin_status", "OK"),

					resource.TestCheckResourceAttrSet("vcda_vcenter_replication_manager.manager_site", "ls_thumbprint"),
				),
			},
			{
				// pair sites
				Config: testAccVcdaPairSiteConfigBasic(os.Getenv(RemoteManagerVMName), os.Getenv(ManagerAddress)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("vcda_pair_site.pair_site", "site_id"),
					resource.TestCheckResourceAttrSet("vcda_pair_site.pair_site", "site_description"),
					resource.TestCheckResourceAttrSet("vcda_pair_site.pair_site", "api_public_url"),
					resource.TestCheckResourceAttrSet("vcda_pair_site.pair_site", "is_provider_deployment"),
				),
			},
		},
	})
}

func testAccVcdaPairSitePreCheck(t *testing.T) {
	if v := os.Getenv(RemoteManagerAddress); v == "" {
		t.Fatal(RemoteManagerAddress + " must be set for vcda_pair_site acceptance tests")
	}
	err := os.Setenv(VcdaIP, os.Getenv(RemoteManagerAddress))
	if err != nil {
		t.Fatal("error setting" + VcdaIP + " to " + RemoteManagerAddress + " for vcda_pair_site acceptance tests")
	}
	if v := os.Getenv(RemoteManagerVMName); v == "" {
		t.Fatal(RemoteManagerVMName + " must be set for vcda_pair_site acceptance tests")
	}
	if v := os.Getenv(ManagerAddress); v == "" {
		t.Fatal(ManagerAddress + " must be set for vcda_pair_site acceptance tests")
	}
	if v := os.Getenv(ManagerVMName); v == "" {
		t.Fatal(ManagerVMName + " must be set for vcda_pair_site acceptance tests")
	}
	if os.Getenv(LookupServiceAddress) == "" {
		t.Fatal(LookupServiceAddress + " must be set for vcda_pair_site acceptance tests")
	}
	if os.Getenv(LicenseKey) == "" {
		t.Fatal(LicenseKey + " must be set for vcda_pair_site acceptance tests")
	}
	if os.Getenv(SsoUser) == "" {
		t.Fatal(SsoUser + " must be set for vcda_pair_site acceptance tests")
	}
	if os.Getenv(SsoPassword) == "" {
		t.Fatal(SsoPassword + " must be set for vcda_pair_site acceptance tests")
	}
}

func testAccVcdaPairSiteConfigBasic(vmName string, remoteAddress string) string {
	return fmt.Sprintf(`
data "vcda_service_cert" "manager_service_cert" {
  datacenter_id = %q
  name          = %q
  type          = "manager"
}

# remote cloud site thumbprint
data "vcda_remote_services_thumbprint" "remote_cloud_thumbprint" {
  address = %q
  port    = "443"
}

resource "vcda_pair_site" "pair_site" {
  service_cert   = data.vcda_service_cert.manager_service_cert.id
  api_thumbprint = data.vcda_remote_services_thumbprint.remote_cloud_thumbprint.id

  api_url     = %q
  pairing_description = "pair site2"
}
`,
		os.Getenv(DatacenterID),
		vmName,
		remoteAddress,
		"https://"+remoteAddress+":8048",
	)
}
