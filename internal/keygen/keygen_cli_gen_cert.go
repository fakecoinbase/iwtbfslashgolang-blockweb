package keygen

/*
 * Copyright 2020 Information Wants To Be Free
 * Visit: https://github.com/iwtbf
 *
 * This project is licensed under the terms of the Apache 2.0 License.
 */

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"time"
)

type genCertCmd struct {
	Out      string `arg help:"Path to the output cert file."`
	Key      string `flag type:"existingfile" help:"Path to the input key file."`
	ValidFor string `flag help:"Validity of the generated certificate." default:"8760h"`
}

func parseDuration(validFor string) time.Duration {
	duration, err := time.ParseDuration(validFor)
	if err != nil {
		panic(err)
	}

	return duration
}

func newTemplate(validFor time.Duration) x509.Certificate {
	notBefore := time.Now()
	notAfter := notBefore.Add(validFor)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		panic(err)
	}

	return x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Acme Co"},
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}
}

func loadEcdsaKeyFromFile(filePath string) ecdsa.PrivateKey {
	keyPairs := LoadKeyPairsFromExistingFile(filePath)
	if len(keyPairs) == 0 {
		panic("No ECDSA private keys found in file!")
	}

	if len(keyPairs) > 1 {
		panic("Multiple ECDSA private keys found in file!")
	}

	return keyPairs[0].PrivateKey
}

func createCertFromKeyFile(path, keyFile string, validFor time.Duration) {
	fmt.Printf("Creating cert file..\n")

	template := newTemplate(validFor)
	privateKey := loadEcdsaKeyFromFile(keyFile)

	certBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, &privateKey)
	if err != nil {
		panic(err)
	}

	certOut, err := os.Create(path)
	if err != nil {
		panic(err)
	}

	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: certBytes}); err != nil {
		panic(err)
	}

	if err := certOut.Close(); err != nil {
		panic(err)
	}

	fmt.Printf("Done.\n")
}

func (genCert *genCertCmd) Run() error {
	if !acceptDisclaimer(disclaimer, retry) {
		os.Exit(0)
	}

	duration := parseDuration(genCert.ValidFor)
	createCertFromKeyFile(genCert.Out, genCert.Key, duration)

	return nil
}
