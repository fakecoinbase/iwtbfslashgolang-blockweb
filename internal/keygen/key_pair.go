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
	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
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
