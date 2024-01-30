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

func (at *AccTests) TestAccVcdaDataSourceCloudHealth_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccVcdaCloudHealthPreCheck(t)
		},
		ProviderFactories: testProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccVcdaCloudHealthConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.vcda_cloud_health.cloud_health", "id"),
					resource.TestCheckResourceAttrSet("data.vcda_cloud_health.cloud_health", "product_name"),
					resource.TestCheckResourceAttrSet("data.vcda_cloud_health.cloud_health", "build_version"),
					resource.TestCheckResourceAttrSet("data.vcda_cloud_health.cloud_health", "build_date"),
					resource.TestCheckResourceAttrSet("data.vcda_cloud_health.cloud_health", "instance_id"),
					resource.TestCheckResourceAttrSet("data.vcda_cloud_health.cloud_health", "runtime_id"),
					resource.TestCheckResourceAttrSet("data.vcda_cloud_health.cloud_health", "current_time"),
					resource.TestCheckResourceAttrSet("data.vcda_cloud_health.cloud_health", "address"),
					resource.TestCheckResourceAttrSet("data.vcda_cloud_health.cloud_health", "service_boot_timestamp"),
					resource.TestCheckResourceAttrSet("data.vcda_cloud_health.cloud_health", "appliance_boot_timestamp"),
					resource.TestCheckResourceAttrSet("data.vcda_cloud_health.cloud_health", "disk_usage.free"),
					resource.TestCheckResourceAttrSet("data.vcda_cloud_health.cloud_health", "disk_usage.usable"),
					resource.TestCheckResourceAttrSet("data.vcda_cloud_health.cloud_health", "disk_usage.total"),
					resource.TestCheckResourceAttrSet("data.vcda_cloud_health.cloud_health", "tunnels_ids.#"),
					resource.TestCheckResourceAttrSet("data.vcda_cloud_health.cloud_health", "manager_id"),
				),
			},
		},
	})
}

func testAccVcdaCloudHealthPreCheck(t *testing.T) {
	if v := os.Getenv(CloudVMName); v == "" {
		t.Fatal(CloudVMName + " must be set for vcda_cloud_health acceptance tests")
	}
	if v := os.Getenv(DatacenterID); v == "" {
		t.Fatal(DatacenterID + " must be set for vcda_cloud_health acceptance tests")
	}
}

func testAccVcdaCloudHealthConfigBasic() string {
	return fmt.Sprintf(`
data "vcda_service_cert" "cloud_service_cert" {
  datacenter_id = %q
  name          = %q
  type          = "cloud"
}

data "vcda_cloud_health" "cloud_health" {
  service_cert   = data.vcda_service_cert.cloud_service_cert.id
}
`,
		os.Getenv(DatacenterID),
		os.Getenv(CloudVMName),
	)
}
