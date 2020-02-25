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

func (cli *CLI) createBlockchain(nodeID string) {
	// TODO: Validate address

	chain := blockchain.CreateBlockchain(nodeID)
	defer chain.CloseDB()

	unspentTransactionOutputSet := blockchain.UnspentTransactionOutputSet{Blockchain: chain}
	unspentTransactionOutputSet.Reindex()

	fmt.Println("Done!")
}
