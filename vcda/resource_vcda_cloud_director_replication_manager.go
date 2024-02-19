// Copyright (c) 2023-2024 Broadcom. All Rights Reserved.
// Broadcom Confidential. The term "Broadcom" refers to Broadcom Inc.
// and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vcda

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceVcdaCloudDirectorReplicationManager() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCloudDirectorReplicationManagerCreate,
		ReadContext:   resourceCloudDirectorReplicationManagerRead,
		UpdateContext: resourceCloudDirectorReplicationManagerUpdate,
		DeleteContext: resourceCloudDirectorReplicationManagerDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"service_cert": {
				Type:        schema.TypeString,
				Description: "The certificate of the Cloud Director Replication Manager Service.",
				Required:    true,
			},
			"vcd_thumbprint": {
				Type: schema.TypeString,
				Description: "The thumbprint of the Cloud Director service. It can either be computed from " +
					"the `vcda_remote_services_thumbprint` data source or provided directly as a SHA-256 fingerprint.",
				Required: true,
			},
			"lookup_service_thumbprint": {
				Type: schema.TypeString,
				Description: "The thumbprint of the vCenter Server Lookup service. It can either be computed from " +
					"the `vcda_remote_services_thumbprint` data source or provided directly as a SHA-256 fingerprint.",
				Required: true,
			},
			"license_key": {
				Type:        schema.TypeString,
				Description: "The license key for VMware Cloud Director Availability.",
				Required:    true,
			},
			"site_name": {
				Type:        schema.TypeString,
				Description: "The site name of the Cloud Director Replication Manager.",
				Required:    true,
			},
			"site_description": {
				Type:        schema.TypeString,
				Description: "The site description of the Cloud Director Replication Manager.",
				Optional:    true,
			},
			"public_endpoint_address": {
				Type:        schema.TypeString,
				Description: "The public API endpoint address.",
				Required:    true,
			},
			"public_endpoint_port": {
				Type:        schema.TypeInt,
				Description: "The public API endpoint port.",
				Required:    true,
			},
			"vcd_username": {
				Type:        schema.TypeString,
				Description: "Cloud Director user name.",
				Required:    true,
			},
			"vcd_password": {
				Type:        schema.TypeString,
				Description: "Cloud Director password.",
				Required:    true,
			},
			"vcd_url": {
				Type: schema.TypeString,
				Description: "This is the URL for the Cloud Director API endpoint. " +
					"For example, https://server.domain.com/api.",
				Required: true,
			},
			"lookup_service_url": {
				Type: schema.TypeString,
				Description: "The URL of the vCenter Server Lookup service. " +
					"For example, https://server.domain.com/lookupservice/sdk.",
				Required: true,
			},

			// computed:
			"is_licensed": {
				Type:        schema.TypeBool,
				Description: "Flag indicating whether the solution is licensed.",
				Computed:    true,
			},
			"expiration_date": {
				Type:        schema.TypeInt,
				Description: "VMware Cloud Director Availability license expiration date.",
				Computed:    true,
			},
			"ls_url": {
				Type:        schema.TypeString,
				Description: "The URL of the vCenter Server Lookup service.",
				Computed:    true,
			},
			"ls_thumbprint": {
				Type:        schema.TypeString,
				Description: "SHA-256 vCenter Server Lookup service thumbprint.",
				Computed:    true,
			},
			"local_site": {
				Type:        schema.TypeString,
				Description: "Cloud Director Replication Manager local site name.",
				Computed:    true,
			},
			"local_site_description": {
				Type:        schema.TypeString,
				Description: "Cloud Director Replication Manager local site description.",
				Computed:    true,
			},
			"vcloud_url": {
				Type:        schema.TypeString,
				Description: "Cloud Director URL.",
				Computed:    true,
			},
			"vcloud_thumbprint": {
				Type:        schema.TypeString,
				Description: "Cloud Director thumbprint.",
				Computed:    true,
			},
			"vcloud_username": {
				Type:        schema.TypeString,
				Description: "Cloud Director user name.",
				Computed:    true,
			},
			"tunnel_url": {
				Type:        schema.TypeString,
				Description: "Tunnel Service URL.",
				Computed:    true,
			},
			"tunnel_certificate": {
				Type:        schema.TypeString,
				Description: "Tunnel Service certificate.",
				Computed:    true,
			},
			"is_combined": {
				Type:        schema.TypeBool,
				Description: "Flag indicating whether the appliance role is Cloud Director Combined Appliance.",
				Computed:    true,
			},

			// endpoints
			"mgmt_address": {
				Type:        schema.TypeString,
				Description: "Effective endpoint management address.",
				Computed:    true,
			},
			"mgmt_port": {
				Type:        schema.TypeInt,
				Description: "Effective endpoint management port.",
				Computed:    true,
			},
			"mgmt_public_address": {
				Type:        schema.TypeString,
				Description: "Effective endpoint management public address.",
				Computed:    true,
			},
			"mgmt_public_port": {
				Type:        schema.TypeInt,
				Description: "Effective endpoint management public port.",
				Computed:    true,
			},
			"api_address": {
				Type:        schema.TypeString,
				Description: "Effective endpoint API address.",
				Computed:    true,
			},
			"api_port": {
				Type:        schema.TypeInt,
				Description: "Effective endpoint API port.",
				Computed:    true,
			},
			"api_public_address": {
				Type:        schema.TypeString,
				Description: "Effective endpoint API public address.",
				Computed:    true,
			},
			"api_public_port": {
				Type:        schema.TypeInt,
				Description: "Effective endpoint API public port.",
				Computed:    true,
			},
		},
	}

}

func resourceCloudDirectorReplicationManagerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	serviceCert := d.Get("service_cert").(string)
	vcdThumbprint := d.Get("vcd_thumbprint").(string)
	lsThumbprint := d.Get("lookup_service_thumbprint").(string)

	licenseKey := d.Get("license_key").(string)
	siteName := d.Get("site_name").(string)
	siteDescription := d.Get("site_description").(string)

	endpointAddress := d.Get("public_endpoint_address").(string)
	endpointPort := d.Get("public_endpoint_port").(int)

	vcdUsername := d.Get("vcd_username").(string)
	vcdPassword := d.Get("vcd_password").(string)
	vcdURL := d.Get("vcd_url").(string)

	lsURL := d.Get("lookup_service_url").(string)

	// set license
	license, err := c.setLicense(serviceCert, licenseKey)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setLicenseData(d, license); err != nil {
		return diag.FromErr(err)
	}

	// set site name
	site, err := c.setCloudSiteName(siteName, siteDescription, serviceCert)
	if err != nil {
		return diag.FromErr(err)
	}

	// set public API endpoint
	if err := c.setPublicEndpoint(endpointAddress, endpointPort, serviceCert); err != nil {
		return diag.FromErr(err)
	}

	// set vcd
	vcdThumb := vcdThumbprint
	if !strings.HasPrefix(vcdThumbprint, "SHA-256:") {
		vcdThumb = "SHA-256:" + vcdThumbprint
	}
	if err := c.setVcloud(vcdUsername, vcdPassword, vcdURL, vcdThumb, serviceCert); err != nil {
		return diag.FromErr(err)
	}

	// set cloud lookup service
	if err := c.setLookupService(lsURL, lsThumbprint, serviceCert); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(site.ID)

	err = retry.RetryContext(context.Background(), d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		isConfigured, err := c.isConfigured(serviceCert)

		if err != nil {
			return retry.NonRetryableError(err)
		}

		if !isConfigured.IsConfigured {
			return retry.RetryableError(fmt.Errorf("service is not configured yet"))
		}

		return nil
	})
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceCloudDirectorReplicationManagerRead(ctx, d, m)
}

func resourceCloudDirectorReplicationManagerRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*Client)

	serviceCert := d.Get("service_cert").(string)

	vcdaSite, err := c.getCloudSiteConfig(serviceCert)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setCloudSiteData(d, vcdaSite); err != nil {
		return diag.FromErr(err)
	}

	endpoints, err := c.getEndpoints(serviceCert)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setEndpointData(d, endpoints.Effective); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceCloudDirectorReplicationManagerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*Client)

	serviceCert := d.Get("service_cert").(string)

	if d.HasChange("license_key") {
		licenseKey := d.Get("license_key").(string)
		if licenseKey != "" {
			vcdaLicense, err := c.setLicense(serviceCert, licenseKey)
			if err != nil {
				return diag.FromErr(err)
			}

			if err := setLicenseData(d, vcdaLicense); err != nil {
				return diag.FromErr(err)
			}
			return resourceCloudDirectorReplicationManagerRead(ctx, d, m)
		}
	}

	if d.HasChange("lookup_service_url") {
		lsURL := d.Get("lookup_service_url").(string)
		lsThumbprint := d.Get("lookup_service_thumbprint").(string)
		if lsURL != "" {
			if err := c.setLookupService(lsURL, lsThumbprint, serviceCert); err != nil {
				return diag.FromErr(err)
			}

			return resourceCloudDirectorReplicationManagerRead(ctx, d, m)
		}
	}

	if d.HasChange("vcd_url") || d.HasChange("vcd_password") || d.HasChange("vcd_username") {
		vcdUsername := d.Get("vcd_username").(string)
		vcdPassword := d.Get("vcd_password").(string)
		vcdURL := d.Get("vcd_url").(string)
		vcdThumbprint := d.Get("vcd_thumbprint").(string)

		if err := c.setVcloud(vcdUsername, vcdPassword, vcdURL, vcdThumbprint, serviceCert); err != nil {
			return diag.FromErr(err)
		}

		return resourceCloudDirectorReplicationManagerRead(ctx, d, m)
	}

	if d.HasChange("public_endpoint_address") || d.HasChange("public_endpoint_port") {
		endpointAddress := d.Get("public_endpoint_address").(string)
		endpointPort := d.Get("public_endpoint_port").(int)

		if err := c.setPublicEndpoint(endpointAddress, endpointPort, serviceCert); err != nil {
			return diag.FromErr(err)
		}

		return resourceCloudDirectorReplicationManagerRead(ctx, d, m)
	}

	return diags
}

func resourceCloudDirectorReplicationManagerDelete(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	d.SetId("")

	return diags
}

func setCloudSiteData(d *schema.ResourceData, site *CloudSiteConfig) error {
	if err := d.Set("ls_url", site.LsURL); err != nil {
		return fmt.Errorf("error setting ls_url field: %s", err)
	}

	if err := d.Set("ls_thumbprint", site.LsThumbprint); err != nil {
		return fmt.Errorf("error setting site field: %s", err)
	}

	if err := d.Set("local_site", site.LocalSite); err != nil {
		return fmt.Errorf("error setting local_site field: %s", err)
	}

	if err := d.Set("local_site_description", site.LocalSiteDescription); err != nil {
		return fmt.Errorf("error setting local_site_description field: %s", err)
	}

	if err := d.Set("vcloud_url", site.VcdURL); err != nil {
		return fmt.Errorf("error setting vcloud_url field: %s", err)
	}

	if err := d.Set("vcloud_thumbprint", site.VcdThumbprint); err != nil {
		return fmt.Errorf("error setting vcloud_thumbprint field: %s", err)
	}

	if err := d.Set("vcloud_username", site.VcdUsername); err != nil {
		return fmt.Errorf("error setting vcloud_username field: %s", err)
	}

	if err := d.Set("tunnel_url", site.TunnelURL); err != nil {
		return fmt.Errorf("error setting tunnel_url field: %s", err)
	}

	if err := d.Set("tunnel_certificate", site.TunnelCertificate); err != nil {
		return fmt.Errorf("error setting tunnel_certificate field: %s", err)
	}

	if err := d.Set("is_combined", site.IsCombined); err != nil {
		return fmt.Errorf("error setting is_combined field: %s", err)
	}

	return nil
}

func setEndpointData(d *schema.ResourceData, endpoint EndpointConfig) error {
	if err := d.Set("mgmt_address", endpoint.MgmtAddress); err != nil {
		return fmt.Errorf("error setting mgmt_address field: %s", err)
	}
	if err := d.Set("mgmt_port", endpoint.MgmtPort); err != nil {
		return fmt.Errorf("error setting mgmt_port field: %s", err)
	}
	if err := d.Set("mgmt_public_address", endpoint.MgmtPublicAddress); err != nil {
		return fmt.Errorf("error setting mgmt_public_address field: %s", err)
	}
	if err := d.Set("mgmt_public_port", endpoint.MgmtPublicPort); err != nil {
		return fmt.Errorf("error setting mgmt_public_port field: %s", err)
	}
	if err := d.Set("api_address", endpoint.APIAddress); err != nil {
		return fmt.Errorf("error setting api_address field: %s", err)
	}
	if err := d.Set("api_port", endpoint.APIPort); err != nil {
		return fmt.Errorf("error setting api_port field: %s", err)
	}
	if err := d.Set("api_public_address", endpoint.APIPublicAddress); err != nil {
		return fmt.Errorf("error setting api_public_address field: %s", err)
	}
	if err := d.Set("api_public_port", endpoint.APIPublicPort); err != nil {
		return fmt.Errorf("error setting api_public_port field: %s", err)
	}

	return nil
}
