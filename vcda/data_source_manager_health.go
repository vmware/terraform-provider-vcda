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

func dataSourceVcdaManagerHealth() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVcdaManagerHealthRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"service_cert": {
				Type:        schema.TypeString,
				Description: "The certificate of the Cloud Director/vCenter Replication Manager Service.",
				Required:    true,
			},
			"manager_id": {
				Type:        schema.TypeString,
				Description: "The cloud manager instance id. **NOTE:** only required for the Cloud Director/Manager Service health info.",
				Optional:    true,
			},
			// Computed
			"id": {
				Type:        schema.TypeString,
				Description: "The health info task ID of the Cloud Director/vCenter Replication Manager Service.",
				Computed:    true,
			},
			"product_name": {
				Type:        schema.TypeString,
				Description: "The product name of the Cloud Director/vCenter Replication Manager Service.",
				Computed:    true,
			},
			"build_version": {
				Type:        schema.TypeString,
				Description: "The build version of the Cloud Director/vCenter Replication Manager Service.",
				Computed:    true,
			},
			"build_date": {
				Type:        schema.TypeFloat,
				Description: "The build date of the Cloud Director/vCenter Replication Manager Service.",
				Computed:    true,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Description: "The instance ID of the Cloud Director/vCenter Replication Manager Service.",
				Computed:    true,
			},
			"runtime_id": {
				Type:        schema.TypeString,
				Description: "The runtime ID of the Cloud Director/vCenter Replication Manager Service.",
				Computed:    true,
			},
			"current_time": {
				Type:        schema.TypeFloat,
				Description: "The current time of the Cloud Director/vCenter Replication Manager Service.",
				Computed:    true,
			},
			"address": {
				Type:        schema.TypeString,
				Description: "The address of the Cloud Director/vCenter Replication Manager Service.",
				Computed:    true,
			},
			"service_boot_timestamp": {
				Type:        schema.TypeInt,
				Description: "The service boot timestamp of the Cloud Director/vCenter Replication Manager Service.",
				Computed:    true,
			},
			"appliance_boot_timestamp": {
				Type:        schema.TypeFloat,
				Description: "The appliance boot timestamp of the Cloud Director/vCenter Replication Manager Service.",
				Computed:    true,
			},
			"disk_usage": {
				Type:        schema.TypeMap,
				Description: "The disk usage of the Cloud Director/vCenter Replication Manager Service.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"local_replicators_ls_mismatch_error_code": {
				Type:        schema.TypeString,
				Description: "The local replicators lookup service mismatch error code of the Cloud Director/vCenter Replication Manager Service.",
				Computed:    true,
			},
			"local_replicators_ls_mismatch_error_msg": {
				Type:        schema.TypeString,
				Description: "The local replicators lookup service mismatch error message of the Cloud Director/vCenter Replication Manager Service.",
				Computed:    true,
			},
			"local_replicators_ls_mismatch_error_args": {
				Type:        schema.TypeList,
				Description: "The local replicators lookup service mismatch error arguments of the Cloud Director/vCenter Replication Manager Service.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"local_replicators_ls_mismatch_error_stacktrace": {
				Type:        schema.TypeString,
				Description: "The local replicators lookup service mismatch error stacktrace of the Cloud Director/vCenter Replication Manager Service.",
				Computed:    true,
			},
			"sso_admin_error_code": {
				Type:        schema.TypeString,
				Description: "The sso admin error code of the Cloud Director/vCenter Replication Manager Service.",
				Computed:    true,
			},
			"sso_admin_error_msg": {
				Type:        schema.TypeString,
				Description: "The sso admin error message of the Cloud Director/vCenter Replication Manager Service.",
				Computed:    true,
			},
			"sso_admin_error_args": {
				Type:        schema.TypeList,
				Description: "The sso admin error arguments of the Cloud Director/vCenter Replication Manager Service.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"sso_admin_error_stacktrace": {
				Type:        schema.TypeString,
				Description: "The sso admin error stacktrace of the Cloud Director/vCenter Replication Manager Service.",
				Computed:    true,
			},
			"ls_error_code": {
				Type:        schema.TypeString,
				Description: "The lookup service error code of the Cloud Director/vCenter Replication Manager Service.",
				Computed:    true,
			},
			"ls_error_msg": {
				Type:        schema.TypeString,
				Description: "The lookup service error message of the Cloud Director/vCenter Replication Manager Service.",
				Computed:    true,
			},
			"ls_error_args": {
				Type:        schema.TypeList,
				Description: "The lookup service error arguments of the Cloud Director/vCenter Replication Manager Service.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"ls_error_stacktrace": {
				Type:        schema.TypeString,
				Description: "The lookup service error stacktrace of the Cloud Director/vCenter Replication Manager Service.",
				Computed:    true,
			},
			"db_error_code": {
				Type:        schema.TypeString,
				Description: "The database error code of the Cloud Director/vCenter Replication Manager Service.",
				Computed:    true,
			},
			"db_error_msg": {
				Type:        schema.TypeString,
				Description: "The database error message of the Cloud Director/vCenter Replication Manager Service.",
				Computed:    true,
			},
			"db_error_args": {
				Type:        schema.TypeList,
				Description: "The database error arguments of the Cloud Director/vCenter Replication Manager Service.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"db_error_stacktrace": {
				Type:        schema.TypeString,
				Description: "The database error stacktrace of the Cloud Director/vCenter Replication Manager Service.",
				Computed:    true,
			},
			"ntp_error_code": {
				Type:        schema.TypeString,
				Description: "The NTP error code of the Cloud Director/vCenter Replication Manager Service.",
				Computed:    true,
			},
			"ntp_error_msg": {
				Type:        schema.TypeString,
				Description: "The NTP error message of the Cloud Director/vCenter Replication Manager Service.",
				Computed:    true,
			},
			"ntp_error_args": {
				Type:        schema.TypeList,
				Description: "The NTP error arguments of the Cloud Director/vCenter Replication Manager Service.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"ntp_error_stacktrace": {
				Type:        schema.TypeString,
				Description: "The NTP error stacktrace of the Cloud Director/vCenter Replication Manager Service.",
				Computed:    true,
			},
			"offline_replicators_ids": {
				Type:        schema.TypeList,
				Description: "A list of the offline replicators IDs of the Cloud Director/vCenter Replication Manager Service.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"online_replicators_ids": {
				Type:        schema.TypeList,
				Description: "A list of the online replicators IDs of the Cloud Director/vCenter Replication Manager Service.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"local_replicators_ids": {
				Type:        schema.TypeList,
				Description: "A list of the local replicators IDs of the Cloud Director/vCenter Replication Manager Service.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"tunnels_ids": {
				Type:        schema.TypeList,
				Description: "A list of the tunnels IDs of the Cloud Director/vCenter Replication Manager Service.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceVcdaManagerHealthRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	return getManagerHealthInfo(c, d)
}

func getManagerHealthInfo(c *Client, d *schema.ResourceData) diag.Diagnostics {
	var diags diag.Diagnostics

	serviceCert := d.Get("service_cert").(string)
	managerID := d.Get("manager_id").(string)
	taskID := d.Get("id").(string)

	task, err := c.getTask(serviceCert, taskID)
	if err != nil {
		return diag.FromErr(err)
	}

	if managerID != "" {
		cloudHealth, ok := task.Result.(map[string]interface{})
		if !ok {
			return diag.FromErr(fmt.Errorf("unexpected type for taskResult.Result"))
		}

		cloudManagerHealth, ok := cloudHealth["managerHealth"].(map[string]interface{})
		if !ok {
			return diag.FromErr(fmt.Errorf("unexpected type for cloudManagerHealth"))
		}

		if err := setManagerHealthInfoData(d, cloudManagerHealth); err != nil {
			return diag.FromErr(err)
		}
	} else {
		managerHealth, ok := task.Result.(map[string]interface{})
		if !ok {
			return diag.FromErr(fmt.Errorf("unexpected type for taskResult.Result"))
		}

		if err := setManagerHealthInfoData(d, managerHealth); err != nil {
			return diag.FromErr(err)
		}
	}

	return diags
}

func setManagerHealthInfoData(d *schema.ResourceData, managerHealth map[string]interface{}) error {
	err := setHealthInfoData(d, managerHealth)
	if err != nil {
		return err
	}

	if localReplicatorsMismatchError, ok := managerHealth["localReplicatorsLsMismatch"].(map[string]interface{}); ok {
		if err := setErrorData(d, localReplicatorsMismatchError, "local_replicators_ls_mismatch"); err != nil {
			return err
		}
	}

	if ssoAdminError, ok := managerHealth["ssoAdminError"].(map[string]interface{}); ok {
		if err := setErrorData(d, ssoAdminError, "sso_admin"); err != nil {
			return err
		}
	}

	if tunnelConnectivity, ok := managerHealth["tunnelConnectivity"].([]interface{}); ok {
		if err := setTunnelIDs(d, tunnelConnectivity); err != nil {
			return err
		}
	}

	if offlineReplicators, ok := managerHealth["offlineReplicators"].([]interface{}); ok {
		if err := setReplicatorsIDs(d, offlineReplicators, false); err != nil {
			return err
		}
	}

	if onlineReplicators, ok := managerHealth["onlineReplicators"].([]interface{}); ok {
		if err := setReplicatorsIDs(d, onlineReplicators, true); err != nil {
			return err
		}
	}

	if localReplicators, ok := managerHealth["localReplicatorsHealth"].([]interface{}); ok {
		if err := setLocalReplicatorsIDs(d, localReplicators); err != nil {
			return err
		}
	}

	return nil
}

func setReplicatorsIDs(d *schema.ResourceData, replicators []interface{}, isOnline bool) error {
	var replicatorsIDs []string
	for _, replicator := range replicators {
		if repMap, ok := replicator.(map[string]interface{}); ok {
			if replicatorID, ok := repMap["id"].(string); ok {
				replicatorsIDs = append(replicatorsIDs, replicatorID)
			}
		}
	}

	if isOnline {
		if err := d.Set("online_replicators_ids", replicatorsIDs); err != nil {
			return fmt.Errorf("error setting online_replicators_ids field: %s", err)
		}
	} else {
		if err := d.Set("offline_replicators_ids", replicatorsIDs); err != nil {
			return fmt.Errorf("error setting offline_replicators_ids field: %s", err)
		}
	}

	return nil
}

func setLocalReplicatorsIDs(d *schema.ResourceData, replicators []interface{}) error {
	var replicatorsIDs []string
	for _, replicator := range replicators {
		if repMap, ok := replicator.(map[string]interface{}); ok {
			if replicatorID, ok := repMap["instanceId"].(string); ok {
				replicatorsIDs = append(replicatorsIDs, replicatorID)
			}
		}
	}

	if err := d.Set("local_replicators_ids", replicatorsIDs); err != nil {
		return fmt.Errorf("error setting local_replicators_ids field: %s", err)
	}

	return nil
}
