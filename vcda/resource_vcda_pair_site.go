// Copyright (c) 2023-2024 Broadcom. All Rights Reserved.
// Broadcom Confidential. The term "Broadcom" refers to Broadcom Inc.
// and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vcda

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceVcdaPairSite() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePairSiteCreate,
		ReadContext:   resourcePairSiteRead,
		UpdateContext: resourcePairSiteUpdate,
		DeleteContext: resourcePairSiteDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"service_cert": {
				Type:        schema.TypeString,
				Description: "The certificate of the Cloud Director/vCenter Replication Management Appliance.",
				Required:    true,
			},
			"api_thumbprint": {
				Type: schema.TypeString,
				Description: "The thumbprint of the to-be paired Cloud Director/vCenter Replication Management Appliance. It can either be computed from " +
					"the `vcda_remote_services_thumbprint` data source or provided directly as a SHA-256 fingerprint.",
				Required: true,
			},
			"api_url": {
				Type:        schema.TypeString,
				Description: "The API URL address/endpoint of the to-be paired Cloud Director/vCenter Replication Management Appliance.",
				Required:    true,
			},
			"pairing_description": {
				Type:        schema.TypeString,
				Description: "The description of the pairing.",
				Optional:    true,
			},
			"site": {
				Type: schema.TypeString,
				Description: "The site name of the to-be paired Cloud Director Replication Management Appliance. " +
					"Only required for pairing a Cloud Director Replication Management Appliance to another Cloud Director Replication Management Appliance.",
				Optional: true,
			},
			//computed
			"site_id": {
				Type: schema.TypeString,
				Description: "The site ID of the paired vCenter Replication Management Appliance. " +
					"Computed only for pairing a vCenter Replication Management Appliance to another vCenter Replication Management Appliance.",
				Computed: true,
			},
			"site_name": {
				Type:        schema.TypeString,
				Description: "The site name of the paired Cloud Director/vCenter Replication Management Appliance.",
				Computed:    true,
			},
			"site_description": {
				Type:        schema.TypeString,
				Description: "The site description of the paired Cloud Director/vCenter Replication Management Appliance.",
				Computed:    true,
			},
			"api_public_url": {
				Type:        schema.TypeString,
				Description: "The public API URL address of the paired Cloud Director/vCenter Replication Management Appliance.",
				Computed:    true,
			},
			"api_version": {
				Type:        schema.TypeString,
				Description: "The API version of the paired Cloud Director/vCenter Replication Management Appliance.",
				Computed:    true,
			},
			"build_version": {
				Type: schema.TypeString,
				Description: "The build version of the paired Cloud Director Replication Management Appliance. " +
					"Computed only for pairing a Cloud Director Replication Management Appliance to another Cloud Director Replication Management Appliance.",
				Computed: true,
			},
			"is_provider_deployment": {
				Type: schema.TypeBool,
				Description: "A flag that indicates whether the paired vCenter Replication Management Appliance is of type provider. " +
					"Computed only for pairing a vCenter Replication Management Appliance to another vCenter Replication Management Appliance.",
				Computed: true,
			},
		},
	}

}

func resourcePairSiteCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	serviceCert := d.Get("service_cert").(string)
	apiThumbprint := d.Get("api_thumbprint").(string)
	apiURL := d.Get("api_url").(string)
	pairingDescription := d.Get("pairing_description").(string)
	siteName := d.Get("site").(string)

	taskID, err := c.pairSite(serviceCert, apiThumbprint, apiURL, pairingDescription, siteName)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(*taskID)

	err = retry.RetryContext(context.Background(), d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		task, err := c.getTask(serviceCert, *taskID)

		if err != nil {
			return retry.NonRetryableError(err)
		}

		if task.State == "FAILED" {
			return retry.NonRetryableError(fmt.Errorf("pair site task failed with Code:" + task.Error.Code + ", Msg: " + task.Error.Msg))
		} else if task.State == "RUNNING" {
			return retry.RetryableError(fmt.Errorf("expected pair site task to be completed but was in state %s", task.State))
		}

		return nil
	})

	if err != nil {
		return diag.FromErr(err)
	}

	return resourcePairSiteRead(ctx, d, m)
}

func resourcePairSiteRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*Client)

	apiURL := d.Get("api_url").(string)
	serviceCert := d.Get("service_cert").(string)
	siteName := d.Get("site").(string)

	if siteName != "" {
		cloudSite, err := c.getCloudSite(serviceCert, apiURL)
		if err != nil {
			return diag.FromErr(err)
		}

		if err := setPairedCloudSiteData(cloudSite, d); err != nil {
			return diag.FromErr(err)
		}
	} else {
		vcenterSite, err := c.getVcenterSite(serviceCert, apiURL)
		if err != nil {
			return diag.FromErr(err)
		}

		if err := setPairedVcenterSiteData(vcenterSite, d); err != nil {
			return diag.FromErr(err)
		}
	}

	return diags
}

func resourcePairSiteUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*Client)

	serviceCert := d.Get("service_cert").(string)
	apiThumbprint := d.Get("api_thumbprint").(string)
	siteName := d.Get("site").(string)
	siteID := d.Get("site_id").(string)

	var site string
	if siteID != "" {
		// only vcenter sites have site id
		site = siteID
	} else {
		// cloud site
		site = siteName
	}

	if d.HasChange("api_url") || d.HasChange("pairing_description") {
		apiURL := d.Get("api_url").(string)
		pairingDescription := d.Get("pairing_description").(string)

		taskID, err := c.repairSite(serviceCert, site, apiThumbprint, apiURL, pairingDescription)
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(*taskID)

		err = retry.RetryContext(context.Background(), d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			task, err := c.getTask(serviceCert, *taskID)

			if err != nil {
				return retry.NonRetryableError(err)
			}

			if task.State == "FAILED" {
				return retry.NonRetryableError(fmt.Errorf("re-pair task failed with Code:" + task.Error.Code + ", Msg: " + task.Error.Msg))
			} else if task.State == "RUNNING" {
				return retry.RetryableError(fmt.Errorf("expected re-pair task to be completed but was in state %s", task.State))
			}

			return nil
		})

		if err != nil {
			return diag.FromErr(err)
		}

		return resourcePairSiteRead(ctx, d, m)
	}

	return diags
}

func resourcePairSiteDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*Client)

	siteID := d.Get("site_id").(string)
	siteName := d.Get("site").(string)
	serviceCert := d.Get("service_cert").(string)

	var site string
	if siteID != "" {
		// only vcenter sites have site id
		site = siteID
	} else {
		// cloud site
		site = siteName
	}

	taskID, err := c.unpairSite(serviceCert, site)
	if err != nil {
		return diag.FromErr(err)
	}

	err = retry.RetryContext(context.Background(), d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		task, err := c.getTask(serviceCert, *taskID)

		if err != nil {
			return retry.NonRetryableError(err)
		}

		if task.State == "FAILED" {
			return retry.NonRetryableError(fmt.Errorf("unpair task failed with Code:" + task.Error.Code + ", Msg: " + task.Error.Msg))
		} else if task.State == "RUNNING" {
			return retry.RetryableError(fmt.Errorf("expected unpair site to be completed but was in state %s", task.State))
		}

		return nil
	})

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}

func setPairedCloudSiteData(site *CloudSite, d *schema.ResourceData) error {
	if err := d.Set("site_name", site.Site); err != nil {
		return fmt.Errorf("error setting site_name field: %s", err)
	}
	if err := d.Set("site_description", site.Description); err != nil {
		return fmt.Errorf("error setting site_description field: %s", err)
	}
	if err := d.Set("api_public_url", site.APIPublicURL); err != nil {
		return fmt.Errorf("error setting api_public_url field: %s", err)
	}
	if err := d.Set("api_version", site.APIVersion); err != nil {
		return fmt.Errorf("error setting api_version field: %s", err)
	}
	if err := d.Set("build_version", site.BuildVersion); err != nil {
		return fmt.Errorf("error setting build_version field: %s", err)
	}

	return nil
}

func setPairedVcenterSiteData(site *VcenterSite, d *schema.ResourceData) error {
	if err := d.Set("site_id", site.ID); err != nil {
		return fmt.Errorf("error setting site_id field: %s", err)
	}
	if err := d.Set("site_name", site.Site); err != nil {
		return fmt.Errorf("error setting site_name field: %s", err)
	}
	if err := d.Set("site_description", site.Description); err != nil {
		return fmt.Errorf("error setting site_description field: %s", err)
	}
	if err := d.Set("api_public_url", site.APIPublicURL); err != nil {
		return fmt.Errorf("error setting api_public_url field: %s", err)
	}
	if err := d.Set("api_version", site.APIVersion); err != nil {
		return fmt.Errorf("error setting api_version field: %s", err)
	}
	if err := d.Set("is_provider_deployment", site.IsProviderDeployment); err != nil {
		return fmt.Errorf("error setting is_provider_deployment field: %s", err)
	}

	return nil
}
