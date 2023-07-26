# Provider configuration
variable "cloud_appliance_management_ip" { default = "" }
variable "first_replicator_management_ip" { default = "" }
variable "second_replicator_management_ip" { default = "" }
variable "tunnel_management_ip" { default = "" }

variable "vcd_address" { default = "" }
variable "lookup_service_address" { default = "" }

variable "local_user" { default = "" }
variable "local_password" { default = "" }
variable "initial_appliance_password" { default = "" }

variable "vsphere_server" { default = "" }
variable "vsphere_user" { default = "" }
variable "vsphere_password" { default = "" }

# Cloud Director Replication Manager configuration
variable "cloud_vm_datacenter_id" { default = "" }
variable "cloud_vm_name" { default = "" }
variable "first_replicator_vm_name" { default = "" }
variable "second_replicator_vm_name" { default = "" }
variable "tunnel_vm_name" { default = "" }

variable "license_key" { default = "" }
variable "site_name" { default = "" }
variable "site_description" { default = "" }
variable "vcd_url" { default = "" }
variable "vcd_password" { default = "" }
variable "vcd_username" { default = "" }
variable "lookup_service_url" { default = "" }

# Replicator configuration
variable "replicator_url" { default = "" }
variable "second_replicator_url" { default = "" }
variable "replicator_lookup_service_url" { default = "" }
variable "replicator_sso_user" { default = "" }
variable "replicator_sso_password" { default = "" }
variable "replicator_root_password" { default = "" }
variable "replicator_owner" { default = "" }

# Tunnel configuration
variable "tunnel_url" { default = "" }
variable "second_tunnel_url" { default = "" }
variable "tunnel_root_password" { default = "" }

# Pairing
variable "remote_cloud_address" { default = "" }
variable "remote_cloud_site_name" { default = "" }
variable "remote_cloud_endpoint_url" { default = "" }
