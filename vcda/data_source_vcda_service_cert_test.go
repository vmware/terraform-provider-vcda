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

func (at *AccTests) TestAccVcdaDataSourceServiceCert_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccVcdaServiceCertConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.vcda_service_cert.service_cert", "id"),
				),
			},
		},
	})
}

func testAccVcdaServiceCertConfigBasic() string {
	return fmt.Sprintf(`
data "vcda_service_cert" "service_cert" {
  datacenter_id = %q
  name          = %q
  type          = "cloud"
}`,
		os.Getenv(DatacenterID),
		os.Getenv(CloudVMName),
	)
}
