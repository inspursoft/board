package helm

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

//HttpClient is the efault HTTP(/S) backend handler
type HttpClient struct { //nolint
	client   *http.Client
	username string
	password string
}

//Get performs a Http Get and returns the body.
func (g *HttpClient) Get(href string) (*bytes.Buffer, error) {
	buf := bytes.NewBuffer(nil)

	// Set a helm specific user agent so that a repo server and metrics can
	// separate helm calls from other tools interacting with repos.
	req, err := http.NewRequest("GET", href, nil)
	if err != nil {
		return buf, err
	}

	if g.username != "" && g.password != "" {
		req.SetBasicAuth(g.username, g.password)
	}

	resp, err := g.client.Do(req)
	if err != nil {
		return buf, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return buf, fmt.Errorf("Failed to fetch %s : %s", href, resp.Status)
	}

	_, err = io.Copy(buf, resp.Body)
	return buf, err
}

//Upload performs a Http Post.
func (g *HttpClient) Upload(href string, body io.Reader) (*bytes.Buffer, error) {
	// Set a helm specific user agent so that a repo server and metrics can
	// separate helm calls from other tools interacting with repos.
	req, err := http.NewRequest("POST", href, body)
	if err != nil {
		return nil, err
	}

	if g.username != "" && g.password != "" {
		req.SetBasicAuth(g.username, g.password)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := g.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 201 {
		return nil, fmt.Errorf("Failed to fetch %s : %s", href, resp.Status)
	}

	buf := bytes.NewBuffer(nil)
	_, err = io.Copy(buf, resp.Body)

	return buf, err
}

//Delete performs a Http Delete.
func (g *HttpClient) Delete(href string) (*bytes.Buffer, error) {
	// Set a helm specific user agent so that a repo server and metrics can
	// separate helm calls from other tools interacting with repos.
	req, err := http.NewRequest("DELETE", href, nil)
	if err != nil {
		return nil, err
	}

	if g.username != "" && g.password != "" {
		req.SetBasicAuth(g.username, g.password)
	}

	resp, err := g.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Failed to fetch %s : %s", href, resp.Status)
	}

	buf := bytes.NewBuffer(nil)
	_, err = io.Copy(buf, resp.Body)

	return buf, nil
}

// NewHTTPClient constructs a valid http/https client as HttpClient
func NewHTTPClient(URL, username, password string, Cert, Key, CA []byte) (*HttpClient, error) {
	tr := &http.Transport{
		DisableCompression: true,
		Proxy:              http.ProxyFromEnvironment,
	}
	if (Cert != nil && Key != nil) || CA != nil {
		tlsConf, err := NewTLSConfig(URL, Cert, Key, CA)
		if err != nil {
			return nil, fmt.Errorf("can't create TLS config: %s", err.Error())
		}
		tr.TLSClientConfig = tlsConf
	}
	client := HttpClient{
		client:   &http.Client{Transport: tr},
		username: username,
		password: password,
	}
	return &client, nil
}

// NewTLSConfig returns tls.Config appropriate for client and/or server auth.
func NewTLSConfig(url string, cert, key, ca []byte) (*tls.Config, error) {
	config, err := newTLSConfigCommon(cert, key, ca)
	if err != nil {
		return nil, err
	}
	config.BuildNameToCertificate()

	serverName, err := ExtractHostname(url)
	if err != nil {
		return nil, err
	}
	config.ServerName = serverName

	return config, nil
}

func newTLSConfigCommon(cert, key, ca []byte) (*tls.Config, error) {
	config := tls.Config{}

	if cert != nil && key != nil {
		certPair, err := CertFromPair(cert, key)
		if err != nil {
			return nil, err
		}
		config.Certificates = []tls.Certificate{*certPair}
	}

	if ca != nil {
		cp, err := CertPool(ca)
		if err != nil {
			return nil, err
		}
		config.RootCAs = cp
	}

	return &config, nil
}

// CertPool returns an x509.CertPool containing the certificates
// in the given PEM-encoded bytes.
// Returns an error if a certificate could not
// be parsed, or if the bytes does not contain any certificates
func CertPool(ca []byte) (*x509.CertPool, error) {
	cp := x509.NewCertPool()
	if !cp.AppendCertsFromPEM(ca) {
		return nil, fmt.Errorf("failed to append certificates %v", ca)
	}
	return cp, nil
}

// CertFromPair returns an tls.Certificate containing the
// certificates public/private key pair from a pair of given PEM-encoded bytes.
// Returns an error if a certificate could not
// be parsed, or if the bytes does not contain any certificates
func CertFromPair(cert, key []byte) (*tls.Certificate, error) {
	certPair, err := tls.X509KeyPair(cert, key)
	if err != nil {
		return nil, fmt.Errorf("can't load key pair from cert %s and key %v: %v", cert, key, err)
	}
	return &certPair, err
}

// ExtractHostname returns hostname from URL
func ExtractHostname(addr string) (string, error) {
	u, err := url.Parse(addr)
	if err != nil {
		return "", err
	}
	return stripPort(u.Host), nil
}

func stripPort(hostport string) string {
	colon := strings.IndexByte(hostport, ':')
	if colon == -1 {
		return hostport
	}
	if i := strings.IndexByte(hostport, ']'); i != -1 {
		return strings.TrimPrefix(hostport[:i], "[")
	}
	return hostport[:colon]

}
