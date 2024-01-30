// Copyright (c) 2024 Broadcom. All Rights Reserved.
// Broadcom Confidential. The term "Broadcom" refers to Broadcom Inc.
// and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vcda

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"time"
)

func dataSourceVcdaCloudHealth() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVcdaCloudHealthRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"service_cert": {
				Type:        schema.TypeString,
				Description: "The service certificate.",
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
				Description: "The product name of the Cloud Director Replication Manager Service.",
				Computed:    true,
			},
			"build_version": {
				Type:        schema.TypeString,
				Description: "The build version of the Cloud Director Replication Manager Service.",
				Computed:    true,
			},
			"build_date": {
				Type:        schema.TypeFloat,
				Description: "The build date of the Cloud Director Replication Manager Service.",
				Computed:    true,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Description: "The instance ID of the Cloud Director Replication Manager Service.",
				Computed:    true,
			},
			"runtime_id": {
				Type:        schema.TypeString,
				Description: "The runtime ID of the Cloud Director Replication Manager Service.",
				Computed:    true,
			},
			"current_time": {
				Type:        schema.TypeFloat,
				Description: "The current time of the Cloud Director Replication Manager Service.",
				Computed:    true,
			},
			"address": {
				Type:        schema.TypeString,
				Description: "The address of the Cloud Director Replication Manager Service.",
				Computed:    true,
			},
			"service_boot_timestamp": {
				Type:        schema.TypeInt,
				Description: "The service boot timestamp of the Cloud Director Replication Manager Service.",
				Computed:    true,
			},
			"appliance_boot_timestamp": {
				Type:        schema.TypeFloat,
				Description: "The appliance boot timestamp of the Cloud Director Replication Manager Service.",
				Computed:    true,
			},
			"disk_usage": {
				Type:        schema.TypeMap,
				Description: "The disk usage of the Cloud Director Replication Manager Service.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"vcd_error_code": {
				Type:        schema.TypeString,
				Description: "The VCD error code of the Cloud Director Replication Manager Service.",
				Computed:    true,
			},
			"vcd_error_msg": {
				Type:        schema.TypeString,
				Description: "The VCD error message of the Cloud Director Replication Manager Service.",
				Computed:    true,
			},
			"vcd_error_args": {
				Type:        schema.TypeList,
				Description: "The VCD error arguments of the Cloud Director Replication Manager Service.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"vcd_error_stacktrace": {
				Type:        schema.TypeString,
				Description: "The VCD error stacktrace of the Cloud Director Replication Manager Service.",
				Computed:    true,
			},
			"manager_error_code": {
				Type:        schema.TypeString,
				Description: "The cloud manager error code of the Cloud Director Replication Manager Service.",
				Computed:    true,
			},
			"manager_error_msg": {
				Type:        schema.TypeString,
				Description: "The cloud manager error message of the Cloud Director Replication Manager Service.",
				Computed:    true,
			},
			"manager_error_args": {
				Type:        schema.TypeList,
				Description: "The cloud manager error arguments of the Cloud Director Replication Manager Service.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"manager_error_stacktrace": {
				Type:        schema.TypeString,
				Description: "The cloud manager error stacktrace of the Cloud Director Replication Manager Service.",
				Computed:    true,
			},
			"ls_error_code": {
				Type:        schema.TypeString,
				Description: "The lookup service error code of the Cloud Director Replication Manager Service.",
				Computed:    true,
			},
			"ls_error_msg": {
				Type:        schema.TypeString,
				Description: "The lookup service error message of the Cloud Director Replication Manager Service.",
				Computed:    true,
			},
			"ls_error_args": {
				Type:        schema.TypeList,
				Description: "The lookup service error arguments of the Cloud Director Replication Manager Service.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"ls_error_stacktrace": {
				Type:        schema.TypeString,
				Description: "The lookup service error stacktrace of the Cloud Director Replication Manager Service.",
				Computed:    true,
			},
			"db_error_code": {
				Type:        schema.TypeString,
				Description: "The database error code of the Cloud Director Replication Manager Service.",
				Computed:    true,
			},
			"db_error_msg": {
				Type:        schema.TypeString,
				Description: "The database error message of the Cloud Director Replication Manager Service.",
				Computed:    true,
			},
			"db_error_args": {
				Type:        schema.TypeList,
				Description: "The database error arguments of the Cloud Director Replication Manager Service.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"db_error_stacktrace": {
				Type:        schema.TypeString,
				Description: "The database error stacktrace of the Cloud Director Replication Manager Service.",
				Computed:    true,
			},
			"ntp_error_code": {
				Type:        schema.TypeString,
				Description: "The NTP error code of the Cloud Director Replication Manager Service.",
				Computed:    true,
			},
			"ntp_error_msg": {
				Type:        schema.TypeString,
				Description: "The NTP error message of the Cloud Director Replication Manager Service.",
				Computed:    true,
			},
			"ntp_error_args": {
				Type:        schema.TypeList,
				Description: "The NTP error arguments of the Cloud Director Replication Manager Service.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"ntp_error_stacktrace": {
				Type:        schema.TypeString,
				Description: "The NTP error stacktrace of the Cloud Director Replication Manager Service.",
				Computed:    true,
			},
			"tunnels_ids": {
				Type:        schema.TypeList,
				Description: "A list of the tunnels IDs of the Cloud Director/vCenter Replication Manager Service.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"manager_id": {
				Type:        schema.TypeString,
				Description: "The cloud manager ID of the Cloud Director Replication Manager Service.",
				Computed:    true,
			},
		},
	}
}

func dataSourceVcdaCloudHealthRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	return getCloudHealthInfo(c, d)
}

func retryHealthTask(c *Client, d *schema.ResourceData, serviceCert string, taskID *string) error {
	err := retry.RetryContext(context.Background(), d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		task, err := c.getTask(serviceCert, *taskID)

		if err != nil {
			return retry.NonRetryableError(err)
		}

		if task.State == "FAILED" {
			return retry.NonRetryableError(fmt.Errorf("health task failed with Code:" + task.Error.Code + ", Msg: " + task.Error.Msg))
		} else if task.State == "RUNNING" {
			return retry.RetryableError(fmt.Errorf("expected health task to be completed but was in state %s", task.State))
		} else if task.State == "QUEUED" {
			return retry.RetryableError(fmt.Errorf("expected health task to be completed but was in state %s", task.State))
		}

		return nil
	})
	return err
}

func getHealthTaskResult(c *Client, d *schema.ResourceData) (map[string]interface{}, error) {
	serviceCert := d.Get("service_cert").(string)
	taskID := d.Get("id").(string)

	task, err := c.getTask(serviceCert, taskID)
	if err != nil {
		return nil, err
	}

	health, ok := task.Result.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected type for taskResult.Result")
	}

	return health, nil
}

func getCloudHealthInfo(c *Client, d *schema.ResourceData) diag.Diagnostics {
	var diags diag.Diagnostics

	health, err := getHealthTaskResult(c, d)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setCloudHealthInfoData(d, health); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func setCloudHealthInfoData(d *schema.ResourceData, cloudHealth map[string]interface{}) error {
	err := setHealthInfoData(d, cloudHealth)
	if err != nil {
		return err
	}

	if vcdError, ok := cloudHealth["vcdError"].(map[string]interface{}); ok {
		if err := setErrorData(d, vcdError, "vcd"); err != nil {
			return err
		}
	}

	if managerError, ok := cloudHealth["managerError"].(map[string]interface{}); ok {
		if err := setErrorData(d, managerError, "manager"); err != nil {
			return err
		}
	}

	if tunnelConnectivity, ok := cloudHealth["tunnelConnectivity"].([]interface{}); ok {
		if err := setTunnelIDs(d, tunnelConnectivity); err != nil {
			return err
		}
	}

	if managerHealth, ok := cloudHealth["managerHealth"].(map[string]interface{}); ok {
		if err := setManagerID(d, managerHealth); err != nil {
			return err
		}
	}

	return nil
}

func setTunnelIDs(d *schema.ResourceData, tunnelsList []interface{}) error {
	var tunnelsIDs []string
	for _, tunnelConnectivity := range tunnelsList {
		if tunMap, ok := tunnelConnectivity.(map[string]interface{}); ok {
			tunnelService := tunMap["tunnelService"].(map[string]interface{})
			if tunnelID, ok := tunnelService["id"].(string); ok {
				tunnelsIDs = append(tunnelsIDs, tunnelID)
			}
		}
	}

	if err := d.Set("tunnels_ids", tunnelsIDs); err != nil {
		return fmt.Errorf("error setting tunnels_ids field: %s", err)
	}

	return nil
}

func setManagerID(d *schema.ResourceData, managerHealth map[string]interface{}) error {
	if managerID, ok := managerHealth["instanceId"].(string); ok {
		if err := d.Set("manager_id", managerID); err != nil {
			return fmt.Errorf("error setting manager_id field: %s", err)
		}
	}

	return nil
}

func setErrorData(d *schema.ResourceData, errorMap map[string]interface{}, errorPrefix string) error {
	if err := d.Set(errorPrefix+"_error_code", errorMap["code"]); err != nil {
		return fmt.Errorf("error setting "+errorPrefix+"_error_code field: %s", err)
	}

	if err := d.Set(errorPrefix+"_error_msg", errorMap["msg"]); err != nil {
		return fmt.Errorf("error setting "+errorPrefix+"_error_msg field: %s", err)
	}

	if err := d.Set(errorPrefix+"_error_args", errorMap["args"]); err != nil {
		return fmt.Errorf("error setting "+errorPrefix+"_error_args field: %s", err)
	}

	if err := d.Set(errorPrefix+"_error_stacktrace", errorMap["stacktrace"]); err != nil {
		return fmt.Errorf("error setting "+errorPrefix+"_error_stacktrace field: %s", err)
	}

	return nil
}

func setHealthInfoData(d *schema.ResourceData, health map[string]interface{}) error {
	if err := d.Set("product_name", health["productName"]); err != nil {
		return fmt.Errorf("error setting product_name field: %s", err)
	}

	if err := d.Set("build_version", health["buildVersion"]); err != nil {
		return fmt.Errorf("error setting build_version field: %s", err)
	}

	if err := d.Set("build_date", health["buildDate"]); err != nil {
		return fmt.Errorf("error setting build_date field: %s", err)
	}

	if err := d.Set("instance_id", health["instanceId"]); err != nil {
		return fmt.Errorf("error setting instance_id field: %s", err)
	}

	if err := d.Set("runtime_id", health["runtimeId"]); err != nil {
		return fmt.Errorf("error setting runtime_id field: %s", err)
	}

	if err := d.Set("current_time", health["currentTime"]); err != nil {
		return fmt.Errorf("error setting current_time field: %s", err)
	}

	if err := d.Set("address", health["address"]); err != nil {
		return fmt.Errorf("error setting address field: %s", err)
	}

	if err := d.Set("service_boot_timestamp", health["serviceBootTimestamp"]); err != nil {
		return fmt.Errorf("error setting service_boot_timestamp field: %s", err)
	}

	if err := d.Set("appliance_boot_timestamp", health["applianceBootTimestamp"]); err != nil {
		return fmt.Errorf("error setting appliance_boot_timestamp field: %s", err)
	}

	if diskData, ok := health["diskUsage"].(map[string]interface{}); ok {
		if err := d.Set("disk_usage", diskData); err != nil {
			return err
		}
	}

	if lsError, ok := health["lsError"].(map[string]interface{}); ok {
		if err := setErrorData(d, lsError, "ls"); err != nil {
			return err
		}
	}

	if dbError, ok := health["dbError"].(map[string]interface{}); ok {
		if err := setErrorData(d, dbError, "db"); err != nil {
			return err
		}
	}

	if ntpError, ok := health["ntpError"].(map[string]interface{}); ok {
		if err := setErrorData(d, ntpError, "ntp"); err != nil {
			return err
		}
	}

	return nil
}
