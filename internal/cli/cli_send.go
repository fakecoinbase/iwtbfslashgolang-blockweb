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

func (cli *CLI) send(from, to string, amount int, nodeID string) {
	// TODO: Validate addresses

	chain := blockchain.NewBlockchain(nodeID)
	unspentTransactionOutputSet := blockchain.UnspentTransactionOutputSet{chain}
	defer chain.CloseDB()

	// TODO: Error handling
	wallets, _ := blockchain.NewWallets(nodeID)
	wallet := wallets.GetWallet(from)

	transaction := blockchain.NewTransaction(&wallet, []byte(to), amount, unspentTransactionOutputSet)

	// Implement mining
	//if mineNow {
	coinbaseTransaction := blockchain.NewCoinbaseTransaction([]byte(from), "")
	transactions := []*blockchain.Transaction{coinbaseTransaction, transaction}

	newBlock := chain.MineBlock(transactions)
	unspentTransactionOutputSet.Update(newBlock)
	//} else {
	//	sendTx(knownNodes[0], transaction)
	//}

	fmt.Println("Success!")
}