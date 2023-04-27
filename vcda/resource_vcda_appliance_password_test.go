/* Copyright 2023 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package vcda

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"
	"testing"
)

func TestAccVcdaAppliancePassword_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccVcdaAppliancePasswordPreCheck(t)
		},
		ProviderFactories: testProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccVcdaAppliancePasswordConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vcda_appliance_password.appliance_password", "root_password_expired", "false"),
					resource.TestCheckResourceAttrSet("vcda_appliance_password.appliance_password", "seconds_until_expiration"),
				),
			},
		},
	})
}

func testAccVcdaAppliancePasswordPreCheck(t *testing.T) {
	if v := os.Getenv(CloudVmName); v == "" {
		t.Fatal(CloudVmName + " must be set for vcda_appliance_password acceptance tests")
	}
	if os.Getenv(RootPassword) == "" {
		t.Fatal(RootPassword + " must be set for vcda_appliance_password acceptance tests")
	}
	if os.Getenv(NewPassword) == "" {
		t.Fatal(NewPassword + " must be set for vcda_appliance_password acceptance tests")
	}
}

func testAccVcdaAppliancePasswordConfigBasic() string {
	return fmt.Sprintf(`
data "vcda_service_cert" "service_cert" {
  datacenter_id = %q
  name          = %q
  type          = "cloud"
}

resource "vcda_appliance_password" "appliance_password" {
  current_password = %q
  new_password     = %q
  appliance_ip     = %q
  service_cert     = data.vcda_service_cert.service_cert.service_cert
}
`,
		os.Getenv(DatacenterID),
		os.Getenv(CloudVmName),
		os.Getenv(RootPassword),
		os.Getenv(NewPassword),
		os.Getenv(VcdaIP),
	)
}
