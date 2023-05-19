terraform {
  required_providers {
    vcda = {
      source  = "vmware/vcda"
      version = ">=1.0"
    }
  }
}

provider "vcda" {
  vcda_ip        = var.cloud_appliance_management_ip
  local_user     = var.local_user
  local_password = var.local_password

  vsphere_user                 = var.vsphere_user
  vsphere_password             = var.vsphere_password
  vsphere_server               = var.vsphere_server
  vsphere_allow_unverified_ssl = true
}

# get the vm's cloud thumbprint
data "vcda_service_cert" "cloud_service_cert" {
  datacenter_id = var.cloud_vm_datacenter_id
  name          = var.cloud_vm_name
  type          = "cloud"
}

output "vcda_cloud_service_cert" {
  value = data.vcda_service_cert.cloud_service_cert.id
}

# get vm's manager thumbprint
data "vcda_service_cert" "manager_service_cert" {
  datacenter_id = var.cloud_vm_datacenter_id
  name          = var.cloud_vm_name
  type          = "manager"
}

output "vcda_manager_service_cert" {
  value = data.vcda_service_cert.manager_service_cert.id
}

// change cloud appliance password - either through new_password or password_file
resource "vcda_appliance_password" "cloud_appliance_password" {
  current_password = "vmware"
  //new_password     = var.local_password
  password_file    = "vcda-pass.txt"
  appliance_ip     = var.cloud_appliance_management_ip
  service_cert     = data.vcda_service_cert.cloud_service_cert.id
}

output "vcda_cloud_appliance_password_is_expired" {
  value = vcda_appliance_password.cloud_appliance_password.root_password_expired
}
output "vcda_cloud_appliance_password_seconds_until_expiration" {
  value = vcda_appliance_password.cloud_appliance_password.seconds_until_expiration
}

### Thumbprints:
# get VCD thumbprint (UNSAFE with address and port)
data "vcda_remote_services_thumbprint" "vcd_thumbprint" {
  depends_on = [vcda_appliance_password.cloud_appliance_password]
  //address    = var.vcd_address
  //port       = "443"
  pem_file   = "vcd-cert.pem"
}

output "vcda_vcd_thumbprint" {
  value = data.vcda_remote_services_thumbprint.vcd_thumbprint.*
}

# compute LS thumbprint (UNSAFE with address and port)
data "vcda_remote_services_thumbprint" "ls_thumbprint" {
  depends_on = [vcda_appliance_password.cloud_appliance_password]
  //address    = var.lookup_service_address
  //port       = "443"
  pem_file   = "ls-cert.pem"
}

output "vcda_ls_thumbprint" {
  value = data.vcda_remote_services_thumbprint.ls_thumbprint.*
}

// config
resource "vcda_cloud_director_replication_manager" "cloud_site" {
  service_cert              = data.vcda_service_cert.cloud_service_cert.id
  lookup_service_thumbprint = data.vcda_remote_services_thumbprint.ls_thumbprint.id
  vcd_thumbprint            = data.vcda_remote_services_thumbprint.vcd_thumbprint.id

  license_key      = var.license_key
  site_name        = var.site_name
  site_description = var.site_description

  public_endpoint_address = "vcda.pub"
  public_endpoint_port    = 443

  vcd_username = var.vcd_username
  vcd_password = var.vcd_password
  vcd_url      = var.vcd_url

  lookup_service_url = var.lookup_service_url
}

output "vcda_cloud_site" {
  value = vcda_cloud_director_replication_manager.cloud_site.*
}

# get first replicator vm thumbprint
data "vcda_service_cert" "replicator_service_cert" {
  datacenter_id = var.cloud_vm_datacenter_id
  name          = var.first_replicator_vm_name
  type          = "replicator"
}

output "vcda_replicator_service_cert" {
  value = data.vcda_service_cert.replicator_service_cert.id
}

# get first replicator thumbprint
data "vcda_remote_services_thumbprint" "replicator_thumbprint" {
  depends_on = [vcda_appliance_password.cloud_appliance_password]
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

output "vcda_replicator_appliance_password_is_expired" {
  value = vcda_appliance_password.replicator_appliance_password.root_password_expired
}

output "vcda_replicator_appliance_password_seconds_until_expiration" {
  value = vcda_appliance_password.replicator_appliance_password.seconds_until_expiration
}

resource "vcda_replicator" "add_replicator" {
  depends_on = [
    vcda_cloud_director_replication_manager.cloud_site,
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
  depends_on    = [vcda_appliance_password.cloud_appliance_password]
  datacenter_id = var.cloud_vm_datacenter_id
  name          = var.second_replicator_vm_name
  type          = "replicator"
}

output "vcda_second_replicator_service_cert" {
  value = data.vcda_service_cert.second_replicator_service_cert.id
}

# get second cloud replicator thumbprint
data "vcda_remote_services_thumbprint" "second_replicator_thumbprint" {
  depends_on = [vcda_appliance_password.cloud_appliance_password]
  //address    = var.second_replicator_management_ip
  //port       = "443"
  pem_file   = "rep2-cert.pem"
}

output "vcda_second_replicator_thumbprint" {
  value = data.vcda_remote_services_thumbprint.second_replicator_thumbprint.*
}

// change second replicator appliance password - either through new_password or password_file
resource "vcda_appliance_password" "second_replicator_appliance_password" {
  appliance_ip     = var.second_replicator_management_ip
  current_password = var.initial_appliance_password
  new_password     = var.replicator_root_password
  //password_file    = "vcda-pass.txt"
  service_cert     = data.vcda_service_cert.second_replicator_service_cert.id
}

output "vcda_second_replicator_appliance_password_is_expired" {
  value = vcda_appliance_password.second_replicator_appliance_password.root_password_expired
}

output "vcda_second_replicator_appliance_password_seconds_until_expiration" {
  value = vcda_appliance_password.second_replicator_appliance_password.seconds_until_expiration
}

resource "vcda_replicator" "add_second_replicator" {
  depends_on = [
    vcda_appliance_password.second_replicator_appliance_password,
    vcda_replicator.add_replicator
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

// tunnel
# get tunnel vm thumbprint
data "vcda_service_cert" "tunnel_service_cert" {
  datacenter_id = var.cloud_vm_datacenter_id
  name          = var.tunnel_vm_name
  type          = "tunnel"
}

output "vcda_tunnel_service_cert" {
  value = data.vcda_service_cert.tunnel_service_cert.id
}

// change tunnel appliance password - either through new_password or password_file
resource "vcda_appliance_password" "tunnel_appliance_password" {
  appliance_ip     = var.tunnel_management_ip
  current_password = var.initial_appliance_password
  new_password     = var.tunnel_root_password
  //password_file    = "vcda-pass.txt"
  service_cert     = data.vcda_service_cert.tunnel_service_cert.id
}

output "vcda_tunnel_appliance_password_is_expired" {
  value = vcda_appliance_password.tunnel_appliance_password.root_password_expired
}

output "vcav_tunnel_appliance_password_seconds_until_expiration" {
  value = vcda_appliance_password.tunnel_appliance_password.seconds_until_expiration
}

resource "vcda_tunnel" "add_tunnel" {
  depends_on = [
    vcda_cloud_director_replication_manager.cloud_site,
    vcda_appliance_password.tunnel_appliance_password
  ]

  service_cert = data.vcda_service_cert.cloud_service_cert.id

  url           = var.tunnel_url
  root_password = var.tunnel_root_password
  certificate   = data.vcda_service_cert.tunnel_service_cert.id
}

output "vcda_add_tunnel" {
  value = vcda_tunnel.add_tunnel.*
}
