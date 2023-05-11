/* Copyright 2023 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package vcda

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"
	"testing"
)

func (at *AccTests) TestAccVcdaCloudDirectorReplicationManager_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccVcdaCloudDirectorReplicationManagerPreCheck(t)
		},
		ProviderFactories: testProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccVcdaCloudDirectorReplicationManagerConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vcda_cloud_director_replication_manager.cloud_site", "is_licensed", "true"),
					resource.TestCheckResourceAttr("vcda_cloud_director_replication_manager.cloud_site", "expiration_date", "0"),
					resource.TestCheckResourceAttr("vcda_cloud_director_replication_manager.cloud_site", "is_combined", "false"),
					resource.TestCheckResourceAttr("vcda_cloud_director_replication_manager.cloud_site",
						"ls_url", "https://"+os.Getenv(LookupServiceAddress)+":443/lookupservice/sdk"),
					resource.TestCheckResourceAttr("vcda_cloud_director_replication_manager.cloud_site",
						"vcloud_url", "https://"+os.Getenv(VcloudDirectorAddress)+"/api"),

					resource.TestCheckResourceAttrSet("vcda_cloud_director_replication_manager.cloud_site", "ls_thumbprint"),
					resource.TestCheckResourceAttrSet("vcda_cloud_director_replication_manager.cloud_site", "local_site"),
					resource.TestCheckResourceAttrSet("vcda_cloud_director_replication_manager.cloud_site", "local_site_description"),
					resource.TestCheckResourceAttrSet("vcda_cloud_director_replication_manager.cloud_site", "vcloud_thumbprint"),
					resource.TestCheckResourceAttrSet("vcda_cloud_director_replication_manager.cloud_site", "vcloud_username"),
					resource.TestCheckResourceAttrSet("vcda_cloud_director_replication_manager.cloud_site", "mgmt_address"),
					resource.TestCheckResourceAttrSet("vcda_cloud_director_replication_manager.cloud_site", "mgmt_port"),
					resource.TestCheckResourceAttrSet("vcda_cloud_director_replication_manager.cloud_site", "mgmt_public_address"),
					resource.TestCheckResourceAttrSet("vcda_cloud_director_replication_manager.cloud_site", "mgmt_public_port"),
					resource.TestCheckResourceAttrSet("vcda_cloud_director_replication_manager.cloud_site", "api_address"),
					resource.TestCheckResourceAttrSet("vcda_cloud_director_replication_manager.cloud_site", "api_port"),
					resource.TestCheckResourceAttrSet("vcda_cloud_director_replication_manager.cloud_site", "api_public_address"),
					resource.TestCheckResourceAttrSet("vcda_cloud_director_replication_manager.cloud_site", "api_public_port"),
				),
			},
		},
	})
}

func testAccVcdaCloudDirectorReplicationManagerPreCheck(t *testing.T) {
	if v := os.Getenv(CloudVmName); v == "" {
		t.Fatal(CloudVmName + " must be set for vcda_cloud_director_replication_manager acceptance tests")
	}
	if os.Getenv(VcloudDirectorAddress) == "" {
		t.Fatal(VcloudDirectorAddress + " must be set for vcda_cloud_director_replication_manager acceptance tests")
	}
	if os.Getenv(LookupServiceAddress) == "" {
		t.Fatal(LookupServiceAddress + " must be set for vcda_cloud_director_replication_manager acceptance tests")
	}
	if os.Getenv(LicenseKey) == "" {
		t.Fatal(LicenseKey + " must be set for vcda_cloud_director_replication_manager acceptance tests")
	}
	if os.Getenv(VcloudDirectorUsername) == "" {
		t.Fatal(VcloudDirectorUsername + " must be set for vcda_cloud_director_replication_manager acceptance tests")
	}
	if os.Getenv(VcloudDirectorPassword) == "" {
		t.Fatal(VcloudDirectorPassword + " must be set for vcda_cloud_director_replication_manager acceptance tests")
	}
}

func testAccVcdaCloudDirectorReplicationManagerConfigBasic() string {
	return fmt.Sprintf(`
data "vcda_service_cert" "cloud_service_cert" {
  datacenter_id = %q
  name          = %q
  type          = "cloud"
}

data "vcda_remote_services_thumbprint" "vcd_thumbprint" {
  address      = %q
  port         = "443"
}

data "vcda_remote_services_thumbprint" "ls_thumbprint" {
  address      = %q
  port         = "443"
}

resource "vcda_cloud_director_replication_manager" "cloud_site" {
  service_cert              = data.vcda_service_cert.cloud_service_cert.id
  lookup_service_thumbprint = data.vcda_remote_services_thumbprint.ls_thumbprint.id
  vcd_thumbprint            = data.vcda_remote_services_thumbprint.vcd_thumbprint.id

  license_key      = %q
  site_name        = "cloud-site1"
  site_description = "cloud site"

  public_endpoint_address = "vcda.pub"
  public_endpoint_port    = 443

  vcd_username = %q
  vcd_password = %q
  vcd_url      = %q

  lookup_service_url = %q
}
`,
		os.Getenv(DatacenterID),
		os.Getenv(CloudVmName),
		os.Getenv(VcloudDirectorAddress),
		os.Getenv(LookupServiceAddress),
		os.Getenv(LicenseKey),
		os.Getenv(VcloudDirectorUsername),
		os.Getenv(VcloudDirectorPassword),
		"https://"+os.Getenv(VcloudDirectorAddress),
		"https://"+os.Getenv(LookupServiceAddress)+":443/lookupservice/sdk",
	)
}
