package tradier

import (
	"crypto/x509"
	"fmt"
)

func showCertDNSNames(rawCerts [][]byte, _ [][]*x509.Certificate) error {
	fmt.Println("DNS Names from server's x509 certificate:")

	x509Cert, err := x509.ParseCertificate(rawCerts[0])
	if err != nil {
		fmt.Printf("Error parsing x509 cert: %s", err.Error())
		return nil
	}

	for i, dnsName := range x509Cert.DNSNames {
		fmt.Printf("  %2d) %s\n", i+1, dnsName)
	}
	fmt.Println()

	return nil
}
