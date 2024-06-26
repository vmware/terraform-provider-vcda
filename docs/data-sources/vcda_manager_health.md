---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "vcda_manager_health Data Source - terraform-provider-for-vmware-cloud-director-availability"
subcategory: ""
description: |-
  VMware Cloud Director Availability Manager Health data source.
---

# vcda_manager_health (Data Source)

The manager health data source obtains the current health info of either:

- already deployed/configured Manager service of the Cloud Director Replication Management Appliance

- already deployed/configured vCenter Replication Management Appliance

## Example Usage

### Health info of the Manager service of an already deployed/configured Cloud Director Replication Management Appliance

```terraform
data "vcda_manager_health" "manager_health" {
  service_cert = data.vcda_service_cert.cloud_service_cert.id
  manager_id   = data.vcda_cloud_health.cloud_health.manager_id
}
```

### Health info of an already deployed/configured vCenter Replication Management Appliance

```terraform
data "vcda_manager_health" "manager_health" {
  service_cert = data.vcda_service_cert.manager_service_cert.id
}
```

<!-- schema generated by tfplugindocs -->

### Required

- `service_cert` (String)  The certificate of the Cloud Director/vCenter Replication Manager Service.

### Optional

- `manager_id` (String)  The cloud manager instance id. **NOTE:** only required for the Cloud Director/Manager Service
  health info. It could be set explicitly or obtained from the `vcda_cloud_health` data source.

### Read-Only

- `id` (String) The health info task ID.
- `product_name` (String) The product name of the Cloud Director/vCenter Replication Manager Service.
- `build_version` (String) The build version of the Cloud Director/vCenter Replication Manager Service.
- `build_date` (String) The build date of the Cloud Director/vCenter Replication Manager Service.
- `instance_id` (String) The instance ID of the Cloud Director/vCenter Replication Manager Service.
- `runtime_id` (String) The runtime ID of the Cloud Director/vCenter Replication Manager Service.
- `current_time` (Number) The current time of the Cloud Director/vCenter Replication Manager Service.
- `address` (String) The address of the Cloud Director/vCenter Replication Manager Service.
- `service_boot_timestamp` (Number) The service boot timestamp of the Cloud Director/vCenter Replication Manager
  Service.
- `appliance_boot_timestamp` (Number) The appliance boot timestamp of the Cloud Director/vCenter Replication Manager
  Service.
- `disk_usage` (Map) The disk usage of the Cloud Director/vCenter Replication Manager Service.
- `local_replicators_ls_mismatch_error_code` (String) The local replicators lookup service mismatch error code of the
  Cloud Director/vCenter Replication Manager Service.
- `local_replicators_ls_mismatch_error_msg` (String) The local replicators lookup service mismatch error message of the
  Cloud Director/vCenter Replication Manager Service.
- `local_replicators_ls_mismatch_error_args` (List) The local replicators lookup service mismatch error arguments of the
  Cloud Director/vCenter Replication Manager Service.
- `local_replicators_ls_mismatch_error_stacktrace` (String) The local replicators lookup service mismatch error
  stacktrace of the Cloud Director/vCenter Replication Manager Service.
- `sso_admin_error_code` (String) The sso admin error code of the Cloud Director/vCenter Replication Manager Service.
- `sso_admin_error_msg` (String) The sso admin error message of the Cloud Director/vCenter Replication Manager Service.
- `sso_admin_error_args` (List) The sso admin error arguments of the Cloud Director/vCenter Replication Manager Service.
- `sso_admin_error_stacktrace` (String) The sso admin error stacktrace of the Cloud Director/vCenter Replication Manager
  Service.
- `ls_error_code` (String) The lookup service error code of the Cloud Director/vCenter Replication Manager Service.
- `ls_error_msg` (String) The lookup service error message of the Cloud Director/vCenter Replication Manager Service.
- `ls_error_args` (List) The lookup service error arguments of the Cloud Director/vCenter Replication Manager Service.
- `ls_error_stacktrace` (String) The lookup service error stacktrace of the Cloud Director/vCenter Replication Manager
  Service.
- `db_error_code` (String) The database error code of the Cloud Director/vCenter Replication Manager Service.
- `db_error_msg` (String) The database error message of the Cloud Director/vCenter Replication Manager Service.
- `db_error_args` (List) The database error arguments of the Cloud Director/vCenter Replication Manager Service.
- `db_error_stacktrace` (String) The database error stacktrace of the Cloud Director/vCenter Replication Manager
  Service.
- `ntp_error_code` (String) The NTP error code of the Cloud Director/vCenter Replication Manager Service.
- `ntp_error_msg` (String) The NTP error message of the Cloud Director/vCenter Replication Manager Service.
- `ntp_error_args` (List) The NTP error arguments of the Cloud Director/vCenter Replication Manager Service.
- `ntp_error_stacktrace` (String) The NTP error stacktrace of the Cloud Director/vCenter Replication Manager Service.
- `offline_replicators_ids` (List) A list of the offline replicators IDs of the Cloud Director/vCenter Replication
  Manager Service.
- `online_replicators_ids` (List) A list of the online replicators IDs of the Cloud Director/vCenter Replication
  Manager Service.
- `local_replicators_ids` (List) A list of the local replicators IDs of the Cloud Director/vCenter Replication Manager
  Service.
- `tunnels_ids` (List) A list of the tunnels IDs of the Cloud Director/vCenter Replication Manager Service.