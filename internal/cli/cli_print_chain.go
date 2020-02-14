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
	"strconv"
)

func (cli *CLI) printChain(nodeID string) {
	chain := blockchain.NewBlockchain(nodeID)
	defer chain.CloseDB()

	blockchainIterator := chain.Iterator()

	for {
		block := blockchainIterator.Next()

		fmt.Printf("============ Block %x ============\n", block.Hash)
		fmt.Printf("Prev. block: %x\n", block.PreviousHash)

		proofOfWork := blockchain.NewProofOfWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(proofOfWork.Validate()))

		fmt.Println("Transactions:")
		for _, transaction := range block.Transactions {
			fmt.Printf("\t%x\n", transaction.ID)
		}

		fmt.Printf("\n")

		if len(block.PreviousHash) == 0 {
			break
		}
	}
}
