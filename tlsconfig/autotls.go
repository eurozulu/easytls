package tlsconfig

import (
	"crypto/tls"
	"golang.org/x/crypto/acme/autocert"
	"log"
	"net/http"
	"strings"
)

type AutoTLS struct {
	certManager autocert.Manager
}

func (at AutoTLS) Config() *tls.Config {
	return at.certManager.TLSConfig()
}

func newAutoTLS(certDir, domain string) *AutoTLS {
	wwwDomain := strings.Join([]string{"www", domain}, ".")
	cm := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(domain, wwwDomain),
		Cache:      autocert.DirCache(certDir), //Folder for storing certificates
	}

	// listen for lets encrypt inbound
	go func() {
		log.Println("starting http listener")
		defer log.Println("closing http listener")
		proxyHandler := &testhandler{proxy: cm.HTTPHandler(nil)}
		if err := http.ListenAndServe(":http", proxyHandler); err != nil {
			log.Printf("Error in http listener  %w", err)
		}
	}()
	return &AutoTLS{certManager: cm}

}

type testhandler struct {
	proxy http.Handler
}

func (t testhandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	log.Printf("Request Started: %v\n", request)
	t.proxy.ServeHTTP(writer, request)
	log.Printf("Request Done: %v\n", request)

}
