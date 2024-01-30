// Copyright (c) 2024 Broadcom. All Rights Reserved.
// Broadcom Confidential. The term "Broadcom" refers to Broadcom Inc.
// and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vcda

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"time"
)

func dataSourceVcdaTunnelConnectivity() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVcdaTunnelConnectivityRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"service_cert": {
				Type:        schema.TypeString,
				Description: "The certificate of the Cloud Director/vCenter Replication Manager Service.",
				Required:    true,
			},
			"tunnel_id": {
				Type:        schema.TypeString,
				Description: "The tunnel service ID.",
				Required:    true,
			},
			// Computed
			"id": {
				Type:        schema.TypeString,
				Description: "The health info task ID of the Cloud Director/vCenter Replication Manager Service.",
				Computed:    true,
			},
			"tunnel_service": {
				Type:        schema.TypeMap,
				Description: "The ID, URL and certificate of the Tunnel Service.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"tunnel_service_error_code": {
				Type:        schema.TypeString,
				Description: "The tunnel service error code.",
				Computed:    true,
			},
			"tunnel_service_error_msg": {
				Type:        schema.TypeString,
				Description: "The tunnel service error message.",
				Computed:    true,
			},
			"tunnel_service_error_args": {
				Type:        schema.TypeList,
				Description: "The tunnel service error arguments.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"tunnel_service_error_stacktrace": {
				Type:        schema.TypeString,
				Description: "The tunnel service error stacktrace.",
				Computed:    true,
			},
		},
	}
}

func dataSourceVcdaTunnelConnectivityRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	serviceCert := d.Get("service_cert").(string)

	taskID, err := c.getCloudHealth(serviceCert)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(*taskID)

	err = retryHealthTask(c, d, serviceCert, taskID)

	if err != nil {
		return diag.FromErr(err)
	}

	return getTunnelConnectivityInfo(c, d)
}

func getTunnelConnectivityInfo(c *Client, d *schema.ResourceData) diag.Diagnostics {
	var diags diag.Diagnostics

	tunnelID := d.Get("tunnel_id").(string)

	health, err := getHealthTaskResult(c, d)
	if err != nil {
		return diag.FromErr(err)
	}

	tunnelConnectivity, ok := health["tunnelConnectivity"].([]interface{})
	if !ok {
		return diag.FromErr(fmt.Errorf("unexpected type for tunnelConnectivity"))
	}

	tunnel, err := findTunnel(tunnelID, tunnelConnectivity)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setTunnelServiceInfoData(d, tunnel); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func setTunnelServiceInfoData(d *schema.ResourceData, tunnel map[string]interface{}) error {
	tunMap, ok := tunnel["tunnelService"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("unexpected type for tunnelService")
	}

	if err := d.Set("tunnel_service", tunMap); err != nil {
		return fmt.Errorf("error setting tunnel_service field: %s", err)
	}

	if tunError, ok := tunnel["error"].(map[string]interface{}); ok {
		if err := setErrorData(d, tunError, "tunnel_service"); err != nil {
			return err
		}
	}

	return nil
}

func findTunnel(tunnelID string, tunnels []interface{}) (map[string]interface{}, error) {
	for _, tun := range tunnels {
		if tunnel, ok := tun.(map[string]interface{}); ok {
			if tunMap, ok := tunnel["tunnelService"].(map[string]interface{}); ok {
				if tunMap["id"] == tunnelID {
					return tunnel, nil
				}
			}
		}
	}

	return nil, fmt.Errorf("tunnel with ID: %s was not found", tunnelID)
}
