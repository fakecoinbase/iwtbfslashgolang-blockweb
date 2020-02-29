package keygen

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
)

/*
 * Copyright 2020 Information Wants To Be Free
 * Visit: https://github.com/iwtbf
 *
 * This project is licensed under the terms of the Apache 2.0 License.
 */

type readCmd struct {
	File string `arg help:"The path to the input file."`
}

func validateReadableFile(path string) {
	_, err := os.Stat(path)
	if err != nil {
		panic("File does not exist!")
	}
}

func readKeyPair(path string) {
	privateKey, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	for block, rest := pem.Decode(privateKey); rest != nil; block, rest = pem.Decode(rest) {
		if block == nil {
			fmt.Println("\nNo more blocks found in file!")
			os.Exit(0)
		}

		x509Encoded := block.Bytes
		privateKey, err := x509.ParseECPrivateKey(x509Encoded)
		if err != nil {
			panic(err)
		}

		publicKey := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)
		readKeyPair := keyPair{privateKey: *privateKey, publicKey: publicKey}

		blockchainAddress := readKeyPair.publicBlockchainAddress()
		fmt.Printf("\nPublic Address:\n%s\n", blockchainAddress)
	}
}

func (read *readCmd) Run() error {
	validateReadableFile(read.File)
	readKeyPair(read.File)

	return nil
}
