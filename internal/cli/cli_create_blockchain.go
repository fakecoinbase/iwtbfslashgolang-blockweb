package cli

/*
 * Copyright 2020 Information Wants To Be Free
 * Visit: https://github.com/iwtbf
 *
 * This project is licensed under the terms of the Apache 2.0 License.
 */

import (
	"fmt"
	"github.com/iwtbf/golang-blockweb/internal/blockchain"
)

func (cli *CLI) createBlockchain(address, nodeID string) {
	// TODO: Validate address

	chain := blockchain.CreateBlockchain(address, nodeID)
	chain.CloseDB()

	unspentTransactionOutputSet := blockchain.UnspentTransactionOutputSet{chain}
	unspentTransactionOutputSet.Reindex()

	fmt.Println("Done!")
}
