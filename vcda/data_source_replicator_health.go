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

func dataSourceVcdaReplicatorHealth() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVcdaReplicatorHealthRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"service_cert": {
				Type:        schema.TypeString,
				Description: "The certificate of the Cloud Director/vCenter Replication Manager Service.",
				Required:    true,
			},
			"replicator_id": {
				Type:        schema.TypeString,
				Description: "The replicator service instance ID.",
				Required:    true,
			},
			// Computed
			"id": {
				Type:        schema.TypeString,
				Description: "The health info task ID.",
				Computed:    true,
			},
			"product_name": {
				Type:        schema.TypeString,
				Description: "The product name of the Replicator Service.",
				Computed:    true,
			},
			"build_version": {
				Type:        schema.TypeString,
				Description: "The build version of the Replicator Service.",
				Computed:    true,
			},
			"build_date": {
				Type:        schema.TypeFloat,
				Description: "The build date of the Replicator Service.",
				Computed:    true,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Description: "The instance ID of the Replicator Service.",
				Computed:    true,
			},
			"runtime_id": {
				Type:        schema.TypeString,
				Description: "The runtime ID of the Replicator Service.",
				Computed:    true,
			},
			"current_time": {
				Type:        schema.TypeFloat,
				Description: "The current time of the Replicator Service.",
				Computed:    true,
			},
			"address": {
				Type:        schema.TypeString,
				Description: "The address of the Replicator Service.",
				Computed:    true,
			},
			"service_boot_timestamp": {
				Type:        schema.TypeInt,
				Description: "The service boot timestamp of the Replicator Service.",
				Computed:    true,
			},
			"appliance_boot_timestamp": {
				Type:        schema.TypeFloat,
				Description: "The appliance boot timestamp of the Replicator Service.",
				Computed:    true,
			},
			"disk_usage": {
				Type:        schema.TypeMap,
				Description: "The disk usage of the Replicator Service.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"lwd_error_code": {
				Type:        schema.TypeString,
				Description: "The LWD error code of the Replicator Service.",
				Computed:    true,
			},
			"lwd_error_msg": {
				Type:        schema.TypeString,
				Description: "The LWD error message of the Replicator Service.",
				Computed:    true,
			},
			"lwd_error_args": {
				Type:        schema.TypeList,
				Description: "The LWD error arguments of the Replicator Service.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"lwd_error_stacktrace": {
				Type:        schema.TypeString,
				Description: "The LWD error stacktrace of the Replicator Service.",
				Computed:    true,
			},
			"hbr_error_code": {
				Type:        schema.TypeString,
				Description: "The HBR error code of the Replicator Service.",
				Computed:    true,
			},
			"hbr_error_msg": {
				Type:        schema.TypeString,
				Description: "The HBR error message of the Replicator Service.",
				Computed:    true,
			},
			"hbr_error_args": {
				Type:        schema.TypeList,
				Description: "The HBR error arguments of the Replicator Service.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"hbr_error_stacktrace": {
				Type:        schema.TypeString,
				Description: "The HBR error stacktrace of the Replicator Service.",
				Computed:    true,
			},
			"h4dm_error_code": {
				Type:        schema.TypeString,
				Description: "The H4DM error code of the Replicator Service.",
				Computed:    true,
			},
			"h4dm_error_msg": {
				Type:        schema.TypeString,
				Description: "The H4DM error message of the Replicator Service.",
				Computed:    true,
			},
			"h4dm_error_args": {
				Type:        schema.TypeList,
				Description: "The H4DM error arguments of the Replicator Service.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"h4dm_error_stacktrace": {
				Type:        schema.TypeString,
				Description: "The H4DM error stacktrace of the Replicator Service.",
				Computed:    true,
			},
			"ls_error_code": {
				Type:        schema.TypeString,
				Description: "The lookup service error code of the Replicator Service.",
				Computed:    true,
			},
			"ls_error_msg": {
				Type:        schema.TypeString,
				Description: "The lookup service error message of the Replicator Service.",
				Computed:    true,
			},
			"ls_error_args": {
				Type:        schema.TypeList,
				Description: "The lookup service error arguments of the Replicator Service.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"ls_error_stacktrace": {
				Type:        schema.TypeString,
				Description: "The lookup service error stacktrace of the Replicator Service.",
				Computed:    true,
			},
			"db_error_code": {
				Type:        schema.TypeString,
				Description: "The database error code of the Replicator Service.",
				Computed:    true,
			},
			"db_error_msg": {
				Type:        schema.TypeString,
				Description: "The database error message of the Replicator Service.",
				Computed:    true,
			},
			"db_error_args": {
				Type:        schema.TypeList,
				Description: "The database error arguments of the Replicator Service.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"db_error_stacktrace": {
				Type:        schema.TypeString,
				Description: "The database error stacktrace of the Replicator Service.",
				Computed:    true,
			},
			"ntp_error_code": {
				Type:        schema.TypeString,
				Description: "The NTP error code of the Replicator Service.",
				Computed:    true,
			},
			"ntp_error_msg": {
				Type:        schema.TypeString,
				Description: "The NTP error message of the Replicator Service.",
				Computed:    true,
			},
			"ntp_error_args": {
				Type:        schema.TypeList,
				Description: "The NTP error arguments of the Replicator Service.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"ntp_error_stacktrace": {
				Type:        schema.TypeString,
				Description: "The NTP error stacktrace of the Replicator Service.",
				Computed:    true,
			},
			"offline_managers_ids": {
				Type:        schema.TypeList,
				Description: "A list of the offline managers IDs of the Replicator Service.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"online_managers_ids": {
				Type:        schema.TypeList,
				Description: "A list of the online managers IDs of the Replicator Service.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceVcdaReplicatorHealthRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	return getReplicatorHealthInfo(c, d)
}

func getReplicatorHealthInfo(c *Client, d *schema.ResourceData) diag.Diagnostics {
	var diags diag.Diagnostics

	replicatorID := d.Get("replicator_id").(string)

	health, err := getHealthTaskResult(c, d)
	if err != nil {
		return diag.FromErr(err)
	}

	localReplicatorsHealth, exists := health["localReplicatorsHealth"].([]interface{})
	if exists {
		localReplicator, err := findReplicator(replicatorID, localReplicatorsHealth)
		if err != nil {
			return diag.FromErr(err)
		}

		if err := setReplicatorInfoData(d, localReplicator); err != nil {
			return diag.FromErr(err)
		}
	} else {
		managerHealth, ok := health["managerHealth"].(map[string]interface{})
		if !ok {
			return diag.FromErr(fmt.Errorf("unexpected type for managerHealth"))
		}

		localReplicatorsHealth, ok = managerHealth["localReplicatorsHealth"].([]interface{})
		if !ok {
			return diag.FromErr(fmt.Errorf("unexpected type for localReplicatorsHealth"))
		}

		localReplicator, err := findReplicator(replicatorID, localReplicatorsHealth)
		if err != nil {
			return diag.FromErr(err)
		}

		if err := setReplicatorInfoData(d, localReplicator); err != nil {
			return diag.FromErr(err)
		}
	}

	return diags
}

func setManagersIDs(d *schema.ResourceData, managers []interface{}, isOnline bool) error {
	var managersIDs []string
	for _, manager := range managers {
		if manMap, ok := manager.(map[string]interface{}); ok {
			if managerID, ok := manMap["id"].(string); ok {
				managersIDs = append(managersIDs, managerID)
			}
		}
	}

	if isOnline {
		if err := d.Set("online_managers_ids", managersIDs); err != nil {
			return fmt.Errorf("error setting online_managers_ids field: %s", err)
		}
	} else {
		if err := d.Set("offline_managers_ids", managersIDs); err != nil {
			return fmt.Errorf("error setting offline_managers_ids field: %s", err)
		}
	}

	return nil
}

func findReplicator(replicatorID string, replicators []interface{}) (map[string]interface{}, error) {
	for _, repl := range replicators {
		if replicator, ok := repl.(map[string]interface{}); ok {
			if replID, ok := replicator["instanceId"].(string); ok {
				if replID == replicatorID {
					return replicator, nil
				}
			}
		}
	}

	return nil, fmt.Errorf("replicator with ID: %s was not found", replicatorID)
}

func setReplicatorInfoData(d *schema.ResourceData, replicator map[string]interface{}) error {
	err := setHealthInfoData(d, replicator)
	if err != nil {
		return err
	}

	if lwdError, ok := replicator["lwdError"].(map[string]interface{}); ok {
		if err := setErrorData(d, lwdError, "lwd"); err != nil {
			return err
		}
	}

	if hbrError, ok := replicator["hbrError"].(map[string]interface{}); ok {
		if err := setErrorData(d, hbrError, "hbr"); err != nil {
			return err
		}
	}

	if h4DmError, ok := replicator["h4dmError"].(map[string]interface{}); ok {
		if err := setErrorData(d, h4DmError, "h4dm"); err != nil {
			return err
		}
	}

	if offlineManagers, ok := replicator["offlineManagers"].([]interface{}); ok {
		if err := setManagersIDs(d, offlineManagers, false); err != nil {
			return err
		}
	}

	if onlineManagers, ok := replicator["onlineManagers"].([]interface{}); ok {
		if err := setManagersIDs(d, onlineManagers, true); err != nil {
			return err
		}
	}

	return nil
}
