// Copyright (c) 2023-2024 Broadcom. All Rights Reserved.
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

func (at *AccTests) TestAccVcdaAppliancePassword_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccVcdaAppliancePasswordPreCheck(t)
		},
		ProviderFactories: testProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccVcdaAppliancePasswordConfigBasic(os.Getenv(CloudVMName), "cloud", os.Getenv(VcdaIP)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vcda_appliance_password.appliance_password", "root_password_expired", "false"),
					resource.TestCheckResourceAttrSet("vcda_appliance_password.appliance_password", "seconds_until_expiration"),
				),
			},
		},
	})
}

func testAccVcdaAppliancePasswordPreCheck(t *testing.T) {
	if os.Getenv(CurrentPassword) == "" {
		t.Fatal(CurrentPassword + " must be set for vcda_appliance_password acceptance tests")
	}
	if os.Getenv(NewPassword) == "" {
		t.Fatal(NewPassword + " must be set for vcda_appliance_password acceptance tests")
	}
}

func testAccVcdaAppliancePasswordConfigBasic(vmName string, applianceType string, applianceIP string) string {
	return fmt.Sprintf(`
data "vcda_service_cert" "service_cert" {
  datacenter_id = %q
  name          = %q
  type          = %q
}

resource "vcda_appliance_password" "appliance_password" {
  current_password = %q
  new_password     = %q
  appliance_ip     = %q
  service_cert     = data.vcda_service_cert.service_cert.id
}
`,
		os.Getenv(DatacenterID),
		vmName,
		applianceType,
		os.Getenv(CurrentPassword),
		os.Getenv(NewPassword),
		applianceIP,
	)
}
