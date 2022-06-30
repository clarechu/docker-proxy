package cmd

import (
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"github.com/spf13/cobra"
	log "k8s.io/klog/v2"
	"math/big"
	"time"
)

var subject *pkix.Name
var expiry = 365 * 24 * time.Hour
var emailAddress string

func GenerateCommand() *cobra.Command {
	subject = new(pkix.Name)
	emailAddress = ""
	genCmd := &cobra.Command{
		Use:   "gen",
		Short: "generate https certificate",
		Long: `
Tips  Find more information at: https://github.com/clarechu/docker-proxy
Example:
	docker-proxy gen 
`,
		PreRun: func(cmd *cobra.Command, args []string) {

		},
		Run: func(cmd *cobra.Command, args []string) {
			notBefore := time.Now()
			notAfter := notBefore.Add(expiry)
			serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
			serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
			if err != nil {
				log.Fatalf("failed to generate serial number: %s", err)
			}
			rootTemplate := x509.Certificate{
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
			log.Infof("root Template", rootTemplate)
		},
	}
	addGenFlag(genCmd)
	return genCmd
}

func addGenFlag(genCmd *cobra.Command) {
	genCmd.PersistentFlags().StringArrayVar(&subject.Country, "country", []string{"CN"}, "Country Name(2 letter code) []AU:")
	genCmd.PersistentFlags().StringArrayVar(&subject.Province, "province", []string{"Guangdong"}, "State or Province Name (full name) [Some-State]:")
	genCmd.PersistentFlags().StringArrayVar(&subject.Locality, "locality", []string{"ShenZheng"}, "Locality Name (eg, city):")
	genCmd.PersistentFlags().StringArrayVar(&subject.Organization, "organization", []string{"demo"}, "Organization Name (eg, company) [Internet Widgits Pty Ltd]:")
	genCmd.PersistentFlags().StringArrayVar(&subject.OrganizationalUnit, "organizational-unit", []string{"demo"}, "Organization Unit Name (eg, section) []:")
	genCmd.PersistentFlags().StringVar(&subject.CommonName, "common-name", "clarechu", " Common Name (e.g. server FQDN or YOUR name) []:")
	genCmd.PersistentFlags().StringVar(&emailAddress, "email", "", "Email Address []")

}
