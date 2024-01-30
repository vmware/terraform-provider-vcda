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

func (at *AccTests) TestAccVcdaDataSourceManagerHealth_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccVcdaManagerHealthPreCheck(t)
		},
		ProviderFactories: testProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccVcdaManagerHealthConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.vcda_manager_health.manager_health", "id"),
					resource.TestCheckResourceAttrSet("data.vcda_manager_health.manager_health", "product_name"),
					resource.TestCheckResourceAttrSet("data.vcda_manager_health.manager_health", "build_version"),
					resource.TestCheckResourceAttrSet("data.vcda_manager_health.manager_health", "build_date"),
					resource.TestCheckResourceAttrSet("data.vcda_manager_health.manager_health", "instance_id"),
					resource.TestCheckResourceAttrSet("data.vcda_manager_health.manager_health", "runtime_id"),
					resource.TestCheckResourceAttrSet("data.vcda_manager_health.manager_health", "current_time"),
					resource.TestCheckResourceAttrSet("data.vcda_manager_health.manager_health", "address"),
					resource.TestCheckResourceAttrSet("data.vcda_manager_health.manager_health", "service_boot_timestamp"),
					resource.TestCheckResourceAttrSet("data.vcda_manager_health.manager_health", "appliance_boot_timestamp"),
					resource.TestCheckResourceAttrSet("data.vcda_manager_health.manager_health", "disk_usage.free"),
					resource.TestCheckResourceAttrSet("data.vcda_manager_health.manager_health", "disk_usage.usable"),
					resource.TestCheckResourceAttrSet("data.vcda_manager_health.manager_health", "disk_usage.total"),
					resource.TestCheckResourceAttrSet("data.vcda_manager_health.manager_health", "offline_replicators_ids.#"),
					resource.TestCheckResourceAttrSet("data.vcda_manager_health.manager_health", "online_replicators_ids.#"),
					resource.TestCheckResourceAttrSet("data.vcda_manager_health.manager_health", "local_replicators_ids.#"),
					resource.TestCheckResourceAttrSet("data.vcda_manager_health.manager_health", "tunnels_ids.#"),
				),
			},
		},
	})
}

func testAccVcdaManagerHealthPreCheck(t *testing.T) {
	if v := os.Getenv(ManagerAddress); v == "" {
		t.Fatal(ManagerAddress + " must be set for vcda_manager_health acceptance tests")
	}

	err := os.Setenv(VcdaIP, os.Getenv(ManagerAddress))
	if err != nil {
		t.Fatal("error setting" + VcdaIP + " to " + ManagerAddress + " for vcda_manager_health acceptance tests")
	}

	if v := os.Getenv(ManagerVMName); v == "" {
		t.Fatal(ManagerVMName + " must be set for vcda_manager_health acceptance tests")
	}

	if v := os.Getenv(DatacenterID); v == "" {
		t.Fatal(DatacenterID + " must be set for vcda_manager_health acceptance tests")
	}
}

func testAccVcdaManagerHealthConfigBasic() string {
	return fmt.Sprintf(`
data "vcda_service_cert" "manager_service_cert" {
  datacenter_id = %q
  name          = %q
  type          = "manager"
}

data "vcda_manager_health" "manager_health" {
  service_cert   = data.vcda_service_cert.manager_service_cert.id
}

output "vcda_manager_health" {
  value = data.vcda_manager_health.manager_health.*
}
`,
		os.Getenv(DatacenterID),
		os.Getenv(ManagerVMName),
	)
}
