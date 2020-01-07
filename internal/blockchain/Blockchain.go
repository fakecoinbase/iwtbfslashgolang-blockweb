package blockchain

/*
 * Copyright 2020 Information Wants To Be Free
 * Visit: https://github.com/iwtbf
 *
 * This project is licensed under the terms of the Apache 2.0 License.
 */

import (
	"../persistence"
	"encoding/hex"
	"fmt"
	"github.com/boltdb/bolt"
	"os"
)

type Blockchain struct {
	tip []byte
	DB  *bolt.DB
}

//func (blockchain *Blockchain) AddBlock(transactions []*Transaction) {
//	var previousHash []byte
//
//	// TODO: Error handling
//	blockchain.DB.View(func(tx *bolt.Tx) error {
//		bucket := tx.Bucket([]byte(persistence.Buckets.BlocksBucket))
//		previousHash = bucket.Get([]byte(persistence.BlocksBucketKeys.LastBlockFileNumber))
//
//		return nil
//	})
//
//	newBlock := NewBlock(transactions, previousHash)
//
//	// TODO: Error handling
//	blockchain.DB.Update(func(tx *bolt.Tx) error {
//		b := tx.Bucket([]byte(persistence.Buckets.BlocksBucket))
//		// TODO: Error handling
//		b.Put(newBlock.Hash, newBlock.Serialize())
//		b.Put([]byte(persistence.BlocksBucketKeys.LastBlockFileNumber), newBlock.Hash)
//		blockchain.tip = newBlock.Hash
//
//		return nil
//	})
//}

func (blockchain *Blockchain) FindUnspentTransactions(address string) []Transaction {
	var unspentTransactions []Transaction
	spentTransactions := make(map[string][]int)
	blockchainIterator := blockchain.Iterator()

	for {
		block := blockchainIterator.Next()

		for _, transaction := range block.Transactions {
			transactionID := hex.EncodeToString(transaction.ID)

		Outputs:
			for outputIndex, transactionOutput := range transaction.Vout {
				if spentTransactions[transactionID] != nil {
					for _, spentOutput := range spentTransactions[transactionID] {
						if spentOutput == outputIndex {
							continue Outputs
						}
					}
				}

				if transactionOutput.CanUnlockUsing(address) {
					unspentTransactions = append(unspentTransactions, *transaction)
				}
			}

			if transaction.IsCoinbase() == false {
				for _, transactionInput := range transaction.Vin {
					if transactionInput.CanUnlockUsing(address) {
						inputTransactionID := hex.EncodeToString(transactionInput.TransactionID)
						spentTransactions[inputTransactionID] = append(spentTransactions[inputTransactionID], transactionInput.Vout)
					}
				}
			}
		}

		if len(block.PreviousHash) == 0 {
			break
		}
	}

	return unspentTransactions
}

func (blockchain *Blockchain) FindUnspentTransactionOutputs(address string) []TransactionOutput {
	var unspentTransactionOutputs []TransactionOutput
	unspentTransactions := blockchain.FindUnspentTransactions(address)

	for _, unspentTransaction := range unspentTransactions {
		for _, transactionOutput := range unspentTransaction.Vout {
			if transactionOutput.CanUnlockUsing(address) {
				unspentTransactionOutputs = append(unspentTransactionOutputs, transactionOutput)
			}
		}
	}

	return unspentTransactionOutputs
}

func NewBlockchain(dbFilePath, address string) *Blockchain {
	var tip []byte
	db, err := bolt.Open(dbFilePath, os.FileMode(0600), nil)
	if err != nil {
		// TODO: Use logger
		fmt.Printf("Error opening Database: %s\n", err)
		os.Exit(1)
	}

	// TODO: Error handling
	db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(persistence.Buckets.BlocksBucket))

		if bucket == nil {
			coinbaseTransaction := NewCoinbaseTransaction(address)
			genesis := NewGenesisBlock(coinbaseTransaction)
			// TODO: Error handling
			newBucket, _ := tx.CreateBucket([]byte(persistence.Buckets.BlocksBucket))
			newBucket.Put(genesis.Hash, genesis.Serialize())
			newBucket.Put([]byte(persistence.BlocksBucketKeys.LastBlockFileNumber), genesis.Hash)
			tip = genesis.Hash
		} else {
			tip = bucket.Get([]byte(persistence.BlocksBucketKeys.LastBlockFileNumber))
		}

		return nil
	})

	blockchain := Blockchain{tip, db}

	return &blockchain
}

func (blockchain *Blockchain) Iterator() *BlockchainIterator {
	blockchainIterator := &BlockchainIterator{blockchain.tip, blockchain.DB}

	return blockchainIterator
}
