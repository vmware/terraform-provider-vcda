// Copyright (c) 2023-2024 Broadcom. All Rights Reserved.
// Broadcom Confidential. The term "Broadcom" refers to Broadcom Inc.
// and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vcda

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceVcdaTunnel() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVcdaTunnelCreate,
		ReadContext:   resourceVcdaTunnelRead,
		UpdateContext: resourceVcdaTunnelUpdate,
		DeleteContext: resourceVcdaTunnelDelete,
		Schema: map[string]*schema.Schema{
			"service_cert": {
				Type: schema.TypeString,
				Description: "The service certificate of the Cloud Director Replication Management Service " +
					"to which the Tunnel Service is being added.",
				Required: true,
			},
			"url": {
				Type:        schema.TypeString,
				Description: "The URL of the Tunnel Service.",
				Required:    true,
			},
			"certificate": {
				Type:        schema.TypeString,
				Description: "The certificate of the Tunnel Service.",
				Required:    true,
			},
			"root_password": {
				Type:        schema.TypeString,
				Description: "The **root** user password of the Tunnel Appliance.",
				Required:    true,
			},

			// computed
			"tunnel_url": {
				Type:        schema.TypeString,
				Description: "The URL of the Tunnel Service.",
				Computed:    true,
			},
			"tunnel_certificate": {
				Type:        schema.TypeString,
				Description: "The certificate of the Tunnel Service.",
				Computed:    true,
			},
		},
	}
}

func resourceVcdaTunnelCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	serviceCert := d.Get("service_cert").(string)

	URL := d.Get("url").(string)
	certificate := d.Get("certificate").(string)
	rootPassword := d.Get("root_password").(string)

	tunnelConfig, err := c.setTunnel(URL, certificate, rootPassword, serviceCert)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(tunnelConfig.ID)

	return resourceVcdaTunnelRead(ctx, d, m)
}

func resourceVcdaTunnelRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*Client)

	serviceCert := d.Get("service_cert").(string)

	tunnel, err := c.getTunnelConfig(serviceCert, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setTunnelData(d, tunnel); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceVcdaTunnelUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*Client)

	serviceCert := d.Get("service_cert").(string)

	if d.HasChange("url") || d.HasChange("root_password") {
		URL := d.Get("url").(string)
		certificate := d.Get("certificate").(string)
		rootPassword := d.Get("root_password").(string)

		tunnelConfig, err := c.setTunnel(URL, certificate, rootPassword, serviceCert)
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(tunnelConfig.ID)

		return resourceVcdaTunnelRead(ctx, d, m)
	}
	return diags
}

func resourceVcdaTunnelDelete(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	d.SetId("")

	return diags
}

func setTunnelData(d *schema.ResourceData, tunnelConfig *TunnelConfig) error {
	if err := d.Set("tunnel_url", tunnelConfig.URL); err != nil {
		return fmt.Errorf("error setting tunnel_url field: %s", err)
	}

	if err := d.Set("tunnel_certificate", tunnelConfig.Certificate); err != nil {
		return fmt.Errorf("error setting tunnel_certificate field: %s", err)
	}

	return nil
}
