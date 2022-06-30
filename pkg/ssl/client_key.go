package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"github.com/clarechu/docker-proxy/pkg/utils/rsa"
	"math/big"
	"os"
	"time"
)

func GenClientKey(path string, rootTemplate x509.Certificate, rootKey *ecdsa.PrivateKey) {
	notBefore := time.Now()
	notAfter := notBefore.Add(expiry)
	clientKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	rsa.KeyToFile(path, "client.key", clientKey)

	clientTemplate := x509.Certificate{
		SerialNumber: new(big.Int).SetInt64(4),
		Subject: pkix.Name{
			Organization: []string{"Acme Co"},
			CommonName:   "client_auth_test_cert",
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
		IsCA:                  false,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &clientTemplate, &rootTemplate, &clientKey.PublicKey, rootKey)
	if err != nil {
		panic(err)
	}
	rsa.DebugCertToFile(path, "client.debug.crt", derBytes)
	rsa.CertToFile(path, "client.pem", derBytes)
	fmt.Fprintf(os.Stdout, `Successfully generated certificates! Here's what you generated.
# Client Certificate - You probably don't need these.

client.key: Secret key for TLS client authentication
client.pem: Public key for TLS client authentication
`)
}
