// Copyright (c) 2023-2024 Broadcom. All Rights Reserved.
// Broadcom Confidential. The term "Broadcom" refers to Broadcom Inc.
// and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vcda

const (
	VcdaAuthTokenHeader = "X-VCAV-Auth"
	ContentTypeHeader   = "Content-Type"
	ConfigSecretHeader  = "Config-Secret"
	AcceptHeader        = "Accept"
	UserAgent           = "User-Agent"

	APIVersion             = "h4-v4.7"
	ContentTypeHeaderValue = "application/json"
	AcceptHeaderValue      = "application/vnd.vmware." + APIVersion + "+json;charset=UTF-8"
	UserType               = "localUser"
	UserAgentValue         = "vcda-terraform-provider/" + APIVersion

	ManagerCertExtraConfigKey    = "guestinfo.manager.certificate"
	CloudCertExtraConfigKey      = "guestinfo.cloud.certificate"
	TunnelCertExtraConfigKey     = "guestinfo.tunnel.certificate"
	ReplicatorCertExtraConfigKey = "guestinfo.replicator.certificate"

	VcdaIP                    = "VCDA_IP"
	LocalUser                 = "LOCAL_USER"
	LocalPassword             = "LOCAL_PASSWORD"
	VsphereUser               = "VSPHERE_USER"
	VspherePassword           = "VSPHERE_PASSWORD"
	VsphereServer             = "VSPHERE_SERVER"
	VsphereAllowUnverifiedSSL = "VSPHERE_ALLOW_UNVERIFIED_SSL"
	DatacenterID              = "DC_ID"
	CloudVMName               = "CLOUD_VM_NAME"
	ManagerVMName             = "MANAGER_VM_NAME"
	RemoteManagerVMName       = "REMOTE_MANAGER_VM_NAME"
	ReplicatorVMName          = "REPLICATOR_VM_NAME"
	TunnelVMName              = "TUNNEL_VM_NAME"
	RootPassword              = "ROOT_PASSWORD"
	CurrentPassword           = "CURRENT_PASSWORD"
	NewPassword               = "NEW_PASSWORD"
	LicenseKey                = "LICENSE_KEY"
	VcloudDirectorUsername    = "VCD_USERNAME"
	VcloudDirectorPassword    = "VCD_PASSWORD"
	VcloudDirectorAddress     = "VCD_ADDRESS"
	LookupServiceAddress      = "LS_ADDRESS"
	ReplicatorAddress         = "REPLICATOR_ADDRESS"
	SsoUser                   = "SSO_USER"
	SsoPassword               = "SSO_PASSWORD"
	TunnelAddress             = "TUNNEL_ADDRESS"
	ManagerAddress            = "MANAGER_ADDRESS"
	RemoteManagerAddress      = "REMOTE_MANAGER_ADDRESS"
)
