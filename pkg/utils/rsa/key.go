package rsa

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	log "k8s.io/klog/v2"
	"os"
	"os/exec"
	"path/filepath"
)

// KeyToFile writes a PEM serialization of |key| to a new file called
// |filename|.
func KeyToFile(path, filename string, key *ecdsa.PrivateKey) {
	file, err := os.Create(filepath.Join(path, filename))
	if err != nil {
		log.Fatalf("failed to write data to cert.key: %s", err)
	}
	defer file.Close()
	b, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		log.Fatalf("Unable to marshal ECDSA private key: %v", err)
	}
	if err := pem.Encode(file, &pem.Block{Type: "EC PRIVATE KEY", Bytes: b}); err != nil {
		log.Fatalf("pem Encode :%v", err)
	}
}

func CertToFile(path, filename string, derBytes []byte) {
	certOut, err := os.Create(filepath.Join(path, filename))
	if err != nil {
		log.Fatalf("failed to open cert.pem for writing: %s", err)
	}
	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		log.Fatalf("failed to write data to cert.pem: %s", err)
	}
	if err := certOut.Close(); err != nil {
		log.Fatalf("error closing cert.pem: %s", err)
	}
}

// DebugCertToFile writes a PEM serialization and OpenSSL debugging dump of
// |derBytes| to a new file called |filename|.
func DebugCertToFile(path, filename string, derBytes []byte) {
	cmd := exec.Command("openssl", "x509", "-text", "-inform", "DER")

	file, err := os.Create(filepath.Join(path, filename))
	if err != nil {
		log.Fatalf("failed to open cert.pem for writing: %s", err)
	}
	defer file.Close()
	cmd.Stdout = file
	cmd.Stderr = os.Stderr

	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatalf("failed to write data to cert.pem: %s", err)
	}

	if err := cmd.Start(); err != nil {
		log.Fatalf("failed to write data to cert.pem: %s", err)
	}
	if _, err := stdin.Write(derBytes); err != nil {
		log.Fatalf("failed to write data to cert.pem: %s", err)
	}
	stdin.Close()
	if err := cmd.Wait(); err != nil {
		log.Fatalf("failed to write data to cert.pem: %s", err)
	}
}
