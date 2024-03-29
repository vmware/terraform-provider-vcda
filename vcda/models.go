// Copyright (c) 2023-2024 Broadcom. All Rights Reserved.
// Broadcom Confidential. The term "Broadcom" refers to Broadcom Inc.
// and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vcda

type AuthTokenData struct {
	Type          string `json:"type"`
	LocalUser     string `json:"localUser"`
	LocalPassword string `json:"localPassword"`
}

type PasswordData struct {
	RootPassword string `json:"rootPassword"`
}

type PasswordExpiration struct {
	RootPasswordExpired    bool  `json:"rootPasswordExpired"`
	SecondsUntilExpiration int64 `json:"secondsUntilExpiration"`
}

type LicenseData struct {
	LicenseKey string `json:"key"`
}

type License struct {
	LicenseKey     string `json:"key"`
	IsLicensed     bool   `json:"isLicensed"`
	ExpirationDate int64  `json:"expirationDate"`
}

type SiteData struct {
	Site string `json:"site"`
}

type EndpointData struct {
	APIAddress        interface{} `json:"apiAddress"`
	APIPort           int         `json:"apiPort"`
	APIPublicAddress  string      `json:"apiPublicAddress"`
	APIPublicPort     int         `json:"apiPublicPort"`
	MgmtAddress       interface{} `json:"mgmtAddress"`
	MgmtPort          int         `json:"mgmtPort"`
	MgmtPublicAddress interface{} `json:"mgmtPublicAddress"`
	MgmtPublicPort    interface{} `json:"mgmtPublicPort"`
}

type Endpoints struct {
	Configured EndpointConfig `json:"configured"`
	Effective  EndpointConfig `json:"effective"`
}

type EndpointConfig struct {
	MgmtAddress       *string `json:"mgmtAddress"`
	MgmtPort          int64   `json:"mgmtPort"`
	MgmtPublicAddress *string `json:"mgmtPublicAddress"`
	MgmtPublicPort    *int64  `json:"mgmtPublicPort"`
	APIAddress        *string `json:"apiAddress"`
	APIPort           int64   `json:"apiPort"`
	APIPublicAddress  string  `json:"apiPublicAddress"`
	APIPublicPort     int64   `json:"apiPublicPort"`
}

type CloudSiteData struct {
	LocalSite            string `json:"localSite"`
	LocalSiteDescription string `json:"localSiteDescription"`
}

type SiteConfig struct {
	ID                string `json:"id"`
	Site              string `json:"site"`
	LsURL             string `json:"lsUrl"`
	LsThumbprint      string `json:"lsThumbprint"`
	TunnelURL         string `json:"tunnelUrl"`
	TunnelCertificate string `json:"tunnelCertificate"`
}

type LookupServiceData struct {
	URL        string `json:"url"`
	Thumbprint string `json:"thumbprint"`
}

type ManagerLookupServiceData struct {
	URL                 string              `json:"url"`
	Thumbprint          string              `json:"thumbprint"`
	SsoAdminCredentials SsoAdminCredentials `json:"ssoAdminCredentials"`
}

type SsoAdminCredentials struct {
	SsoUser     string `json:"ssoUser"`
	SsoPassword string `json:"ssoPassword"`
}

type LookupService struct {
	LsURL        string `json:"lsUrl"`
	LsThumbprint string `json:"lsThumbprint"`
}

type Replicator struct {
	ID                  string      `json:"id"`
	Owner               string      `json:"owner"`
	Site                string      `json:"site"`
	Description         string      `json:"description"`
	APIURL              string      `json:"apiUrl"`
	CERTThumbprint      string      `json:"certThumbprint"`
	PairingCookie       interface{} `json:"pairingCookie"`
	State               State       `json:"state"`
	IsInMaintenanceMode bool        `json:"isInMaintenanceMode"`
	APIVersion          string      `json:"apiVersion"`
	DataAddress         interface{} `json:"dataAddress"`
	BuildVersion        interface{} `json:"buildVersion"`
}

type State struct {
	IncomingCommError interface{} `json:"incomingCommError"`
	OutgoingCommError interface{} `json:"outgoingCommError"`
}

type CloudSiteConfig struct {
	ID                   string `json:"id"`
	LsURL                string `json:"lsUrl"`
	LsThumbprint         string `json:"lsThumbprint"`
	LocalSite            string `json:"localSite"`
	LocalSiteDescription string `json:"localSiteDescription"`
	VcdURL               string `json:"vcdUrl"`
	VcdThumbprint        string `json:"vcdThumbprint"`
	VcdUsername          string `json:"vcdUsername"`
	TunnelURL            string `json:"tunnelUrl"`
	TunnelCertificate    string `json:"tunnelCertificate"`
	IsCombined           bool   `json:"isCombined"`
}

type VcloudConfigData struct {
	VcdPassword   string `json:"vcdPassword"`
	VcdThumbprint string `json:"vcdThumbprint"`
	VcdURL        string `json:"vcdUrl"`
	VcdUsername   string `json:"vcdUsername"`
}

type ReplicatorLookupServiceData struct {
	LsURL         string `json:"lsUrl"`
	LsThumbprint  string `json:"lsThumbprint"`
	APIURL        string `json:"apiUrl"`
	APIThumbprint string `json:"apiThumbprint"`
	RootPassword  string `json:"rootPassword"`
}

type ReplicatorConfigData struct {
	APIURL        string `json:"apiUrl"`
	APIThumbprint string `json:"apiThumbprint"`
	RootPassword  string `json:"rootPassword"`
	SsoUser       string `json:"ssoUser"`
	SsoPassword   string `json:"ssoPassword"`
}

type ReplicatorData struct {
	Description  string               `json:"description"`
	Owner        string               `json:"owner"`
	Site         string               `json:"site"`
	ReplicatorID interface{}          `json:"replicatorId"`
	Details      ReplicatorConfigData `json:"details"`
}

type TunnelData struct {
	Certificate  string `json:"certificate"`
	RootPassword string `json:"rootPassword"`
	URL          string `json:"url"`
}

type Tunnels struct {
	Tunnels []TunnelConfig `json:"tunnels"`
}

type TunnelConfig struct {
	ID          string `json:"id"`
	URL         string `json:"url"`
	Certificate string `json:"certificate"`
}

type VspherePluginStatus struct {
	Status string `json:"status"`
}

type IsServiceConfigured struct {
	IsConfigured bool `json:"isConfigured"`
}

type PairCloudSiteData struct {
	APIThumbprint string `json:"apiThumbprint"`
	APIURL        string `json:"apiUrl"`
	Description   string `json:"description"`
	Site          string `json:"site"`
}

type PairVcenterSiteData struct {
	APIThumbprint string `json:"apiThumbprint"`
	APIURL        string `json:"apiUrl"`
	Description   string `json:"description"`
}

type Task struct {
	ID           string        `json:"id"`
	User         string        `json:"user"`
	WorkflowInfo interface{}   `json:"workflowInfo"`
	Progress     int64         `json:"progress"`
	State        string        `json:"state"`
	LastUpdated  int64         `json:"lastUpdated"`
	StartTime    int64         `json:"startTime"`
	EndTime      int64         `json:"endTime"`
	ResultType   string        `json:"resultType"`
	Result       interface{}   `json:"result"`
	Error        Error         `json:"error"`
	Warnings     []interface{} `json:"warnings"`
	Site         string        `json:"site"`
}

type Error struct {
	Code       string        `json:"code"`
	Msg        string        `json:"msg"`
	Args       []interface{} `json:"args"`
	Stacktrace string        `json:"stacktrace"`
}

type VcenterSites []VcenterSite

type VcenterSite struct {
	ID                   string      `json:"id"`
	Site                 string      `json:"site"`
	Description          string      `json:"description"`
	APIURL               string      `json:"apiUrl"`
	APIPublicURL         string      `json:"apiPublicUrl"`
	APIThumbprint        string      `json:"apiThumbprint"`
	IsLocal              bool        `json:"isLocal"`
	State                State       `json:"state"`
	APIVersion           string      `json:"apiVersion"`
	IsProviderDeployment bool        `json:"isProviderDeployment"`
	PeerTunnelCERT       interface{} `json:"peerTunnelCert"`
}

type CloudSites []CloudSite

type CloudSite struct {
	Site          string `json:"site"`
	Description   string `json:"description"`
	APIURL        string `json:"apiUrl"`
	APIPublicURL  string `json:"apiPublicUrl"`
	APIThumbprint string `json:"apiThumbprint"`
	IsLocal       bool   `json:"isLocal"`
	State         State  `json:"state"`
	APIVersion    string `json:"apiVersion"`
	BuildVersion  string `json:"buildVersion"`
}
