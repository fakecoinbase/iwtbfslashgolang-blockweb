package network

/*
 * Copyright 2020 Information Wants To Be Free
 * Visit: https://github.com/iwtbf
 *
 * This project is licensed under the terms of the Apache 2.0 License.
 */

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"github.com/iwtbf/golang-blockweb/internal/blockchain"
	inv "github.com/iwtbf/golang-blockweb/internal/blockchain/network/inventory"
)

type transaction struct {
	AddressFrom string
	Transaction []byte
}

func sendTransaction(address string, blockchainTransaction *blockchain.Transaction) {
	data := transaction{AddressFrom: nodeAddress, Transaction: blockchainTransaction.Serialize()}
	payload := gobEncode(data)
	request := append(commandToBytes("tx"), payload...)

	sendData(address, request)
}

func handleTransaction(request []byte, chain *blockchain.Blockchain) {
	var buffer bytes.Buffer
	var payload transaction

	buffer.Write(request[commandLength:])
	decoder := gob.NewDecoder(&buffer)
	// TODO: Error handling
	decoder.Decode(&payload)

	transactionData := payload.Transaction
	serializedTransaction := *blockchain.DeserializeTransaction(transactionData)
	mempool[hex.EncodeToString(serializedTransaction.ID)] = serializedTransaction

	if nodeAddress == knownNodes[0] {
		for _, knownNode := range knownNodes {
			if knownNode != nodeAddress && knownNode != payload.AddressFrom {
				sendInventory(knownNode, inv.Transaction, [][]byte{serializedTransaction.ID})
			}
		}
	} else {
		if len(mempool) >= 2 && len(miningAddress) > 0 {
		MineTransactions:
			var transactions []*blockchain.Transaction

			for id := range mempool {
				tx := mempool[id]
				if chain.VerifyTransaction(&tx) {
					transactions = append(transactions, &tx)
				}
			}

			if len(transactions) == 0 {
				fmt.Println("All transactions are invalid! Waiting for new ones...")
				return
			}

			coinbaseTransaction := blockchain.NewCoinbaseTransaction([]byte(miningAddress))
			transactions = append(transactions, coinbaseTransaction)

			newBlock := chain.MineBlock(transactions)
			unspentTransactionOutputSet := blockchain.UnspentTransactionOutputSet{Blockchain: chain}
			// TODO: Reindex should be Update()
			unspentTransactionOutputSet.Reindex()

			fmt.Println("New block is mined!")

			for _, transaction := range transactions {
				txID := hex.EncodeToString(transaction.ID)
				delete(mempool, txID)
			}

			for _, knownNode := range knownNodes {
				if knownNode != nodeAddress {
					sendInventory(knownNode, inv.Block, [][]byte{newBlock.Hash})
				}
			}

			if len(mempool) > 0 {
				goto MineTransactions
			}
		}
	}
}
