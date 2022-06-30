package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"flag"
	"fmt"
	"github.com/clarechu/docker-proxy/pkg/utils/rsa"
	log "k8s.io/klog/v2"
	"math/big"
	"net"
	"os"
	"strings"
	"time"
)

var (
	host = flag.String("host", "", "Comma-separated hostnames and IPs to generate a certificate for")
	path = flag.String("path", "./", "path to generate a certificate for")
)

var (
	expiry = 365 * 24 * time.Hour
)

func main() {
	flag.Parse()
	if len(*host) == 0 {
		log.Fatalf("Missing required --host parameter")
	}
	GenKey()
	fmt.Fprintf(os.Stdout, `Successfully generated certificates! Here's what you generated.

# Root CA

root.key
	The private key for the root Certificate Authority. Keep this private.

root.pem
	The public key for the root Certificate Authority. Clients should load the
	certificate in this file to connect to the server.

root.debug.crt
	information about the generated certificate.

# Leaf Certificate - Use these to serve TLS traffic.

server.key
	Private key (PEM-encoded) for terminating TLS traffic on the server.

server.pem
	Public key for terminating TLS traffic on the server.

server.crt
	information about the generated certificate

`)
}

func GenKey() {
	notBefore := time.Now()
	notAfter := notBefore.Add(expiry)
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		log.Fatalf("failed to generate serial number: %s", err)
	}
	rootKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	rsa.KeyToFile(*path, "root.key", rootKey)
	rootTemplate := x509.Certificate{
		SerialNumber:          serialNumber,
		NotBefore:             notBefore,
		NotAfter:              notBefore.Add(expiry).UTC(),
		BasicConstraintsValid: true,
		IsCA:                  true,
		KeyUsage: x509.KeyUsageDigitalSignature |
			x509.KeyUsageKeyEncipherment | x509.KeyUsageCertSign |
			x509.KeyUsageCRLSign,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageAny},
		Subject: pkix.Name{
			Country:            []string{"CN"},
			Locality:           []string{"zhongguancun"},
			Province:           []string{"Beijing"},
			OrganizationalUnit: []string{"tect"},
			Organization:       []string{"paradise"},
			StreetAddress:      []string{"street", "address", "demo"},
			PostalCode:         []string{"310000"},
			CommonName:         "localhost",
		},
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &rootTemplate, &rootTemplate, &rootKey.PublicKey, rootKey)
	if err != nil {
		panic(err)
	}
	rsa.DebugCertToFile(*path, "root.crt", derBytes)
	rsa.CertToFile(*path, "root.pem", derBytes)

	leafKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	rsa.KeyToFile(*path, "server.key", leafKey)

	serialNumber, err = rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		log.Fatalf("failed to generate serial number: %s", err)
	}
	leafTemplate := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Acme Co"},
			CommonName:   *host,
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  false,
	}
	hosts := strings.Split(*host, ",")
	for _, h := range hosts {
		if ip := net.ParseIP(h); ip != nil {
			leafTemplate.IPAddresses = append(leafTemplate.IPAddresses, ip)
		} else {
			leafTemplate.DNSNames = append(leafTemplate.DNSNames, h)
		}
	}

	derBytes, err = x509.CreateCertificate(rand.Reader, &leafTemplate, &rootTemplate, &leafKey.PublicKey, rootKey)
	if err != nil {
		panic(err)
	}
	rsa.DebugCertToFile(*path, "server.crt", derBytes)
	rsa.CertToFile(*path, "server.pem", derBytes)
}
