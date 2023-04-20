package vcda

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"

	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	VimClient     VimClient
	VcdaIP        string
	LocalUser     string
	LocalPassword string
}

func (c *Client) NewHttpClientConfig(serviceCert string) (*http.Client, error) {
	if serviceCert == "" {
		return nil, fmt.Errorf("vcda service certificate is required")
	}

	data, err := base64.StdEncoding.DecodeString(serviceCert)
	if err != nil {
		return nil, fmt.Errorf("could not decode vcda service certificate: %s", err)
	}
	block := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: data,
	}
	buf := new(bytes.Buffer)
	if err := pem.Encode(buf, block); err != nil {
		return nil, fmt.Errorf("could not encode vcda service certificate: %s", err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(buf.Bytes())

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{RootCAs: caCertPool},
	}
	client := &http.Client{Timeout: 10 * time.Second, Transport: tr}

	return client, nil
}

func (c *Client) DoRequest(host string, req *http.Request, serviceCert string) ([]byte, error) {
	authToken, err := c.GetAuthToken(host, c.LocalPassword, serviceCert)

	if err != nil {
		return nil, err
	}

	req.Header.Set(VcdaAuthTokenHeader, *authToken)
	req.Header.Set(ContentTypeHeader, ContentTypeHeaderValue)
	req.Header.Set(AcceptHeader, AcceptHeaderValue)

	hcl, err := c.NewHttpClientConfig(serviceCert)
	if err != nil {
		return nil, err
	}

	r, err := hcl.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error executing request: %s", err)
	}
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %s", err)
	}

	if !successCheck(r.StatusCode) {
		return nil, fmt.Errorf("request: %s finished with status: %d, body: %s", req.URL.String(), r.StatusCode, body)
	}

	return body, err
}

func (c *Client) doRequest(req *http.Request, serviceCert string) ([]byte, error) {
	return c.DoRequest(c.VcdaIP, req, serviceCert)
}

func successCheck(code int) bool {
	return code >= 200 && code < 300
}

func (c *Client) BuildRequestURL(host string, path string) (*string, error) {
	apiURL := url.URL{
		Scheme: "https",
		Host:   host,
		Path:   path,
	}

	u, err := url.Parse(apiURL.String())

	if err != nil {
		return nil, fmt.Errorf("could not parse API URL: %s", err)
	}

	reqURL := u.String()

	return &reqURL, nil
}

func (c *Client) buildRequestURL(path string) (*string, error) {
	return c.BuildRequestURL(c.VcdaIP, path)
}

func (c *Client) GetAuthToken(host string, password string, serviceCert string) (*string, error) {
	reqURL, err := c.BuildRequestURL(host, "/sessions")

	if err != nil {
		return nil, err
	}

	reqData := AuthTokenData{Type: UserType, LocalUser: c.LocalUser, LocalPassword: password}

	rb, err := json.Marshal(reqData)
	if err != nil {
		return nil, fmt.Errorf("could not marshal request data: %s", err)
	}

	req, err := http.NewRequest(http.MethodPost, *reqURL, strings.NewReader(string(rb)))
	req.Header.Set(ContentTypeHeader, ContentTypeHeaderValue)
	req.Header.Set(AcceptHeader, AcceptHeaderValue)

	if err != nil {
		return nil, err
	}

	hcl, err := c.NewHttpClientConfig(serviceCert)
	if err != nil {
		return nil, err
	}
	r, err := hcl.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error executing request: %s", err)
	}

	vcdaToken := r.Header.Get(VcdaAuthTokenHeader)

	return &vcdaToken, nil
}

func (c *Client) getAuthToken(password string, serviceCert string) (*string, error) {
	return c.GetAuthToken(c.VcdaIP, password, serviceCert)
}

// c4/h4 client methods
func (c *Client) changePassword(host string, currentPassword string, newPassword string, serviceCert string) error {
	reqURL, err := c.BuildRequestURL(host, "/config/root-password")

	if err != nil {
		return err
	}

	reqData := PasswordData{RootPassword: newPassword}

	rb, err := json.Marshal(reqData)
	if err != nil {
		return fmt.Errorf("could not marshal request data: %s", err)
	}

	req, err := http.NewRequest(http.MethodPost, *reqURL, strings.NewReader(string(rb)))
	if err != nil {
		return fmt.Errorf("error creating new request: %s", err)
	}

	token, err := c.GetAuthToken(host, currentPassword, serviceCert)
	if err != nil {
		return err
	}

	req.Header.Set(VcdaAuthTokenHeader, *token)
	req.Header.Set(ContentTypeHeader, ContentTypeHeaderValue)
	req.Header.Set(AcceptHeader, AcceptHeaderValue)
	req.Header.Set(ConfigSecretHeader, currentPassword)

	hcl, err := c.NewHttpClientConfig(serviceCert)
	if err != nil {
		return err
	}
	r, err := hcl.Do(req)
	if err != nil {
		return fmt.Errorf("error executing request: %s", err)
	}
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %s", err)
	}

	if r.StatusCode != http.StatusNoContent {
		return fmt.Errorf("change password failed with status: %d, body: %s", r.StatusCode, body)
	}

	c.LocalPassword = newPassword

	return nil
}

func (c *Client) checkPasswordExpired(host string, serviceCert string) (*PasswordExpiration, error) {
	reqURL, err := c.BuildRequestURL(host, "/appliance/root-password-expired")

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, *reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating new request: %s", err)
	}

	body, err := c.DoRequest(host, req, serviceCert)
	if err != nil {
		return nil, err
	}

	rootExpiration := PasswordExpiration{}
	err = json.Unmarshal(body, &rootExpiration)

	if err != nil {
		return nil, fmt.Errorf("could not unmarshal response body: %s", err)
	}

	return &rootExpiration, nil
}

func (c *Client) setLicense(serviceCert string, licenseKey string) (*License, error) {
	reqURL, err := c.buildRequestURL("/license")

	if err != nil {
		return nil, err
	}

	reqData := LicenseData{LicenseKey: licenseKey}

	rb, err := json.Marshal(reqData)
	if err != nil {
		return nil, fmt.Errorf("could not marshal request data: %s", err)
	}

	req, err := http.NewRequest(http.MethodPost, *reqURL, strings.NewReader(string(rb)))
	if err != nil {
		return nil, fmt.Errorf("error creating new request: %s", err)
	}

	body, err := c.doRequest(req, serviceCert)
	if err != nil {
		return nil, err
	}

	vcdaLicense := License{}
	err = json.Unmarshal(body, &vcdaLicense)

	if err != nil {
		return nil, fmt.Errorf("could not unmarshal response body: %s", err)
	}

	return &vcdaLicense, nil
}

func (c *Client) setSiteName(siteName string, serviceCert string) (*SiteConfig, error) {
	reqURL, err := c.buildRequestURL("/config/site")

	if err != nil {
		return nil, err
	}

	reqData := SiteData{Site: siteName}

	rb, err := json.Marshal(reqData)
	if err != nil {
		return nil, fmt.Errorf("could not marshal request data: %s", err)
	}

	req, err := http.NewRequest(http.MethodPost, *reqURL, strings.NewReader(string(rb)))
	req.Header.Set(ConfigSecretHeader, c.LocalPassword)

	if err != nil {
		return nil, fmt.Errorf("error creating new request: %s", err)
	}

	body, err := c.doRequest(req, serviceCert)
	if err != nil {
		return nil, err
	}

	vcdaSite := SiteConfig{}

	err = json.Unmarshal(body, &vcdaSite)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal response body: %s", err)
	}

	return &vcdaSite, nil
}

func (c *Client) setCloudSiteName(siteName string, description string, serviceCert string) (*CloudSiteConfig, error) {
	reqURL, err := c.buildRequestURL("/config/site")

	if err != nil {
		return nil, err
	}

	reqData := CloudSiteData{LocalSite: siteName, LocalSiteDescription: description}

	rb, err := json.Marshal(reqData)
	if err != nil {
		return nil, fmt.Errorf("could not marshal request data: %s", err)
	}

	req, err := http.NewRequest(http.MethodPost, *reqURL, strings.NewReader(string(rb)))
	req.Header.Set(ConfigSecretHeader, c.LocalPassword)

	if err != nil {
		return nil, fmt.Errorf("error creating new request: %s", err)
	}

	body, err := c.doRequest(req, serviceCert)
	if err != nil {
		return nil, err
	}

	vcdaSite := CloudSiteConfig{}

	err = json.Unmarshal(body, &vcdaSite)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal response body: %s", err)
	}

	return &vcdaSite, nil
}

func (c *Client) setPublicEndpoint(address string, port int, serviceCert string) error {
	reqURL, err := c.buildRequestURL("/config/endpoints")

	if err != nil {
		return err
	}

	reqData := EndpointData{APIAddress: nil, APIPort: 8443, APIPublicAddress: address, APIPublicPort: port,
		MgmtAddress: nil, MgmtPort: 8046, MgmtPublicAddress: nil, MgmtPublicPort: nil}

	rb, err := json.Marshal(reqData)
	if err != nil {
		return fmt.Errorf("could not marshal request data: %s", err)
	}

	req, err := http.NewRequest(http.MethodPost, *reqURL, strings.NewReader(string(rb)))

	if err != nil {
		return fmt.Errorf("error creating new request: %s", err)
	}

	body, err := c.doRequest(req, serviceCert)
	if err != nil {
		return err
	}

	endpoints := Endpoints{}
	err = json.Unmarshal(body, &endpoints)
	if err != nil {
		return fmt.Errorf("could not unmarshal response body: %s", err)
	}

	return nil
}

func (c *Client) getEndpoints(serviceCert string) (*Endpoints, error) {
	reqURL, err := c.buildRequestURL("/config/endpoints")

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, *reqURL, nil)

	if err != nil {
		return nil, fmt.Errorf("error creating new request: %s", err)
	}

	body, err := c.doRequest(req, serviceCert)
	if err != nil {
		return nil, err
	}

	endpoints := Endpoints{}

	err = json.Unmarshal(body, &endpoints)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal response body: %s", err)
	}

	return &endpoints, nil
}

func (c *Client) setLookupService(lsURL string, lsThumbprint string, serviceCert string) error {
	reqURL, err := c.buildRequestURL("config/lookup-service")

	if err != nil {
		return err
	}

	reqData := LookupServiceData{URL: lsURL, Thumbprint: lsThumbprint}

	rb, err := json.Marshal(reqData)
	if err != nil {
		return fmt.Errorf("could not marshal request data: %s", err)
	}

	req, err := http.NewRequest(http.MethodPost, *reqURL, strings.NewReader(string(rb)))
	if err != nil {
		return fmt.Errorf("error creating new request: %s", err)
	}

	body, err := c.doRequest(req, serviceCert)
	if err != nil {
		return err
	}

	lookupService := LookupService{}
	err = json.Unmarshal(body, &lookupService)

	if err != nil {
		return fmt.Errorf("could not unmarshal response body: %s", err)
	}

	return nil
}

func (c *Client) setReplicatorLookupService(host string, lsURL string, lsThumbprint string, apiURL string, apiThumbprint string, rootPassword string, serviceCert string) (*LookupService, error) {
	reqURL, err := c.BuildRequestURL(host, "/config/replicators/lookup-service")

	if err != nil {
		return nil, err
	}

	reqData := ReplicatorLookupServiceData{LsURL: lsURL, LsThumbprint: lsThumbprint, APIURL: apiURL, APIThumbprint: apiThumbprint, RootPassword: rootPassword}

	rb, err := json.Marshal(reqData)
	if err != nil {
		return nil, fmt.Errorf("could not marshal request data: %s", err)
	}

	req, err := http.NewRequest(http.MethodPost, *reqURL, strings.NewReader(string(rb)))
	if err != nil {
		return nil, fmt.Errorf("error creating new request: %s", err)
	}

	body, err := c.DoRequest(host, req, serviceCert)
	if err != nil {
		return nil, err
	}

	lookupService := LookupService{}
	err = json.Unmarshal(body, &lookupService)

	if err != nil {
		return nil, fmt.Errorf("could not unmarshal response body: %s", err)
	}

	return &lookupService, nil
}

func (c *Client) setVcloud(vcdUsername string, vcdPassword string, vcdURL string, vcdThumbprint string, serviceCert string) error {
	reqData := VcloudConfigData{VcdPassword: vcdPassword, VcdThumbprint: vcdThumbprint, VcdURL: vcdURL + "/api", VcdUsername: vcdUsername}

	rb, err := json.Marshal(reqData)
	if err != nil {
		return fmt.Errorf("could not marshal request data: %s", err)
	}

	reqURL, err := c.buildRequestURL("/config/vcloud")

	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, *reqURL, strings.NewReader(string(rb)))
	if err != nil {
		return fmt.Errorf("error creating new request: %s", err)
	}

	body, err := c.doRequest(req, serviceCert)
	if err != nil {
		return err
	}

	if err != nil {
		return fmt.Errorf("error reading response body: %s", body)
	}

	vcloudConfig := CloudSiteConfig{}
	err = json.Unmarshal(body, &vcloudConfig)

	if err != nil {
		return fmt.Errorf("could not unmarshal response body: %s", err)
	}

	return nil
}

func (c *Client) setTunnel(tunnelURL string, tunnelCertificate string, tunnelRootPassword string, serviceCert string) (*CloudSiteConfig, error) {
	reqURL, err := c.buildRequestURL("/config/tunnel-service")

	if err != nil {
		return nil, err
	}

	reqData := TunnelData{Certificate: tunnelCertificate, RootPassword: tunnelRootPassword, URL: tunnelURL}

	rb, err := json.Marshal(reqData)
	if err != nil {
		return nil, fmt.Errorf("could not marshal request data: %s", err)
	}

	req, err := http.NewRequest(http.MethodPost, *reqURL, strings.NewReader(string(rb)))
	if err != nil {
		return nil, fmt.Errorf("error creating new request: %s", err)
	}

	body, err := c.doRequest(req, serviceCert)
	if err != nil {
		return nil, err
	}

	tunnelConfig := CloudSiteConfig{}
	err = json.Unmarshal(body, &tunnelConfig)

	if err != nil {
		return nil, fmt.Errorf("could not unmarshal response body: %s", err)
	}

	return &tunnelConfig, nil
}

func (c *Client) getManagerSiteConfig(serviceCert string) (*SiteConfig, error) {
	reqURL, err := c.buildRequestURL("/config")

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, *reqURL, nil)

	if err != nil {
		return nil, fmt.Errorf("error creating new request: %s", err)
	}

	body, err := c.doRequest(req, serviceCert)
	if err != nil {
		return nil, err
	}

	vcdaSite := SiteConfig{}

	err = json.Unmarshal(body, &vcdaSite)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal response body: %s", err)
	}

	return &vcdaSite, nil
}

func (c *Client) getCloudSiteConfig(serviceCert string) (*CloudSiteConfig, error) {
	reqURL, err := c.buildRequestURL("/config")

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, *reqURL, nil)

	if err != nil {
		return nil, fmt.Errorf("error creating new request: %s", err)
	}

	body, err := c.doRequest(req, serviceCert)
	if err != nil {
		return nil, err
	}

	vcdaSite := CloudSiteConfig{}

	err = json.Unmarshal(body, &vcdaSite)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal response body: %s", err)
	}

	return &vcdaSite, nil
}

func (c *Client) addReplicator(host string, serviceCert string, description string, owner string, siteName string, details ReplicatorConfigData) (*Replicator, error) {
	reqData := ReplicatorData{Description: description, Owner: owner, Site: siteName, ReplicatorID: nil, Details: details}

	rb, err := json.Marshal(reqData)
	if err != nil {
		return nil, fmt.Errorf("could not marshal request data: %s", err)
	}

	reqURL, err := c.BuildRequestURL(host, "/replicators")
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, *reqURL, strings.NewReader(string(rb)))
	if err != nil {
		return nil, fmt.Errorf("error creating new request: %s", err)
	}

	body, err := c.DoRequest(host, req, serviceCert)
	if err != nil {
		return nil, err
	}

	replicator := Replicator{}
	err = json.Unmarshal(body, &replicator)

	if err != nil {
		return nil, fmt.Errorf("could not unmarshal response body: %s", err)
	}

	return &replicator, nil
}

func (c *Client) getReplicator(host string, serviceCert string, replicatorID string) (*Replicator, error) {
	reqURL, err := c.BuildRequestURL(host, "/replicators")

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, *reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating new request: %s", err)
	}

	body, err := c.DoRequest(host, req, serviceCert)
	if err != nil {
		return nil, err
	}

	var replicators []Replicator
	err = json.Unmarshal(body, &replicators)

	if err != nil {
		return nil, fmt.Errorf("could not unmarshal response body: %s", err)
	}

	var replicator *Replicator
	for _, r := range replicators {
		if r.ID == replicatorID {
			replicator = &r
			break
		}
	}
	if replicator == nil {
		return nil, fmt.Errorf("replicator with ID: %s was not found", replicatorID)
	}

	return replicator, nil
}

func (c *Client) repairReplicator(host string, serviceCert string, replicatorID string, apiURL string, apiThumbprint string, rootPassword string, ssoUser string, ssoPassword string) error {
	reqData := ReplicatorConfigData{APIURL: apiURL, APIThumbprint: apiThumbprint, RootPassword: rootPassword, SsoUser: ssoUser, SsoPassword: ssoPassword}

	rb, err := json.Marshal(reqData)
	if err != nil {
		return fmt.Errorf("could not marshal request data: %s", err)
	}

	path := "/replicators" + replicatorID + "/reset-cookie"
	reqURL, err := c.BuildRequestURL(host, path)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, *reqURL, strings.NewReader(string(rb)))
	if err != nil {
		return fmt.Errorf("error creating new request: %s", err)
	}

	body, err := c.DoRequest(host, req, serviceCert)
	if err != nil {
		return err
	}

	replicator := Replicator{}
	err = json.Unmarshal(body, &replicator)

	if err != nil {
		return fmt.Errorf("could not unmarshal response body: %s", err)
	}

	return nil
}

func (c *Client) deleteReplicator(host string, serviceCert string, replicatorID string) error {
	reqURL, err := c.BuildRequestURL(host, "/replicators/"+replicatorID)

	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodDelete, *reqURL, nil)
	if err != nil {
		return fmt.Errorf("error creating new request: %s", err)
	}

	body, err := c.DoRequest(host, req, serviceCert)
	if err != nil {
		return err
	}

	if len(body) > 0 {
		resBody := make(map[string]interface{})
		err = json.Unmarshal(body, &resBody)

		if err != nil {
			return fmt.Errorf("could not unmarshal response body: %s", err)
		}
		return fmt.Errorf("error deleting replicator: %s with response body: %s", replicatorID, resBody)
	}

	return nil
}

func (c *Client) setVspherePlugin(ssoUser string, ssoPassword string, serviceCert string) error {
	reqURL, err := c.buildRequestURL("config/vsphere-ui/register")

	if err != nil {
		return err
	}

	reqData := VspherePluginData{SsoUser: ssoUser, SsoPassword: ssoPassword}

	rb, err := json.Marshal(reqData)
	if err != nil {
		return fmt.Errorf("could not marshal request data: %s", err)
	}

	req, err := http.NewRequest(http.MethodPost, *reqURL, strings.NewReader(string(rb)))
	if err != nil {
		return fmt.Errorf("error creating new request: %s", err)
	}

	_, err = c.doRequest(req, serviceCert)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) removeVspherePlugin(ssoUser string, ssoPassword string, serviceCert string) error {
	reqURL, err := c.buildRequestURL("config/vsphere-ui/unregister")

	if err != nil {
		return err
	}

	reqData := VspherePluginData{SsoUser: ssoUser, SsoPassword: ssoPassword}

	rb, err := json.Marshal(reqData)
	if err != nil {
		return fmt.Errorf("could not marshal request data: %s", err)
	}

	req, err := http.NewRequest(http.MethodPost, *reqURL, strings.NewReader(string(rb)))
	if err != nil {
		return fmt.Errorf("error creating new request: %s", err)
	}

	_, err = c.doRequest(req, serviceCert)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) isConfigured(serviceCert string) (*IsServiceConfigured, error) {
	reqURL, err := c.buildRequestURL("/config/is-configured")

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, *reqURL, nil)

	if err != nil {
		return nil, fmt.Errorf("error creating new request: %s", err)
	}

	body, err := c.doRequest(req, serviceCert)
	if err != nil {
		return nil, err
	}

	isServiceConfigured := IsServiceConfigured{}

	err = json.Unmarshal(body, &isServiceConfigured)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal response body: %s", err)
	}

	return &isServiceConfigured, nil
}
