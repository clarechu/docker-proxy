package ssl

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"fmt"
	"github.com/clarechu/docker-proxy/pkg/utils/rsa"
	"net"
	"os"
	"strings"
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
	leafTemplate x509.Certificate
)

func GenCACertificate(template *x509.Certificate, filepath string) (key *ecdsa.PrivateKey) {
	rootKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	rsa.KeyToFile(filepath, rootPrivateKey, rootKey)
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
	rsa.KeyToFile(filepath, serverPrivateKey, leafKey)

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
