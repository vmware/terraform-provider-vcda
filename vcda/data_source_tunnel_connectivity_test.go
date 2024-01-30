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

func (at *AccTests) TestAccVcdaDataSourceTunnelConnectivity_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccVcdaManagerHealthPreCheck(t)
		},
		ProviderFactories: testProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccVcdaTunnelConnectivityConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.vcda_tunnel_connectivity.tunnel_connectivity", "tunnel_service.id"),
					resource.TestCheckResourceAttrSet("data.vcda_tunnel_connectivity.tunnel_connectivity", "tunnel_service.url"),
					resource.TestCheckResourceAttrSet("data.vcda_tunnel_connectivity.tunnel_connectivity", "tunnel_service.certificate"),
				),
			},
		},
	})
}

func testAccVcdaTunnelConnectivityConfigBasic() string {
	return fmt.Sprintf(`
data "vcda_service_cert" "manager_service_cert" {
  datacenter_id = %q
  name          = %q
  type          = "manager"
}

data "vcda_manager_health" "manager_health" {
  service_cert   = data.vcda_service_cert.manager_service_cert.id
}

data "vcda_tunnel_connectivity" "tunnel_connectivity" {
  service_cert = data.vcda_service_cert.manager_service_cert.id
  tunnel_id = data.vcda_manager_health.manager_health.tunnels_ids[0]
}

`,
		os.Getenv(DatacenterID),
		os.Getenv(ManagerVMName),
	)
}
