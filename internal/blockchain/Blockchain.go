package blockchain

/*
 * Copyright 2020 Information Wants To Be Free
 * Visit: https://github.com/iwtbf
 *
 * This project is licensed under the terms of the Apache 2.0 License.
 */

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/iwtbf/golang-blockweb/internal/persistence"
	"os"
)

const genesisCoinbaseData = "she isn't human; she is art, with a heart."

type Blockchain struct {
	tip []byte
	db  *bolt.DB
}

func (blockchain *Blockchain) MineBlock(transactions []*Transaction) {
	var previousHash []byte

	// TODO: Error handling
	blockchain.db.View(func(boltTransaction *bolt.Tx) error {
		bucket := boltTransaction.Bucket([]byte(persistence.BlocksBucket))
		previousHash = bucket.Get([]byte("l"))

		return nil
	})

	newBlock := NewBlock(transactions, previousHash)

	// TODO: Error handling
	blockchain.db.Update(func(boltTransaction *bolt.Tx) error {
		bucket := boltTransaction.Bucket([]byte(persistence.BlocksBucket))
		// TODO: Error handling
		bucket.Put(newBlock.Hash, newBlock.Serialize())

		// TODO: Error handling
		bucket.Put([]byte("l"), newBlock.Hash)

		blockchain.tip = newBlock.Hash

		return nil
	})
}

func (blockchain *Blockchain) FindTransaction(ID []byte) (Transaction, error) {
	blockchainIterator := blockchain.Iterator()

	for {
		block := blockchainIterator.Next()

		for _, transaction := range block.Transactions {
			if bytes.Compare(transaction.ID, ID) == 0 {
				return *transaction, nil
			}
		}

		if len(block.PreviousHash) == 0 {
			break
		}
	}

	return Transaction{}, errors.New("Transaction was not found")
}

func (blockchain *Blockchain) SignTransaction(transaction *Transaction, privateKey ecdsa.PrivateKey) {
	previousTransactions := make(map[string]Transaction)

	for _, transactionInput := range transaction.TransactionInputs {
		// TODO: Error handling
		previousTransaction, _ := blockchain.FindTransaction(transactionInput.TransactionID)
		previousTransactions[hex.EncodeToString(previousTransaction.ID)] = previousTransaction
	}

	transaction.Sign(privateKey, previousTransactions)
}

func (blockchain *Blockchain) VerifyTransaction(transaction *Transaction) bool {
	prevTXs := make(map[string]Transaction)

	for _, transactionInput := range transaction.TransactionInputs {
		// TODO: Error handling
		prevTX, _ := blockchain.FindTransaction(transactionInput.TransactionID)
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}

	return transaction.Verify(prevTXs)
}

func (blockchain *Blockchain) FindSpendableOutputs(publicKeyHash []byte, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	unspentTransactions := blockchain.FindUnspentTransactions(publicKeyHash)
	accumulated := 0

Work:
	for _, unspentTransaction := range unspentTransactions {
		transactionID := hex.EncodeToString(unspentTransaction.ID)

		for outputIterator, transactionOutput := range unspentTransaction.TransactionOutputs {
			if transactionOutput.IsLockedWithKey(publicKeyHash) && accumulated < amount {
				accumulated += transactionOutput.Value
				unspentOutputs[transactionID] = append(unspentOutputs[transactionID], outputIterator)

				if accumulated >= amount {
					break Work
				}
			}
		}
	}

	return accumulated, unspentOutputs
}

func (blockchain *Blockchain) FindUnspentOutputs(pubKeyHash []byte) []TransactionOutput {
	var unspentTransactionOutputs []TransactionOutput
	unspentTransactions := blockchain.FindUnspentTransactions(pubKeyHash)

	for _, unspentTransaction := range unspentTransactions {
		for _, transactionOutput := range unspentTransaction.TransactionOutputs {
			if transactionOutput.IsLockedWithKey(pubKeyHash) {
				unspentTransactionOutputs = append(unspentTransactionOutputs, transactionOutput)
			}
		}
	}

	return unspentTransactionOutputs
}

func (blockchain *Blockchain) FindUnspentTransactions(pubKeyHash []byte) []Transaction {
	var unspentTransactions []Transaction
	spentTransactionOutputs := make(map[string][]int)
	blockchainIterator := blockchain.Iterator()

	for {
		block := blockchainIterator.Next()

		for _, transaction := range block.Transactions {
			transactionID := hex.EncodeToString(transaction.ID)

		Outputs:
			for outputIterator, transactionOutput := range transaction.TransactionOutputs {
				// Was the output spent?
				if spentTransactionOutputs[transactionID] != nil {
					for _, spentOutIdx := range spentTransactionOutputs[transactionID] {
						if spentOutIdx == outputIterator {
							continue Outputs
						}
					}
				}

				if transactionOutput.IsLockedWithKey(pubKeyHash) {
					unspentTransactions = append(unspentTransactions, *transaction)
				}
			}

			if transaction.IsCoinbase() == false {
				for _, transactionInput := range transaction.TransactionInputs {
					if transactionInput.UsesKey(pubKeyHash) {
						transactionInputID := hex.EncodeToString(transactionInput.TransactionID)
						spentTransactionOutputs[transactionInputID] = append(spentTransactionOutputs[transactionInputID], transactionInput.transactionOutputID)
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

func (blockchain *Blockchain) CloseDB() {
	blockchain.db.Close()
}

func NewBlockchain(dbFilePath, genesisAddress string) *Blockchain {
	var tip []byte
	db, err := bolt.Open(dbFilePath, os.FileMode(0600), nil)
	if err != nil {
		// TODO: Use logger
		fmt.Printf("Error opening Database: %s\n", err)
		os.Exit(1)
	}

	// TODO: Error handling
	db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(persistence.BlocksBucket))

		if bucket == nil {
			coinbaseTransaction := NewCoinbaseTransaction(genesisAddress, genesisCoinbaseData)
			genesis := NewGenesisBlock(coinbaseTransaction)
			// TODO: Error handling
			newBucket, _ := tx.CreateBucket([]byte(persistence.BlocksBucket))
			newBucket.Put(genesis.Hash, genesis.Serialize())
			newBucket.Put([]byte(persistence.LastBlockFileNumber), genesis.Hash)
			tip = genesis.Hash
		} else {
			tip = bucket.Get([]byte(persistence.LastBlockFileNumber))
		}

		return nil
	})

	blockchain := Blockchain{tip: tip, db: db}

	return &blockchain
}

func (blockchain *Blockchain) Iterator() *BlockchainIterator {
	blockchainIterator := &BlockchainIterator{blockchain.tip, blockchain.db}

	return blockchainIterator
}
