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

func resourceVcdaVcenterReplicationManager() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVcenterReplicationManagerCreate,
		ReadContext:   resourceVcenterReplicationManagerRead,
		UpdateContext: resourceVcenterReplicationManagerUpdate,
		DeleteContext: resourceVcenterReplicationManagerDelete,
		Schema: map[string]*schema.Schema{
			"service_cert": {
				Type:        schema.TypeString,
				Description: "The service certificate of the vCenter Replication Manager.",
				Required:    true,
			},
			"lookup_service_thumbprint": {
				Type: schema.TypeString,
				Description: "The thumbprint of the vCenter Server Lookup service. It can either be computed from " +
					"the `vcda_remote_services_thumbprint` data source or provided directly as a SHA-256 fingerprint.",
				Required: true,
			},
			"license_key": {
				Type:        schema.TypeString,
				Description: "The license key of VMware Cloud Director Availability.",
				Required:    true,
			},
			"site_name": {
				Type:        schema.TypeString,
				Description: "The site name of the vCenter Replication Manager.",
				Required:    true,
			},
			"lookup_service_url": {
				Type: schema.TypeString,
				Description: "The URL of the vCenter Server Lookup service. " +
					"For example, https://server.domain.com/lookupservice/sdk.",
				Required: true,
			},
			"sso_user": {
				Type:        schema.TypeString,
				Description: "The user name of a single sign-on (SSO) administrator.",
				Required:    true,
			},
			"sso_password": {
				Type:        schema.TypeString,
				Description: "The password of the SSO administrator.",
				Required:    true,
			},

			// computed:
			"is_licensed": {
				Type:        schema.TypeBool,
				Description: "Flag indicating whether the service is licensed.",
				Computed:    true,
			},
			"expiration_date": {
				Type:        schema.TypeInt,
				Description: "The expiration date of the license.",
				Computed:    true,
			},
			"site": {
				Type:        schema.TypeString,
				Description: "The site name of the vCenter Replication Manager.",
				Computed:    true,
			},
			"ls_url": {
				Type:        schema.TypeString,
				Description: "The URL of the vCenter Server Lookup service.",
				Computed:    true,
			},
			"ls_thumbprint": {
				Type:        schema.TypeString,
				Description: "The thumbprint of the vCenter Server Lookup service.",
				Computed:    true,
			},
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
			"vsphere_plugin_status": {
				Type:        schema.TypeString,
				Description: "The status of the Vsphere Plugin.",
				Computed:    true,
			},
		},
	}

}

func resourceVcenterReplicationManagerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	serviceCert := d.Get("service_cert").(string)
	licenseKey := d.Get("license_key").(string)
	siteName := d.Get("site_name").(string)
	lsURL := d.Get("lookup_service_url").(string)
	lsThumbprint := d.Get("lookup_service_thumbprint").(string)
	ssoUser := d.Get("sso_user").(string)
	ssoPassword := d.Get("sso_password").(string)

	// set license
	license, err := c.setLicense(serviceCert, licenseKey)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setLicenseData(d, license); err != nil {
		return diag.FromErr(err)
	}

	// set site name
	site, err := c.setSiteName(siteName, serviceCert)
	if err != nil {
		return diag.FromErr(err)
	}

	// set manager lookup service
	if err := c.setManagerLookupService(lsURL, lsThumbprint, ssoUser, ssoPassword, serviceCert); err != nil {
		return diag.FromErr(err)
	}

	pluginStatus, err := c.setVspherePlugin(serviceCert)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("vsphere_plugin_status", pluginStatus.Status); err != nil {
		return diag.Errorf("error setting vsphere_plugin_status field: %s", err)
	}

	d.SetId(site.ID)

	return resourceVcenterReplicationManagerRead(ctx, d, m)
}

func resourceVcenterReplicationManagerRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	serviceCert := d.Get("service_cert").(string)

	managerSite, err := c.getManagerSiteConfig(serviceCert)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setSiteData(d, managerSite); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceVcenterReplicationManagerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	serviceCert := d.Get("service_cert").(string)

	if d.HasChange("license_key") {
		licenseKey := d.Get("license_key").(string)
		if licenseKey != "" {
			license, err := c.setLicense(serviceCert, licenseKey)
			if err != nil {
				return diag.FromErr(err)
			}

			if err := setLicenseData(d, license); err != nil {
				return diag.FromErr(err)
			}
			return resourceVcenterReplicationManagerRead(ctx, d, m)
		}
	}

	if d.HasChange("lookup_service_url") {
		lsURL := d.Get("lookup_service_url").(string)
		lsThumbprint := d.Get("lookup_service_thumbprint").(string)
		if lsURL != "" {
			if err := c.setLookupService(lsURL, lsThumbprint, serviceCert); err != nil {
				return diag.FromErr(err)
			}

			return resourceVcenterReplicationManagerRead(ctx, d, m)
		}
	}

	return resourceVcenterReplicationManagerRead(ctx, d, m)
}

func resourceVcenterReplicationManagerDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*Client)

	serviceCert := d.Get("service_cert").(string)

	if err := c.removeVspherePlugin(serviceCert); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}

// util methods
func setLicenseData(d *schema.ResourceData, license *License) error {
	if err := d.Set("is_licensed", license.IsLicensed); err != nil {
		return fmt.Errorf("error setting is_licensed field: %s", err)
	}

	if err := d.Set("expiration_date", license.ExpirationDate); err != nil {
		return fmt.Errorf("error setting expiration_date field: %s", err)
	}

	return nil
}

func setSiteData(d *schema.ResourceData, site *SiteConfig) error {
	if err := d.Set("site", site.Site); err != nil {
		return fmt.Errorf("error setting site field: %s", err)
	}

	if err := d.Set("ls_url", site.LsURL); err != nil {
		return fmt.Errorf("error setting ls_url field: %s", err)
	}

	if err := d.Set("ls_thumbprint", site.LsThumbprint); err != nil {
		return fmt.Errorf("error setting ls_thumbprint field: %s", err)
	}

	if err := d.Set("tunnel_url", site.TunnelURL); err != nil {
		return fmt.Errorf("error setting tunnel_url field: %s", err)
	}

	if err := d.Set("tunnel_certificate", site.TunnelCertificate); err != nil {
		return fmt.Errorf("error setting tunnel_certificate field: %s", err)
	}

	return nil
}
