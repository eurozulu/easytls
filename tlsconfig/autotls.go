package tlsconfig

import (
	"crypto/tls"
	"golang.org/x/crypto/acme/autocert"
	"net/http"
	"strings"
)

type AutoTLS struct {
	certManager autocert.Manager
}

func (at AutoTLS) Config() *tls.Config {
	return &tls.Config{GetCertificate: at.certManager.GetCertificate}
}

func newAutoTLS(certDir, domain string) *AutoTLS {
	wwwDomain := strings.Join([]string{"www", domain}, ".")
	cm := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(domain, wwwDomain),
		Cache:      autocert.DirCache(certDir), //Folder for storing certificates
	}

	// listen for lets encrypt inbound
	go http.ListenAndServe(":http", cm.HTTPHandler(nil))
	return &AutoTLS{certManager: cm}

}
