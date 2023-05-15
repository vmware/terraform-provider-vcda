/* Copyright 2023 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package vcda

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

func dataSourceVcdaServiceCert() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVcdaServiceCertRead,
		Schema: map[string]*schema.Schema{
			"datacenter_id": {
				Type:        schema.TypeString,
				Description: "The managed object ID of the datacenter where the virtual machine resides in.",
				Required:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The VM name of the appliance.",
				Required:    true,
			},
			"type": {
				Type: schema.TypeString,
				Description: "The type of the appliance role: manager, cloud, tunnel, replicator. " +
					"When not set returns an error.",
				Required: true,
			},
		},
	}
}

func dataSourceVcdaServiceCertRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*Client)
	vimClient := c.VimClient

	name := d.Get("name").(string)
	vmType := d.Get("type").(string)
	var vm *object.VirtualMachine
	var err error

	// TODO: implement looking for VM or template by UUID
	log.Printf("[DEBUG] Looking for VM or template by name/path %q", name)
	var dc *object.Datacenter
	if dcID, ok := d.GetOk("datacenter_id"); ok {
		dc, err = datacenterFromID(vimClient.vimClient, dcID.(string))
		if err != nil {
			return diag.FromErr(fmt.Errorf("cannot locate datacenter: %s", err))
		}
		log.Printf("[DEBUG] Datacenter for VM/template search: %s", dc.InventoryPath)
	}
	vm, err = FromPath(vimClient.vimClient, name, dc)

	if err != nil {
		return diag.FromErr(fmt.Errorf("error fetching virtual machine: %s", err))
	}

	props, err := Properties(vm)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error fetching virtual machine properties: %s", err))
	}

	if props.Config == nil {
		return diag.FromErr(fmt.Errorf("no configuration returned for virtual machine %q", vm.InventoryPath))
	}

	extraConfig := props.Config.ExtraConfig

	var extraConfigKey string

	switch vmType {
	case "manager":
		extraConfigKey = ManagerCertExtraConfigKey
	case "cloud":
		extraConfigKey = CloudCertExtraConfigKey
	case "tunnel":
		extraConfigKey = TunnelCertExtraConfigKey
	case "replicator":
		extraConfigKey = ReplicatorCertExtraConfigKey
	default:
		return diag.Errorf("unknown VM appliance type")
	}

	var applianceCert string
	for _, v := range extraConfig {
		ov := v.GetOptionValue()
		if ov.Key == extraConfigKey {
			applianceCert = ov.Value.(string)
		}
	}

	if applianceCert == "" {
		return diag.FromErr(fmt.Errorf("applaince certificate for %s was not found in virtual machine extraConfig", extraConfigKey))
	}

	d.SetId(applianceCert)

	return diags
}

// datacenterFromID locates a Datacenter by its managed object reference ID.
func datacenterFromID(client *govmomi.Client, id string) (*object.Datacenter, error) {
	finder := find.NewFinder(client.Client, false)

	ref := types.ManagedObjectReference{
		Type:  "Datacenter",
		Value: id,
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultAPITimeout)
	defer cancel()
	ds, err := finder.ObjectReference(ctx, ref)
	if err != nil {
		return nil, fmt.Errorf("could not find datacenter with id: %s: %s", id, err)
	}
	return ds.(*object.Datacenter), nil
}

// FromPath returns a VirtualMachine via its supplied path.
func FromPath(client *govmomi.Client, path string, dc *object.Datacenter) (*object.VirtualMachine, error) {
	finder := find.NewFinder(client.Client, false)
	if dc != nil {
		finder.SetDatacenter(dc)
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultAPITimeout)
	defer cancel()
	return finder.VirtualMachine(ctx, path)
}

// Properties is a convenience method that wraps fetching the
// VirtualMachine MO from its higher-level object.
func Properties(vm *object.VirtualMachine) (*mo.VirtualMachine, error) {
	log.Printf("[DEBUG] Fetching properties for VM %q", vm.InventoryPath)
	ctx, cancel := context.WithTimeout(context.Background(), defaultAPITimeout)
	defer cancel()
	var props mo.VirtualMachine
	if err := vm.Properties(ctx, vm.Reference(), nil, &props); err != nil {
		return nil, err
	}
	return &props, nil
}
