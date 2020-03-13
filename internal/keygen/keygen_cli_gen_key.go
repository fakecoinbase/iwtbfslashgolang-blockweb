package keygen

/*
 * Copyright 2020 Information Wants To Be Free
 * Visit: https://github.com/iwtbf
 *
 * This project is licensed under the terms of the Apache 2.0 License.
 */

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

type genKeyCmd struct {
	Out string `arg help:"Path to the output key file."`
}

func createKeyPair(path string) {
	fmt.Printf("Creating key file..")

	curve := elliptic.P256()

	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		panic(err)
	}

	publicKey := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)

	newKeyPair := KeyPair{PrivateKey: *privateKey, publicKey: publicKey}

	x509Encoded, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		panic(err)
	}

	keyOut, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		panic(err)
	}

	if err := pem.Encode(keyOut, &pem.Block{Type: "PRIVATE KEY", Bytes: x509Encoded}); err != nil {
		panic(err)
	}

	blockchainAddress := newKeyPair.publicBlockchainAddress()
	fmt.Printf("\nYour public Address is:\n%s\n", blockchainAddress)

	fmt.Println("\nDone.")
}

func (genKey *genKeyCmd) Run() error {
	if !acceptDisclaimer(disclaimer, retry) {
		os.Exit(0)
	}

	createKeyPair(genKey.Out)

	return nil
}
