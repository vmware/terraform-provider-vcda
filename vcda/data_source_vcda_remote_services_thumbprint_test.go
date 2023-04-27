/* Copyright 2023 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package vcda

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"
	"testing"
)

func TestAccVcdaDataSourceRemoteServicesThumbprint_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccVcdaRemoteServicesThumbprintConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.vcda_remote_services_thumbprint.thumbprint", "id"),
				),
			},
		},
	})
}

func testAccVcdaRemoteServicesThumbprintConfigBasic() string {
	return fmt.Sprintf(`
data "vcda_remote_services_thumbprint" "thumbprint" {
  address      = %q
  port         = "443"
}
`,
		os.Getenv(VcdaIP),
	)
}
