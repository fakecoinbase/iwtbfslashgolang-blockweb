package keygen

/*
 * Copyright 2020 Information Wants To Be Free
 * Visit: https://github.com/iwtbf
 *
 * This project is licensed under the terms of the Apache 2.0 License.
 */

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
	"io/ioutil"
	"os"
)

const (
	version               = byte(0x00)
	addressChecksumLength = 4
)

type KeyPair struct {
	PrivateKey ecdsa.PrivateKey
	publicKey  []byte
}

func (keyPair *KeyPair) hashPublicKey() []byte {
	publicSHA256 := sha256.Sum256(keyPair.publicKey)

	RIPEMD160Hasher := ripemd160.New()

	_, err := RIPEMD160Hasher.Write(publicSHA256[:])
	if err != nil {
		panic(err)
	}

	return RIPEMD160Hasher.Sum(nil)
}

func (keyPair *KeyPair) publicBlockchainAddress() string {
	versionedPayload := append([]byte{version}, keyPair.hashPublicKey()...)
	checksum := checksum(versionedPayload)

	return base58.Encode(append(versionedPayload, checksum...))
}

func checksum(versionedPublicKey []byte) []byte {
	firstSHA := sha256.Sum256(versionedPublicKey)
	secondSHA := sha256.Sum256(firstSHA[:])

	return secondSHA[:addressChecksumLength]
}

func validateReadableFile(path string) {
	_, err := os.Stat(path)
	if err != nil {
		panic("Key does not exist!")
	}
}

func LoadKeyPairsFromExistingFile(path string) []KeyPair {
	validateReadableFile(path)

	var keyPairs []KeyPair

	privateKeys, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	for block, rest := pem.Decode(privateKeys); block != nil && rest != nil; block, rest = pem.Decode(rest) {
		x509Encoded := block.Bytes
		privateKey, err := x509.ParseECPrivateKey(x509Encoded)
		if err != nil {
			panic(err)
		}

		publicKey := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)
		readKeyPair := KeyPair{PrivateKey: *privateKey, publicKey: publicKey}

		keyPairs = append(keyPairs, readKeyPair)
	}

	return keyPairs
}
