terraform {
  required_providers {
    vcda = {
      source  = "vmware/vcda"
      version = ">=1.0"
    }
  }
}

provider "vcda" {
  vcda_ip        = var.manager_appliance_management_ip
  local_user     = var.local_user
  local_password = var.local_password

  vsphere_user                 = var.vsphere_user
  vsphere_password             = var.vsphere_password
  vsphere_server               = var.vsphere_server
  vsphere_allow_unverified_ssl = true
}

# get the manager thumbprint
data "vcda_service_cert" "manager_service_cert" {
  datacenter_id = var.manager_vm_datacenter_id
  name          = var.manager_vm_name
  type          = "manager"
}

output "vcda_manager_service_cert" {
  value = data.vcda_service_cert.manager_service_cert.id
}

// change manager appliance password - either through new_password or password_file
resource "vcda_appliance_password" "manager_appliance_password" {
  current_password = var.initial_appliance_password
  new_password     = var.local_password
  //password_file    = "vcda-pass.txt"
  appliance_ip     = var.manager_appliance_management_ip
  service_cert     = data.vcda_service_cert.manager_service_cert.id
}

output "vcda_appliance_password_is_expired" {
  value = vcda_appliance_password.manager_appliance_password.root_password_expired
}

output "vcda_appliance_password_seconds_until_expiration" {
  value = vcda_appliance_password.manager_appliance_password.seconds_until_expiration
}

# get LS thumbprint
data "vcda_remote_services_thumbprint" "ls_thumbprint" {
  depends_on = [vcda_appliance_password.manager_appliance_password]
  //address    = var.lookup_service_address
  //port       = "443"
  pem_file   = "ls-cert.pem"
}

output "vcda_ls_thumbprint" {
  value = data.vcda_remote_services_thumbprint.ls_thumbprint.*
}

resource "vcda_vcenter_replication_manager" "manager_site" {
  service_cert              = data.vcda_service_cert.manager_service_cert.id
  lookup_service_thumbprint = data.vcda_remote_services_thumbprint.ls_thumbprint.id

  license_key        = var.license_key
  site_name          = var.site_name
  lookup_service_url = var.lookup_service_url
  sso_user           = var.replicator_sso_user
  sso_password       = var.replicator_sso_password
}

output "vcda_manager_site" {
  value = vcda_vcenter_replication_manager.manager_site.*
}

// add two replicators
# get first replicator vm thumbprint
data "vcda_service_cert" "replicator_service_cert" {
  datacenter_id = var.manager_vm_datacenter_id
  name          = var.first_replicator_vm_name
  type          = "replicator"
}

output "vcda_replicator_api_thumbprint" {
  value = data.vcda_service_cert.replicator_service_cert.id
}

# get first replicator thumbprint
data "vcda_remote_services_thumbprint" "replicator_thumbprint" {
  depends_on = [vcda_appliance_password.manager_appliance_password]
  //address    = var.first_replicator_management_ip
  //port       = "443"
  pem_file   = "rep-cert.pem"
}

output "vcda_replicator_thumbprint" {
  value = data.vcda_remote_services_thumbprint.replicator_thumbprint.*
}

// change first replicator appliance password - either through new_password or password_file
resource "vcda_appliance_password" "replicator_appliance_password" {
  appliance_ip     = var.first_replicator_management_ip
  current_password = var.initial_appliance_password
  new_password     = var.replicator_root_password
  //password_file    = "vcda-pass.txt"
  service_cert     = data.vcda_service_cert.replicator_service_cert.id
}

output "vcda_replicator_appliance_password" {
  value     = vcda_appliance_password.replicator_appliance_password.*
  sensitive = true
}

resource "vcda_replicator" "add_replicator" {
  depends_on = [
    vcda_vcenter_replication_manager.manager_site,
    vcda_appliance_password.replicator_appliance_password
  ]

  lookup_service_url = var.replicator_lookup_service_url
  api_url            = var.replicator_url
  sso_user           = var.replicator_sso_user
  sso_password       = var.replicator_sso_password
  root_password      = var.replicator_root_password
  owner              = var.replicator_owner
  site_name          = var.site_name

  api_thumbprint            = data.vcda_remote_services_thumbprint.replicator_thumbprint.id
  service_cert              = data.vcda_service_cert.manager_service_cert.id
  lookup_service_thumbprint = data.vcda_remote_services_thumbprint.ls_thumbprint.id
}

output "vcda_add_replicator" {
  value = vcda_replicator.add_replicator.*
}

# get second replicator vm thumbprint
data "vcda_service_cert" "second_replicator_service_cert" {
  datacenter_id = var.manager_vm_datacenter_id
  name          = var.second_replicator_vm_name
  type          = "replicator"
}

output "vcda_second_replicator_service_cert" {
  value = data.vcda_service_cert.second_replicator_service_cert.id
}

# get second replicator thumbprint
data "vcda_remote_services_thumbprint" "second_replicator_thumbprint" {
  depends_on = [vcda_appliance_password.manager_appliance_password]
  //address    = var.second_replicator_management_ip
  //port       = "443"
  pem_file   = "rep2-cert.pem"
}

output "vcda_second_replicator_thumbprint" {
  value = data.vcda_remote_services_thumbprint.second_replicator_thumbprint.id
}

// change second replicator appliance password - either through new_password or password_file
resource "vcda_appliance_password" "second_replicator_appliance_password" {
  appliance_ip     = var.second_replicator_management_ip
  current_password = var.initial_appliance_password
  new_password     = var.replicator_root_password
  //password_file        = "vcda-pass.txt"
  service_cert     = data.vcda_service_cert.second_replicator_service_cert.id
}

output "vcda_second_replicator_appliance_password" {
  value     = vcda_appliance_password.second_replicator_appliance_password.*
  sensitive = true
}

resource "vcda_replicator" "add_second_replicator" {
  depends_on = [
    vcda_vcenter_replication_manager.manager_site,
    vcda_appliance_password.second_replicator_appliance_password
  ]

  lookup_service_url = var.replicator_lookup_service_url
  api_url            = var.second_replicator_url
  sso_user           = var.replicator_sso_user
  sso_password       = var.replicator_sso_password
  root_password      = var.replicator_root_password
  owner              = var.replicator_owner
  site_name          = var.site_name

  api_thumbprint            = data.vcda_remote_services_thumbprint.second_replicator_thumbprint.id
  service_cert              = data.vcda_service_cert.manager_service_cert.id
  lookup_service_thumbprint = data.vcda_remote_services_thumbprint.ls_thumbprint.id
}

output "vcda_add_second_replicator" {
  value = vcda_replicator.add_second_replicator.*
}

# remote site thumbprint
data "vcda_remote_services_thumbprint" "remote_cloud_thumbprint" {
  address = var.remote_site_address
  port    = "443"
  //pem_file   = "vcd-cert.pem"
}

output "vcda_remote_cloud_thumbprint" {
  value = data.vcda_remote_services_thumbprint.remote_cloud_thumbprint.*
}

# Pair site
resource "vcda_pair_site" "pair_site" {
  depends_on = [
    vcda_replicator.add_second_replicator
  ]

  service_cert   = data.vcda_service_cert.manager_service_cert.id
  api_thumbprint = data.vcda_remote_services_thumbprint.remote_cloud_thumbprint.id

  api_url             = var.remote_site_endpoint_url
  pairing_description = "pair site"
}

output "vcda_pair_site1" {
  value = vcda_pair_site.pair_site.*
}