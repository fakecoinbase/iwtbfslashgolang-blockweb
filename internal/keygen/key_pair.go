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
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
	"io/ioutil"
	"os"
)

const (
	version               = byte(0x00)
	AddressChecksumLength = 4
)

type keyPair struct {
	privateKey ecdsa.PrivateKey
	publicKey  []byte
}

func (keyPair *keyPair) hashPublicKey() []byte {
	publicSHA256 := sha256.Sum256(keyPair.publicKey)

	RIPEMD160Hasher := ripemd160.New()

	_, err := RIPEMD160Hasher.Write(publicSHA256[:])
	if err != nil {
		panic(err)
	}

	return RIPEMD160Hasher.Sum(nil)
}

func (keyPair *keyPair) publicBlockchainAddress() string {
	versionedPayload := append([]byte{version}, keyPair.hashPublicKey()...)
	checksum := checksum(versionedPayload)

	return base58.Encode(append(versionedPayload, checksum...))
}

func checksum(versionedPublicKey []byte) []byte {
	firstSHA := sha256.Sum256(versionedPublicKey)
	secondSHA := sha256.Sum256(firstSHA[:])

	return secondSHA[:AddressChecksumLength]
}

func shouldOverrideFile(path string) bool {
	question := fmt.Sprintf("The file '%s' does already exist? Do you want to override it (y / N)?", path)
	return acceptDisclaimer(question, disclaimer)
}

func isWriteableFile(fileInfo os.FileInfo) bool {
	// TODO: Check file permissions
	return !fileInfo.IsDir()
}

func validateFilePath(path string) {
	fileInfo, err := os.Stat(path)

	if err == nil && !shouldOverrideFile(path) {
		panic("Will not override existing file!")
	}

	if err == nil && !isWriteableFile(fileInfo) {
		panic("File '" + path + "' is not writeable!")
	}
}

func newKeyPair(path string, print bool) {
	validateFilePath(path)

	curve := elliptic.P256()

	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		panic(err)
	}

	publicKey := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)

	newKeyPair := keyPair{privateKey: *privateKey, publicKey: publicKey}

	x509Encoded, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		panic(err)
	}

	pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: x509Encoded})

	println("Writing file..")

	err = ioutil.WriteFile(path, pemEncoded, os.ModeExclusive)
	if err != nil {
		panic(err)
	}

	if print {
		fmt.Printf("\n Your private Key:\n%s\n", string(pemEncoded))
	}

	blockchainAddress := newKeyPair.publicBlockchainAddress()
	fmt.Printf("\nYour public Address:\n%s\n", blockchainAddress)
}
