/* Copyright 2023 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package vcda

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"
	"testing"
)

func TestAccVcdaDataSourceServiceCert_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccVcdaServiceCertConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.vcda_service_cert.service_cert", "service_cert"),
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
		os.Getenv(CloudVmName),
	)
}
