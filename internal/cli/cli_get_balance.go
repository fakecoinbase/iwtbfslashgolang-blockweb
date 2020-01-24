package cli

/*
 * Copyright 2020 Information Wants To Be Free
 * Visit: https://github.com/iwtbf
 *
 * This project is licensed under the terms of the Apache 2.0 License.
 */

import (
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"github.com/iwtbf/golang-blockweb/internal/blockchain"
)

func (cli *CLI) getBalance(address, nodeID string) {
	// TODO: Validate address

	chain := blockchain.NewBlockchain(nodeID)
	unspentTransactionOutputSet := blockchain.UnspentTransactionOutputSet{chain}
	defer chain.CloseDB()

	balance := 0
	pubKeyHash := base58.Decode(address)
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-blockchain.AddressChecksumLength]
	unspentTransactionOutputs := unspentTransactionOutputSet.FindUnspentTransactionOutputs(pubKeyHash)

	for _, unspentTransactionOutput := range unspentTransactionOutputs {
		balance += unspentTransactionOutput.Value
	}

	fmt.Printf("Balance of '%s': %d\n", address, balance)
}
