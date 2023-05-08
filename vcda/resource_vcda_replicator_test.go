/* Copyright 2023 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package vcda

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"
	"testing"
)

func (at *AccTests) TestAccVcdaReplicator_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccVcdaAppliancePasswordPreCheck(t)
			testAccVcdaReplicatorPreCheck(t)
		},
		ProviderFactories: testProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccVcdaAppliancePasswordConfigBasic(os.Getenv(ReplicatorVmName), "replicator", os.Getenv(ReplicatorAddress)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vcda_appliance_password.appliance_password", "root_password_expired", "false"),
					resource.TestCheckResourceAttrSet("vcda_appliance_password.appliance_password", "seconds_until_expiration"),
				),
			},
			{
				Config: testAccVcdaReplicatorConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vcda_replicator.add_replicator", "is_in_maintenance_mode", "false"),
					resource.TestCheckResourceAttr("vcda_replicator.add_replicator",
						"replicator_ls_url", "https://"+os.Getenv(LookupServiceAddress)+":443/lookupservice/sdk"),

					resource.TestCheckResourceAttrSet("vcda_replicator.add_replicator", "data_address"),
					resource.TestCheckResourceAttrSet("vcda_replicator.add_replicator", "build_version"),
					resource.TestCheckResourceAttrSet("vcda_replicator.add_replicator", "replicator_ls_thumbprint"),
				),
			},
		},
	})
}

func testAccVcdaReplicatorPreCheck(t *testing.T) {
	if v := os.Getenv(ManagerAddress); v == "" {
		t.Fatal(ManagerAddress + " must be set for vcda_replicator acceptance tests")
	}
	err := os.Setenv(VcdaIP, os.Getenv(ManagerAddress))
	if err != nil {
		t.Fatal("error setting" + VcdaIP + " to " + ManagerAddress + " for vcda_replicator acceptance tests")
	}
	if v := os.Getenv(ReplicatorVmName); v == "" {
		t.Fatal(ReplicatorVmName + " must be set for vcda_replicator acceptance tests")
	}
	if v := os.Getenv(ManagerVmName); v == "" {
		t.Fatal(ManagerVmName + " must be set for vcda_replicator acceptance tests")
	}
	if os.Getenv(LookupServiceAddress) == "" {
		t.Fatal(LookupServiceAddress + " must be set for vcda_replicator acceptance tests")
	}
	if os.Getenv(ReplicatorAddress) == "" {
		t.Fatal(ReplicatorAddress + " must be set for vcda_replicator acceptance tests")
	}
	if os.Getenv(SsoUser) == "" {
		t.Fatal(SsoUser + " must be set for vcda_replicator acceptance tests")
	}
	if os.Getenv(SsoPassword) == "" {
		t.Fatal(SsoPassword + " must be set for vcda_replicator acceptance tests")
	}
	if os.Getenv(RootPassword) == "" {
		t.Fatal(RootPassword + " must be set for vcda_replicator acceptance tests")
	}
}

func testAccVcdaReplicatorConfigBasic() string {
	return fmt.Sprintf(`

variable "datacenter_id" {
  type    = string
  default = %q
}

data "vcda_service_cert" "manager_service_cert" {
  datacenter_id = var.datacenter_id
  name          = %q
  type          = "manager"
}

data "vcda_remote_services_thumbprint" "ls_thumbprint" {
  address      = %q
  port         = "443"
}

data "vcda_remote_services_thumbprint" "replicator_thumbprint" {
  address      = %q
  port         = "443"
}

resource "vcda_replicator" "add_replicator" {
  lookup_service_url = %q
  api_url            = %q
  sso_user           = %q
  sso_password       = %q
  root_password      = %q
  owner              = "*"
  site_name          = "manager-site1"

  api_thumbprint            = data.vcda_remote_services_thumbprint.replicator_thumbprint.id
  service_cert              = data.vcda_service_cert.manager_service_cert.service_cert
  lookup_service_thumbprint = data.vcda_remote_services_thumbprint.ls_thumbprint.id
}
`,
		os.Getenv(DatacenterID),
		os.Getenv(ManagerVmName),
		os.Getenv(LookupServiceAddress),
		os.Getenv(ReplicatorAddress),
		"https://"+os.Getenv(LookupServiceAddress)+":443/lookupservice/sdk",
		"https://"+os.Getenv(ReplicatorAddress)+":8043",
		os.Getenv(SsoUser),
		os.Getenv(SsoPassword),
		os.Getenv(RootPassword),
	)
}
