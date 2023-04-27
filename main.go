/* Copyright 2023 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"

	"terraform-provider-for-vmware-cloud-director-availability/vcda"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() *schema.Provider {
			return vcda.Provider()
		},
	})
}
