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

func (cli *CLI) createWallet(nodeID string) {
	wallets, _ := blockchain.NewWallets(nodeID)
	address := wallets.CreateWallet()

	wallets.SaveToFile(nodeID)

	fmt.Printf("Your new address: %s\n", address)
}
