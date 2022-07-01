package cmd

import (
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"github.com/clarechu/docker-proxy/pkg/ssl"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	log "k8s.io/klog/v2"
	"math/big"
	"os"
	"time"
)

var (
	expiry       = 365 * 24 * time.Hour
	subject      *pkix.Name
	emailAddress = ""
	filepath     = ""
	hostname     = ""
)

func GenerateCommand() *cobra.Command {
	subject = new(pkix.Name)
	genCmd := &cobra.Command{
		Use:   "gen",
		Short: "generate https certificate",
		Long: `
Tips  Find more information at: https://github.com/clarechu/docker-proxy
Example:
	docker-proxy gen 

Certificate Host Name (localhost) []: localhost
Country Name(2 letter code) []AU: CN
State or Province Name (full name) [Some-State]: Guangdong
Locality Name (eg, city): ShenZheng
Organization Name (eg, company) [Internet Widgits Pty Ltd]: demo
Common Name (e.g. server FQDN or YOUR name) []: hh
Email Address []: 1062186165@qq.com

`,
		PreRun: func(cmd *cobra.Command, args []string) {
			if hostname == "" {
				hostname, _ = setValueInteractive("Certificate Host Name (localhost) []", "localhost")
			}
			if len(subject.Country) == 0 {
				country, _ := setValueInteractive("Country Name(2 letter code) []AU", "CN")
				subject.Country = []string{country}
			}
			if len(subject.Province) == 0 {
				province, _ := setValueInteractive("State or Province Name (full name) [Some-State]", "Guangdong")
				subject.Province = []string{province}
			}
			if len(subject.Locality) == 0 {
				locality, _ := setValueInteractive("Locality Name (eg, city)", "ShenZheng")
				subject.Locality = []string{locality}
			}
			if len(subject.Organization) == 0 {
				organization, _ := setValueInteractive("Organization Name (eg, company) [Internet Widgits Pty Ltd]", "demo")
				subject.Organization = []string{organization}
			}
			if len(subject.OrganizationalUnit) == 0 {
				organizationalUnit, _ := setValueInteractive("Organization Unit Name (eg, section) []", "demo1")
				subject.Province = []string{organizationalUnit}
			}
			if subject.CommonName == "" {
				commonName, _ := setValueInteractive("Common Name (e.g. server FQDN or YOUR name) []", "")
				subject.Province = []string{commonName}
			}
			if emailAddress == "" {
				emailAddress, _ = setValueInteractive("Email Address []", "")
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			_, err := os.Stat(filepath)
			if os.IsNotExist(err) {
				_ = os.Mkdir(filepath, os.ModeDir)
			}
			notBefore := time.Now()
			notAfter := notBefore.Add(expiry)
			serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
			serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
			if err != nil {
				log.Fatalf("failed to generate serial number: %s", err)
			}
			rootTemplate := &x509.Certificate{
				SerialNumber:          serialNumber,
				NotBefore:             notBefore,
				NotAfter:              notAfter.UTC(),
				BasicConstraintsValid: true,
				IsCA:                  true,
				KeyUsage: x509.KeyUsageDigitalSignature |
					x509.KeyUsageKeyEncipherment | x509.KeyUsageCertSign |
					x509.KeyUsageCRLSign,
				ExtKeyUsage:    []x509.ExtKeyUsage{x509.ExtKeyUsageAny},
				EmailAddresses: []string{emailAddress},
				Subject:        *subject,
			}
			// generate ca root Certificate
			key := ssl.GenCACertificate(rootTemplate, filepath)
			// generate server Certificate
			ssl.GenServerCertificate(rootTemplate, rootTemplate, key, hostname, filepath)
		},
	}
	addGenFlag(genCmd)
	return genCmd
}

func addGenFlag(genCmd *cobra.Command) {
	genCmd.PersistentFlags().StringArrayVar(&subject.Country, "country", []string{}, "Country Name(2 letter code) []AU:")
	genCmd.PersistentFlags().StringArrayVar(&subject.Province, "province", []string{}, "State or Province Name (full name) [Some-State]:")
	genCmd.PersistentFlags().StringArrayVar(&subject.Locality, "locality", []string{}, "Locality Name (eg, city):")
	genCmd.PersistentFlags().StringArrayVar(&subject.Organization, "organization", []string{}, "Organization Name (eg, company) [Internet Widgits Pty Ltd]:")
	genCmd.PersistentFlags().StringArrayVar(&subject.OrganizationalUnit, "organizational-unit", []string{""}, "Organization Unit Name (eg, section) []:")
	genCmd.PersistentFlags().StringVar(&subject.CommonName, "common-name", "", " Common Name (e.g. server FQDN or YOUR name) []:")
	genCmd.PersistentFlags().StringVar(&emailAddress, "email", "", "Email Address []:")
	genCmd.PersistentFlags().StringVar(&filepath, "filepath", "./", "CA File Path []:")

}

func setValueInteractive(message, value string) (string, error) {
	prompt := promptui.Prompt{
		Label:   message,
		Default: value,
	}
	result, err := prompt.Run()
	if err != nil {
		panic(err)
	}
	return result, nil
}

func selectValueInteractive(message string, options interface{}) (string, error) {
	prompt := promptui.Select{
		Label: message,
		Items: options,
	}
	_, result, err := prompt.Run()
	if err != nil {
		return "", err
	}
	return result, nil
}
