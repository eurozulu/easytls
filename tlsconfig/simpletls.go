package tlsconfig

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io/ioutil"
	"math/big"
	"os"
	"path"
	"time"
)

const keyName = "private.key"
const certName = "cert.crt"
const keyLength = 2048

type SimpleTls struct {
	certs []tls.Certificate
}

func (st SimpleTls) Config() *tls.Config {
	return &tls.Config{
		Certificates:             st.certs,
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
	}
}

func newSimpleTLS(certPath, host string) (*SimpleTls, error) {
	p := path.Join(certPath, host)
	if err := ensureHost(p); err != nil {
		return nil, err
	}
	kp := path.Join(p, keyName)
	cp := path.Join(p, certName)
	if err := ensureKey(kp); err != nil {
		return nil, err
	}
	if err := ensureCertificate(cp); err != nil {
		return nil, err
	}
	c, err := tls.LoadX509KeyPair(cp, kp)
	if err != nil {
		return nil, err
	}
	return &SimpleTls{certs: []tls.Certificate{c}}, nil
}

func ensureHost(p string) error {
	return os.MkdirAll(p, 0750)
}

func ensureKey(p string) error {
	_, err := os.Stat(p)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		return generateKey(p)
	}
	_, err = readKey(p)
	return err
}

func readKey(p string) (crypto.PrivateKey, error) {
	by, err := ioutil.ReadFile(p)
	if err != nil {
		return nil, err
	}
	pb, _ := pem.Decode(by)
	return x509.ParsePKCS1PrivateKey(pb.Bytes)
}

func generateKey(p string) error {
	k, err := rsa.GenerateKey(rand.Reader, keyLength)
	if err != nil {
		return err
	}
	der := x509.MarshalPKCS1PrivateKey(k)
	by := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: der,
	})
	return ioutil.WriteFile(p, by, 0600)
}

func ensureCertificate(p string) error {
	_, err := os.Stat(p)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		return generateCertificate(p)
	}
	_, err = readCertificate(p)
	return err
}

func readCertificate(p string) (*x509.Certificate, error) {
	by, err := ioutil.ReadFile(p)
	if err != nil {
		return nil, err
	}
	pb, _ := pem.Decode(by)
	return x509.ParseCertificate(pb.Bytes)
}

func generateCertificate(p string) error {
	k, err := readKey(path.Join(path.Dir(p), keyName))
	if err != nil {
		return err
	}
	prk := k.(*rsa.PrivateKey)
	puk := &prk.PublicKey

	n := pkix.Name{
		CommonName: path.Base(path.Dir(p)),
	}
	c := &x509.Certificate{
		SerialNumber:                big.NewInt(1),
		SignatureAlgorithm:          x509.SHA512WithRSA,
		PublicKey:                   puk,
		Version:                     1,
		Issuer:                      n,
		Subject:                     n,
		NotBefore:                   time.Now(),
		NotAfter:                    time.Now().AddDate(1, 0, 0),
		KeyUsage:                    x509.KeyUsageKeyAgreement | x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		UnhandledCriticalExtensions: nil,
		ExtKeyUsage:                 []x509.ExtKeyUsage{x509.ExtKeyUsageEmailProtection, x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageOCSPSigning},
		IsCA:                        false,
	}
	by, err := x509.CreateCertificate(rand.Reader, c, c, puk, prk)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(p, pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: by,
	}), 0644)
}
