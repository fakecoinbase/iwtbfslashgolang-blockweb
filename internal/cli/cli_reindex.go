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

func (cli *CLI) reindexUTXO(nodeID string) {
	chain := blockchain.NewBlockchain(nodeID)
	unspentTransactionOutputSet := blockchain.UnspentTransactionOutputSet{chain}
	unspentTransactionOutputSet.Reindex()

	defer chain.CloseDB()

	amountOfTransactions := unspentTransactionOutputSet.CountTransactions()
	fmt.Printf("Done! There are %d transactions in the UTXO set.\n", amountOfTransactions)
}
