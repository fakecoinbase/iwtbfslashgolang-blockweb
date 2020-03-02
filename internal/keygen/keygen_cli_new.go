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
	"io/ioutil"
	"os"
)

const (
	disclaimer = `DISCLAIMER:
This program * can * be used to generate key pairs for 'golang-blockweb'.
It is open source software, but there is no liability of the contributors or guarantee for the security of the generated keys. Please report any problems or vulnerabilities at https://github.com/iwtbf/golang-blockweb/issues/new.
You could use any other key tool with elliptic P256 curve cryptography (ECDSA) + x509 + PEM. Use the command 'read' to verify the generated key.
By continuing, you confirm that you have understood the above points (y / N):`
	retry = `Sorry, I didn't understand this. Please try again (y / N):`
)

type newCmd struct {
	Out   string `arg help:"Path to the output file."`
	Print bool   `flag help:"Print the *private key* to System.Out as well."`
}

func shouldOverrideFile(path string) bool {
	question := fmt.Sprintf("The file '%s' does already exist? Do you want to override it (y / N)?", path)
	return acceptDisclaimer(question, disclaimer)
}

func isWriteableFile(fileInfo os.FileInfo) bool {
	// TODO: Check file permissions
	return !fileInfo.IsDir()
}

func validateWriteableFilePath(path string) {
	fileInfo, err := os.Stat(path)

	if err == nil && !shouldOverrideFile(path) {
		panic("Will not override existing file!")
	}

	if err == nil && !isWriteableFile(fileInfo) {
		panic("File '" + path + "' is not writeable!")
	}
}

func createKeyPair(path string, print bool) {
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

	pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: x509Encoded})

	fmt.Println("Writing file..")

	if err = ioutil.WriteFile(path, pemEncoded, os.ModeExclusive); err != nil {
		panic(err)
	}

	if print {
		fmt.Printf("\n Your private Key:\n%s\n", string(pemEncoded))
	}

	blockchainAddress := newKeyPair.publicBlockchainAddress()
	fmt.Printf("\nYour public Address:\n%s\n", blockchainAddress)

	fmt.Println("\nDone.")
}

func (new *newCmd) Run() error {
	if !acceptDisclaimer(disclaimer, retry) {
		os.Exit(0)
	}

	validateWriteableFilePath(new.Out)
	createKeyPair(new.Out, new.Print)

	return nil
}
