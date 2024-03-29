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
const ENV_VERBOSE = "EASYTLS_VERBOSE"

// used : https://stackoverflow.com/questions/37321760/how-to-set-up-lets-encrypt-for-a-go-server-application for lets encrypt
func main() {
	verbose := os.Getenv(ENV_VERBOSE) != ""

	certPath := path.Join(os.ExpandEnv("$PWD"), "certs")
	domain := os.Getenv(ENV_DOMAINNAME)
	cfg, err := tlsconfig.NewTLSConfig(certPath, domain)
	if err != nil {
		log.Fatalln(err)
	}

	h := NewHandler()
	h.Verbose = verbose

	mux := http.NewServeMux()
	mux.HandleFunc("/", h.handleRequest)

	srv := &http.Server{
		Addr:         ":https",
		Handler:      mux,
		TLSConfig:    cfg.Config(),
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}
	log.Printf("Starting to listen for domain %s\n", domain)
	err = srv.ListenAndServeTLS("", "")
	if err != http.ErrServerClosed {
		log.Fatalln(err)
	}
}
