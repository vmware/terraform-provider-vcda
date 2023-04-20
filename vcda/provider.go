package vcda

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// defaultAPITimeout is a default timeout value that is passed to functions
// requiring contexts, and other various waiters.
var defaultAPITimeout = time.Minute * 5

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"vcda_ip": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc(VcdaIP, nil),
				Description: "The VCDA appliance IP address.",
			},
			"local_user": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc(LocalUser, nil),
				Description: "The VCDA appliance local user.",
			},
			"local_password": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc(LocalPassword, nil),
				Description: "The VCDA appliance local password.",
			},
			"vsphere_user": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc(VsphereUser, nil),
				Description: "The user name for vSphere API operations.",
			},
			"vsphere_password": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc(VspherePassword, nil),
				Description: "The user password for vSphere API operations.",
			},
			"vsphere_server": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc(VsphereServer, nil),
				Description: "The vSphere Server name for vSphere API operations.",
			},
			"vsphere_allow_unverified_ssl": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(VsphereAllowUnverifiedSSL, true),
				Description: "If set, VMware vSphere client will establish an insecure TLS connection without performing certificate validations.",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"vcda_appliance_password":                 resourceVcdaAppliancePassword(),
			"vcda_vcenter_replication_manager":        resourceVcdaVcenterReplicationManager(),
			"vcda_cloud_director_replication_manager": resourceVcdaCloudDirectorReplicationManager(),
			"vcda_replicator":                         resourceVcdaReplicator(),
			"vcda_tunnel":                             resourceVcdaTunnel(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"vcda_remote_services_thumbprint": dataSourceVcdaRemoteServicesThumbprint(),
			"vcda_service_cert":               dataSourceVcdaServiceCert(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	vcdaIP := d.Get("vcda_ip").(string)

	localUser := d.Get("local_user").(string)
	localPassword := d.Get("local_password").(string)
	if len(localPassword) <= 0 {
		return nil, diag.Errorf("local_password cannot be empty")
	}

	c, err := NewConfig(d)
	if err != nil {
		return nil, diag.FromErr(err)
	}
	vimClient, err := c.VimClient()
	if err != nil {
		return nil, diag.Errorf("could not initialize vim client: %s", err)
	}
	client := Client{*vimClient, vcdaIP, localUser, localPassword}

	return &client, nil
}