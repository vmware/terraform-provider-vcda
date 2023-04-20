# Provider configuration
variable "manager_appliance_management_ip" { default = "" }
variable "first_replicator_management_ip" { default = "" }
variable "second_replicator_management_ip" { default = "" }

variable "lookup_service_address" { default = "" }

variable "local_user" { default = "" }
variable "local_password" { default = "" }
variable "initial_appliance_password" { default = "" }

variable "vsphere_server" { default = "" }
variable "vsphere_user" { default = "" }
variable "vsphere_password" { default = "" }

# Vcenter Replication Manager configuration
variable "manager_vm_datacenter_id" { default = "" }
variable "manager_vm_name" { default = "" }
variable "first_replicator_vm_name" { default = "" }
variable "second_replicator_vm_name" { default = "" }

variable "license_key" { default = "" }
variable "site_name" { default = "" }
variable "lookup_service_url" { default = "" }

# Replicator configuration
variable "replicator_url" { default = "" }
variable "second_replicator_url" { default = "" }
variable "replicator_lookup_service_url" { default = "" }
variable "replicator_sso_user" { default = "" }
variable "replicator_sso_password" { default = "" }
variable "replicator_root_password" { default = "" }
variable "replicator_owner" { default = "" }
