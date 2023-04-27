/* Copyright 2023 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package vcda

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceVcdaReplicator() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVcdaReplicatorCreate,
		ReadContext:   resourceVcdaReplicatorRead,
		UpdateContext: resourceVcdaReplicatorUpdate,
		DeleteContext: resourceVcdaReplicatorDelete,
		Schema: map[string]*schema.Schema{
			"service_cert": {
				Type:        schema.TypeString,
				Description: "Replicator appliance VM thumbprint.",
				Required:    true,
			},
			"lookup_service_url": {
				Type:        schema.TypeString,
				Description: "Lookup service URL.",
				Required:    true,
			},
			"lookup_service_thumbprint": {
				Type:        schema.TypeString,
				Description: "Lookup service thumbprint.",
				Required:    true,
			},
			"api_url": {
				Type:        schema.TypeString,
				Description: "Replicator service API URL.",
				Required:    true,
			},
			"api_thumbprint": {
				Type:        schema.TypeString,
				Description: "Replicator service API thumbprint.",
				Required:    true,
			},
			"sso_user": {
				Type:        schema.TypeString,
				Description: "Replicator service SSO user.",
				Required:    true,
			},
			"sso_password": {
				Type:        schema.TypeString,
				Description: "Replicator service SSO password.",
				Required:    true,
			},
			"root_password": {
				Type:        schema.TypeString,
				Description: "Replicator service root password.",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Replicator service description.",
				Optional:    true,
			},
			"owner": {
				Type:        schema.TypeString,
				Description: "Replicator service owner.",
				Required:    true,
			},
			"site_name": {
				Type:        schema.TypeString,
				Description: "Replicator service site name.",
				Required:    true,
			},

			// computed
			"is_in_maintenance_mode": {
				Type:        schema.TypeBool,
				Description: "Flag indicating whether replicator service is in maintenance mode.",
				Computed:    true,
			},
			"data_address": {
				Type:        schema.TypeString,
				Description: "Replicator service data address.",
				Computed:    true,
			},
			"build_version": {
				Type:        schema.TypeString,
				Description: "Replicator service build version.",
				Computed:    true,
			},
			"replicator_ls_url": {
				Type:        schema.TypeString,
				Description: "Replicator service lookup service URL.",
				Computed:    true,
			},
			"replicator_ls_thumbprint": {
				Type:        schema.TypeString,
				Description: "Replicator service lookup service thumbprint.",
				Computed:    true,
			},
		},
	}

}

func resourceVcdaReplicatorCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	serviceCert := d.Get("service_cert").(string)
	lsURL := d.Get("lookup_service_url").(string)
	lsThumbprint := d.Get("lookup_service_thumbprint").(string)
	apiURL := d.Get("api_url").(string)
	apiThumbprint := d.Get("api_thumbprint").(string)
	rootPassword := d.Get("root_password").(string)
	ssoUser := d.Get("sso_user").(string)
	ssoPassword := d.Get("sso_password").(string)
	description := d.Get("description").(string)
	owner := d.Get("owner").(string)
	siteName := d.Get("site_name").(string)
	host := c.VcdaIP + ":8441"

	// set replicator lookup service
	replicatorLookupService, err := c.setReplicatorLookupService(host, lsURL, lsThumbprint, apiURL, apiThumbprint, rootPassword, serviceCert)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := setReplicatorLookupServiceData(d, replicatorLookupService); err != nil {
		return diag.FromErr(err)
	}

	// add replicator
	details := ReplicatorConfigData{APIURL: apiURL, APIThumbprint: apiThumbprint, RootPassword: rootPassword, SsoUser: ssoUser, SsoPassword: ssoPassword}

	replicator, err := c.addReplicator(host, serviceCert, description, owner, siteName, details)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(replicator.ID)

	return resourceVcdaReplicatorRead(ctx, d, m)
}

func resourceVcdaReplicatorRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*Client)

	host := c.VcdaIP + ":8441"
	serviceCert := d.Get("service_cert").(string)

	replicator, err := c.getReplicator(host, serviceCert, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("is_in_maintenance_mode", replicator.IsInMaintenanceMode); err != nil {
		return diag.FromErr(fmt.Errorf("error setting is_in_maintenance_mode field: %s", err))
	}

	if err := d.Set("data_address", replicator.DataAddress); err != nil {
		return diag.FromErr(fmt.Errorf("error setting data_address field: %s", err))
	}

	if err := d.Set("build_version", replicator.BuildVersion); err != nil {
		return diag.FromErr(fmt.Errorf("error setting build_version field: %s", err))
	}

	return diags
}

func resourceVcdaReplicatorUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*Client)

	if d.HasChange("root_password") || d.HasChange("sso_user") || d.HasChange("sso_password") {
		rootPassword := d.Get("root_password").(string)
		ssoUser := d.Get("sso_user").(string)
		ssoPassword := d.Get("sso_password").(string)

		replicatorID := d.Id()
		apiURL := d.Get("api_url").(string)
		apiThumbprint := d.Get("api_thumbprint").(string)
		serviceCert := d.Get("service_cert").(string)
		host := c.VcdaIP + ":8441"

		if err := c.repairReplicator(host, serviceCert, replicatorID, apiURL, apiThumbprint, rootPassword, ssoUser, ssoPassword); err != nil {
			return diag.FromErr(err)
		}

		return resourceVcdaReplicatorRead(ctx, d, m)
	}
	return diags
}

func resourceVcdaReplicatorDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*Client)

	host := c.VcdaIP + ":8441"
	serviceCert := d.Get("service_cert").(string)
	replicatorID := d.Id()

	if err := c.deleteReplicator(host, serviceCert, replicatorID); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")

	return diags
}

func setReplicatorLookupServiceData(d *schema.ResourceData, lookupService *LookupService) error {
	if err := d.Set("replicator_ls_url", lookupService.LsURL); err != nil {
		return fmt.Errorf("error setting ls_url field: %s", err)
	}

	if err := d.Set("replicator_ls_thumbprint", lookupService.LsThumbprint); err != nil {
		return fmt.Errorf("error setting ls_thumbprint field: %s", err)
	}

	return nil
}
