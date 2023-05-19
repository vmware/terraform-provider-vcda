/* Copyright 2023 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package vcda

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"
	"testing"
)

func (at *AccTests) TestAccVcdaTunnel_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccVcdaAppliancePasswordPreCheck(t)
			testAccVcdaTunnelPreCheck(t)
		},
		ProviderFactories: testProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccVcdaAppliancePasswordConfigBasic(os.Getenv(TunnelVMName), "tunnel", os.Getenv(TunnelAddress)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vcda_appliance_password.appliance_password", "root_password_expired", "false"),
					resource.TestCheckResourceAttrSet("vcda_appliance_password.appliance_password", "seconds_until_expiration"),
				),
			},
			{
				Config: testAccVcdaTunnelConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vcda_tunnel.add_tunnel", "tunnel_url", "https://"+os.Getenv(TunnelAddress)+":8047"),
					resource.TestCheckResourceAttrSet("vcda_tunnel.add_tunnel", "tunnel_certificate"),
				),
			},
		},
	})
}

func testAccVcdaTunnelPreCheck(t *testing.T) {
	if v := os.Getenv(CloudVMName); v == "" {
		t.Fatal(CloudVMName + " must be set for vcda_tunnel acceptance tests")
	}
	if os.Getenv(TunnelVMName) == "" {
		t.Fatal(TunnelVMName + " must be set for vcda_tunnel acceptance tests")
	}
	if os.Getenv(TunnelAddress) == "" {
		t.Fatal(TunnelAddress + " must be set for vcda_tunnel acceptance tests")
	}
	if os.Getenv(RootPassword) == "" {
		t.Fatal(RootPassword + " must be set for vcda_tunnel acceptance tests")
	}
}

func testAccVcdaTunnelConfigBasic() string {
	return fmt.Sprintf(`

variable "datacenter_id" {
  type    = string
  default = %q
}

data "vcda_service_cert" "cloud_service_cert" {
  datacenter_id = var.datacenter_id
  name          = %q
  type          = "cloud"
}

data "vcda_service_cert" "tunnel_service_cert" {
  datacenter_id = var.datacenter_id
  name          = %q
  type          = "tunnel"
}

resource "vcda_tunnel" "add_tunnel" {
  service_cert = data.vcda_service_cert.cloud_service_cert.id

  url           = %q
  root_password = %q
  certificate   = data.vcda_service_cert.tunnel_service_cert.id
}
`,
		os.Getenv(DatacenterID),
		os.Getenv(CloudVMName),
		os.Getenv(TunnelVMName),
		"https://"+os.Getenv(TunnelAddress)+":8047",
		os.Getenv(RootPassword),
	)
}
