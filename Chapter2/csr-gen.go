// Run it with:
//        go run csr-gen.go client <USERNAME>;
// It will create clientkey.pem and clinet.csr 

package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		panic("Usage: go run main.go <name> <user>")
	}

	name := os.Args[1]
	user := os.Args[2]

	// Generate RSA private key
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}

	// Encode private key to PEM format
	keyDer := x509.MarshalPKCS1PrivateKey(key)
	keyBlock := pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: keyDer,
	}

	keyFile, err := os.Create(name + "-key.pem")
	if err != nil {
		panic(err)
	}
	pem.Encode(keyFile, &keyBlock)
	keyFile.Close()

	// Define certificate subject info
	commonName := user
	emailAddress := "someone@myco.com"
	org := "MyCo, Inc."
	orgUnit := "Widget Farmers"
	city := "Seattle"
	state := "WA"
	country := "US"

	subject := pkix.Name{
		CommonName:         commonName,
		Country:            []string{country},
		Locality:           []string{city},
		Organization:       []string{org},
		OrganizationalUnit: []string{orgUnit},
		Province:           []string{state},
	}

	// Convert subject to ASN.1
	asn1Subject, err := asn1.Marshal(subject.ToRDNSequence())
	if err != nil {
		panic(err)
	}

	// Create CSR
	csr := x509.CertificateRequest{
		RawSubject:         asn1Subject,
		EmailAddresses:     []string{emailAddress},
		SignatureAlgorithm: x509.SHA256WithRSA,
	}

	csrBytes, err := x509.CreateCertificateRequest(rand.Reader, &csr, key)
	if err != nil {
		panic(err)
	}

	// Write CSR to file
	csrFile, err := os.Create(name + ".csr")
	if err != nil {
		panic(err)
	}
	pem.Encode(csrFile, &pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: csrBytes,
	})
	csrFile.Close()
}
