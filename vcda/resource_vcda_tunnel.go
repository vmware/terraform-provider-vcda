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
				Type:        schema.TypeString,
				Description: "Tunnel appliance VM thumbprint.",
				Required:    true,
			},
			"url": {
				Type:        schema.TypeString,
				Description: "Tunnel service URL.",
				Required:    true,
			},
			"certificate": {
				Type:        schema.TypeString,
				Description: "Tunnel service certificate.",
				Required:    true,
			},
			"root_password": {
				Type:        schema.TypeString,
				Description: "Tunnel service root password.",
				Required:    true,
			},

			// computed
			"tunnel_url": {
				Type:        schema.TypeString,
				Description: "Tunnel service URL.",
				Computed:    true,
			},
			"tunnel_certificate": {
				Type:        schema.TypeString,
				Description: "Tunnel service certificate.",
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

	d.SetId(tunnelConfig.TunnelURL)

	return resourceVcdaTunnelRead(ctx, d, m)
}

func resourceVcdaTunnelRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*Client)

	serviceCert := d.Get("service_cert").(string)

	site, err := c.getCloudSiteConfig(serviceCert)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setTunnelData(d, site); err != nil {
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

		d.SetId(tunnelConfig.TunnelURL)

		return resourceVcdaTunnelRead(ctx, d, m)
	}
	return diags
}

func resourceVcdaTunnelDelete(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	d.SetId("")

	return diags
}

func setTunnelData(d *schema.ResourceData, tunnelConfig *CloudSiteConfig) error {
	if err := d.Set("tunnel_url", tunnelConfig.TunnelURL); err != nil {
		return fmt.Errorf("error setting tunnel_url field: %s", err)
	}

	if err := d.Set("tunnel_certificate", tunnelConfig.TunnelCertificate); err != nil {
		return fmt.Errorf("error setting tunnel_certificate field: %s", err)
	}

	return nil
}
