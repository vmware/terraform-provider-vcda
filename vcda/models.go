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

type VspherePluginData struct {
	SsoUser     string `json:"ssoUser"`
	SsoPassword string `json:"ssoPassword"`
}

type IsServiceConfigured struct {
	IsConfigured bool `json:"isConfigured"`
}
