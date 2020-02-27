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
)

func newKeyPair() {
	curve := elliptic.P256()
	// TODO: Error handling
	privateKey, _ := ecdsa.GenerateKey(curve, rand.Reader)
	//publicKey := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)

	// TODO: Error handling
	pkcs8Formatted, _ := x509.MarshalPKCS8PrivateKey(privateKey)
	fmt.Printf("PKCS8 Private Key:\n%s\n", string(pkcs8Formatted))

	x509Encoded, _ := x509.MarshalECPrivateKey(privateKey)
	fmt.Printf("x509 Private Key:\n%s\n", string(x509Encoded))

	pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: x509Encoded})
	fmt.Printf("PEM Private Key:\n%s\n", string(pemEncoded))

	x509EncodedPub, _ := x509.MarshalPKIXPublicKey(privateKey.PublicKey)
	fmt.Printf("Px509 Public Key:\n%s\n", string(x509EncodedPub))

	pemEncodedPub := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: x509EncodedPub})
	fmt.Printf("PEM Public Key:\n%s\n", string(pemEncodedPub))
}
