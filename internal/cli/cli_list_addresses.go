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

func (cli *CLI) listAddresses(nodeID string) {
	// TODO: Error handling
	wallets, _ := blockchain.NewWallets(nodeID)
	addresses := wallets.GetAddresses()

	for _, address := range addresses {
		fmt.Println(address)
	}
}
