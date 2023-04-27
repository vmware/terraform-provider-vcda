/* Copyright 2023 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package vcda

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"
	"testing"
)

func TestAccVcdaTunnel_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccVcdaTunnelPreCheck(t)
		},
		ProviderFactories: testProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccVcdaTunnelConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vcda_tunnel.add_tunnel", "tunnel_url", os.Getenv(TunnelURL)),
					resource.TestCheckResourceAttrSet("vcda_tunnel.add_tunnel", "tunnel_certificate"),
				),
			},
		},
	})
}

func testAccVcdaTunnelPreCheck(t *testing.T) {
	if v := os.Getenv(CloudVmName); v == "" {
		t.Fatal(CloudVmName + " must be set for vcda_tunnel acceptance tests")
	}
	if os.Getenv(TunnelVmName) == "" {
		t.Fatal(TunnelVmName + " must be set for vcda_tunnel acceptance tests")
	}
	if os.Getenv(TunnelURL) == "" {
		t.Fatal(TunnelURL + " must be set for vcda_tunnel acceptance tests")
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
  service_cert = data.vcda_service_cert.cloud_service_cert.service_cert

  url           = %q
  root_password = %q
  certificate   = data.vcda_service_cert.tunnel_service_cert.service_cert
}
`,
		os.Getenv(DatacenterID),
		os.Getenv(CloudVmName),
		os.Getenv(TunnelVmName),
		os.Getenv(TunnelURL),
		os.Getenv(RootPassword),
	)
}
