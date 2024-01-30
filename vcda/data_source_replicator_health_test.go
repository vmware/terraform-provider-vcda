// Copyright (c) 2024 Broadcom. All Rights Reserved.
// Broadcom Confidential. The term "Broadcom" refers to Broadcom Inc.
// and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vcda

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"
	"testing"
)

func (at *AccTests) TestAccVcdaDataSourceReplicatorHealth_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccVcdaReplicatorPreCheck(t)
			testAccVcdaReplicatorHealthPreCheck(t)
		},
		ProviderFactories: testProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccVcdaReplicatorHealthConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.vcda_replicator_health.replicator_health", "id"),
					resource.TestCheckResourceAttrSet("data.vcda_replicator_health.replicator_health", "product_name"),
					resource.TestCheckResourceAttrSet("data.vcda_replicator_health.replicator_health", "build_version"),
					resource.TestCheckResourceAttrSet("data.vcda_replicator_health.replicator_health", "build_date"),
					resource.TestCheckResourceAttrSet("data.vcda_replicator_health.replicator_health", "instance_id"),
					resource.TestCheckResourceAttrSet("data.vcda_replicator_health.replicator_health", "runtime_id"),
					resource.TestCheckResourceAttrSet("data.vcda_replicator_health.replicator_health", "current_time"),
					resource.TestCheckResourceAttrSet("data.vcda_replicator_health.replicator_health", "address"),
					resource.TestCheckResourceAttrSet("data.vcda_replicator_health.replicator_health", "service_boot_timestamp"),
					resource.TestCheckResourceAttrSet("data.vcda_replicator_health.replicator_health", "appliance_boot_timestamp"),
					resource.TestCheckResourceAttrSet("data.vcda_replicator_health.replicator_health", "disk_usage.free"),
					resource.TestCheckResourceAttrSet("data.vcda_replicator_health.replicator_health", "disk_usage.usable"),
					resource.TestCheckResourceAttrSet("data.vcda_replicator_health.replicator_health", "disk_usage.total"),
					resource.TestCheckResourceAttrSet("data.vcda_replicator_health.replicator_health", "offline_managers_ids.#"),
					resource.TestCheckResourceAttrSet("data.vcda_replicator_health.replicator_health", "online_managers_ids.#"),
				),
			},
		},
	})
}

func testAccVcdaReplicatorHealthPreCheck(t *testing.T) {
	if v := os.Getenv(ManagerAddress); v == "" {
		t.Fatal(ManagerAddress + " must be set for vcda_replicator_health acceptance tests")
	}

	err := os.Setenv(VcdaIP, os.Getenv(ManagerAddress))
	if err != nil {
		t.Fatal("error setting" + VcdaIP + " to " + ManagerAddress + " for vcda_replicator_health acceptance tests")
	}

	if v := os.Getenv(DatacenterID); v == "" {
		t.Fatal(DatacenterID + " must be set for vcda_replicator_health acceptance tests")
	}
}

func testAccVcdaReplicatorHealthConfigBasic() string {
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
  service_cert              = data.vcda_service_cert.manager_service_cert.id
  lookup_service_thumbprint = data.vcda_remote_services_thumbprint.ls_thumbprint.id
}

data "vcda_manager_health" "manager_health" {
  depends_on = [vcda_replicator.add_replicator]
  service_cert   = data.vcda_service_cert.manager_service_cert.id
}

data "vcda_replicator_health" "replicator_health" {
  depends_on = [data.vcda_manager_health.manager_health]
  service_cert   = data.vcda_service_cert.manager_service_cert.id
  replicator_id = data.vcda_manager_health.manager_health.local_replicators_ids[0]
}
`,
		os.Getenv(DatacenterID),
		os.Getenv(ManagerVMName),
		os.Getenv(LookupServiceAddress),
		os.Getenv(ReplicatorAddress),
		"https://"+os.Getenv(LookupServiceAddress)+":443/lookupservice/sdk",
		"https://"+os.Getenv(ReplicatorAddress)+":8043",
		os.Getenv(SsoUser),
		os.Getenv(SsoPassword),
		os.Getenv(RootPassword),
	)
}
