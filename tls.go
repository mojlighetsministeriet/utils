package utils

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
)

// GetCACertificatesTLSConfig will read and return a configuration for the root certificates from /etc/ssl/certs/ca-certificates.crt that can be mounted from the host system.
func GetCACertificatesTLSConfig() (config *tls.Config, err error) {
	config, err = GetTLSConfigFromFilename("/etc/ssl/certs/ca-certificates.crt")
	return
}

// GetTLSConfigFromFilename will read and return a configuration for the certificates in a file
func GetTLSConfigFromFilename(filename string) (config *tls.Config, err error) {
	certs, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(certs)
	config = &tls.Config{RootCAs: pool}

	return
}
