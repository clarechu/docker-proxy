package testdata

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
	rootPrivateKey  = "root.key"
	rootPublicKey   = "root.pem"
	rootCertificate = "root.crt"

	serverPrivateKey  = "server.key"
	serverPublicKey   = "server.pem"
	serverCertificate = "server.crt"
)

var (
	expiry            = 365 * 24 * time.Hour
	notBefore         time.Time
	serialNumberLimit *big.Int
	serialNumber      *big.Int

	leafTemplate x509.Certificate
	rootTemplate x509.Certificate
)

func init() {
	notBefore = time.Now()
	serialNumberLimit = new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, _ = rand.Int(rand.Reader, serialNumberLimit)
	leafTemplate = x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Acme Co"},
			CommonName:   *host,
		},
		NotBefore:             notBefore,
		NotAfter:              notBefore.Add(expiry).UTC(),
		EmailAddresses:        []string{"xxx@gmail.com"},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  false,
	}
	rootTemplate = x509.Certificate{
		SerialNumber:          serialNumber,
		NotBefore:             notBefore,
		NotAfter:              notBefore.Add(expiry).UTC(),
		BasicConstraintsValid: true,
		IsCA:                  true,
		KeyUsage: x509.KeyUsageDigitalSignature |
			x509.KeyUsageKeyEncipherment | x509.KeyUsageCertSign |
			x509.KeyUsageCRLSign,
		ExtKeyUsage:    []x509.ExtKeyUsage{x509.ExtKeyUsageAny},
		EmailAddresses: []string{"xxx@gmail.com"},
		Subject: pkix.Name{
			Country:            []string{"CN"},
			Locality:           []string{"zhongguancun"},
			Province:           []string{"Beijing"},
			OrganizationalUnit: []string{"tect"},
			Organization:       []string{"paradise"},
			StreetAddress:      []string{"street"},
			PostalCode:         []string{"310000"},
			CommonName:         "localhost",
		},
	}
}

func main() {
	flag.Parse()
	if len(*host) == 0 {
		log.Fatalf("Missing required --host parameter")
	}
	GenKey1()
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

func GenCACertificate(template *x509.Certificate, filepath string) (key *ecdsa.PrivateKey) {
	rootKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	rsa.KeyToFile(*path, rootPrivateKey, rootKey)
	derBytes, err := x509.CreateCertificate(rand.Reader, template, template, &rootKey.PublicKey, rootKey)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	rsa.DebugCertToFile(filepath, rootCertificate, derBytes)
	rsa.CertToFile(filepath, rootPublicKey, derBytes)
	return rootKey
}

func GenServerCertificate(root, server *x509.Certificate, rootKey *ecdsa.PrivateKey, hostname, filepath string) {
	leafKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	rsa.KeyToFile(*path, serverPrivateKey, leafKey)

	serialNumber, err = rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		fmt.Printf("failed to generate serial number: %s", err)
		os.Exit(2)
	}

	hosts := strings.Split(hostname, ",")
	for _, h := range hosts {
		if ip := net.ParseIP(h); ip != nil {
			leafTemplate.IPAddresses = append(leafTemplate.IPAddresses, ip)
		} else {
			leafTemplate.DNSNames = append(leafTemplate.DNSNames, h)
		}
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, server, root, &leafKey.PublicKey, rootKey)
	if err != nil {
		fmt.Printf("create Certificate %s", err)
		os.Exit(2)
	}
	rsa.DebugCertToFile(filepath, serverCertificate, derBytes)
	rsa.CertToFile(filepath, serverPublicKey, derBytes)
}

func GenKey1() {
	rootKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	rsa.KeyToFile(*path, rootPrivateKey, rootKey)

	derBytes, err := x509.CreateCertificate(rand.Reader, &rootTemplate, &rootTemplate, &rootKey.PublicKey, rootKey)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	rsa.DebugCertToFile(*path, rootCertificate, derBytes)
	rsa.CertToFile(*path, rootPublicKey, derBytes)

	leafKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	rsa.KeyToFile(*path, serverPrivateKey, leafKey)

	serialNumber, err = rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		fmt.Printf("failed to generate serial number: %s", err)
		os.Exit(2)
	}

	hosts := strings.Split(*host, ",")
	for _, h := range hosts {
		if ip := net.ParseIP(h); ip != nil {
			leafTemplate.IPAddresses = append(leafTemplate.IPAddresses, ip)
		} else {
			leafTemplate.DNSNames = append(leafTemplate.DNSNames, h)
		}
	}
	// x509.ParseCertificate()
	derBytes, err = x509.CreateCertificate(rand.Reader, &leafTemplate, &rootTemplate, &leafKey.PublicKey, rootKey)
	if err != nil {
		fmt.Printf("create Certificate %s", err)
		os.Exit(2)
	}
	rsa.DebugCertToFile(*path, serverCertificate, derBytes)
	rsa.CertToFile(*path, serverPublicKey, derBytes)
}
