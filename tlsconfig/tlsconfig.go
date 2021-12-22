package tlsconfig

import (
	"crypto/tls"
	"log"
	"os"
)

type TLSConfig interface {
	Config() *tls.Config
}

func NewTLSConfig(certPath, domain string) (TLSConfig, error) {
	if domain != "" {
		return newAutoTLS(certPath, domain), nil
	} else {
		hn, err := os.Hostname()
		if err != nil {
			log.Fatalln(err)
		}
		return newSimpleTLS(certPath, hn)
	}
}
