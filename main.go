package main

import (
	"crypto/tls"
	"easytls/tlsconfig"
	"log"
	"net/http"
	"os"
	"path"
)

const ENV_DOMAINNAME = "EASYTLS_DOMAIN"

// used : https://stackoverflow.com/questions/37321760/how-to-set-up-lets-encrypt-for-a-go-server-application for lets encrypt
func main() {
	certPath := path.Join(os.ExpandEnv("$PWD"), "certs")
	domain := os.Getenv(ENV_DOMAINNAME)
	cfg, err := tlsconfig.NewTLSConfig(certPath, domain)
	if err != nil {
		log.Fatalln(err)
	}

	h := NewHandler()
	mux := http.NewServeMux()
	mux.HandleFunc("/", h.handleRequest)

	srv := &http.Server{
		Addr:         ":https",
		Handler:      mux,
		TLSConfig:    cfg.Config(),
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}

	err = srv.ListenAndServeTLS("", "")
	if err != http.ErrServerClosed {
		log.Fatalln(err)
	}
}
