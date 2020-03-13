package keygen

import (
	"fmt"
)

/*
 * Copyright 2020 Information Wants To Be Free
 * Visit: https://github.com/iwtbf
 *
 * This project is licensed under the terms of the Apache 2.0 License.
 */

type readKeyCmd struct {
	File string `arg type:"existingfile" help:"Path to the input key file."`
}

func readAndPrintKeyPairs(path string) {
	fmt.Println("\nReading key pairs from file..")

	keyPairs := LoadKeyPairsFromExistingFile(path)

	if len(keyPairs) == 0 {
		fmt.Println("No ECDSA private keys found in file.")
		return
	}

	fmt.Printf("Found %d key pair(s). Public addresses:\n", len(keyPairs))

	for _, keyPair := range keyPairs {
		fmt.Printf("\t- %s\n", keyPair.publicBlockchainAddress())
	}
}

func (readKey *readKeyCmd) Run() error {
	readAndPrintKeyPairs(readKey.File)

	return nil
}
