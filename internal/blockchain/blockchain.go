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
	"log"
	"os"
)

const dbFile = "golang_blockweb_%s.db"
const genesisCoinbaseData = "she isn't human; she is art, with a heart."

type Blockchain struct {
	tip []byte
	db  *bolt.DB
}

func (blockchain *Blockchain) CloseDB() {
	blockchain.db.Close()
}

func (blockchain *Blockchain) GetBlock(blockHash []byte) (Block, error) {
	var block Block

	// TODO: Error handling
	blockchain.db.View(func(transaction *bolt.Tx) error {
		bucket := transaction.Bucket([]byte(persistence.BlocksBucket))

		blockData := bucket.Get(blockHash)

		if blockData == nil {
			return errors.New("Block was not found.")
		}

		block = *DeserializeBlock(blockData)

		return nil
	})

	return block, nil
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
	if transaction.IsCoinbase() {
		return true
	}

	previousTransactions := make(map[string]Transaction)

	for _, transactionInput := range transaction.TransactionInputs {
		// TODO: Error handling
		previousTransaction, _ := blockchain.FindTransaction(transactionInput.TransactionID)
		previousTransactions[hex.EncodeToString(previousTransaction.ID)] = previousTransaction
	}

	return transaction.Verify(previousTransactions)
}

func (blockchain *Blockchain) FindUnspentTransactionOutputs() map[string]TransactionOutputSet {
	unspentTransactionOutput := make(map[string]TransactionOutputSet)
	spentTransactionOutputs := make(map[string][]int)
	blockchainIterator := blockchain.Iterator()

	for {
		block := blockchainIterator.Next()

		for _, transaction := range block.Transactions {
			transactionID := hex.EncodeToString(transaction.ID)

		Outputs:
			for outputIterator, transactionOutput := range transaction.TransactionOutputs {
				if spentTransactionOutputs[transactionID] != nil {
					for _, spentOutputIndex := range spentTransactionOutputs[transactionID] {
						if spentOutputIndex == outputIterator {
							continue Outputs
						}
					}
				}

				transactionOutputs := unspentTransactionOutput[transactionID]
				transactionOutputs.TransactionOutputs = append(transactionOutputs.TransactionOutputs, transactionOutput)
				unspentTransactionOutput[transactionID] = transactionOutputs
			}

			if transaction.IsCoinbase() == false {
				for _, transactionInput := range transaction.TransactionInputs {
					transactionInputID := hex.EncodeToString(transactionInput.TransactionID)
					spentTransactionOutputs[transactionInputID] = append(spentTransactionOutputs[transactionInputID], transactionInput.transactionOutputID)
				}
			}
		}

		if len(block.PreviousHash) == 0 {
			break
		}
	}

	return unspentTransactionOutput
}

func (blockchain *Blockchain) AddBlock(block *Block) {
	//TODO: Error handling
	blockchain.db.Update(func(transaction *bolt.Tx) error {
		bucket := transaction.Bucket([]byte(persistence.BlocksBucket))
		blockInDb := bucket.Get(block.Hash)

		if blockInDb != nil {
			return nil
		}

		blockData := block.Serialize()
		// TODO: Error handling
		bucket.Put(block.Hash, blockData)

		lastHash := bucket.Get([]byte(persistence.LastBlockFileNumber))
		lastBlockData := bucket.Get(lastHash)
		lastBlock := DeserializeBlock(lastBlockData)

		if block.Height > lastBlock.Height {
			// TODO: Error handling
			bucket.Put([]byte(persistence.LastBlockFileNumber), block.Hash)
			blockchain.tip = block.Hash
		}

		return nil
	})
}

func (blockchain *Blockchain) MineBlock(transactions []*Transaction) *Block {
	var lastHash []byte
	var lastHeight int

	for _, transaction := range transactions {
		// TODO: Ignore transaction if it's not valid

		if blockchain.VerifyTransaction(transaction) != true {
			log.Panic("ERROR: Invalid transaction")
		}
	}

	// TODO: Error handling
	blockchain.db.View(func(transaction *bolt.Tx) error {
		bucket := transaction.Bucket([]byte(persistence.BlocksBucket))
		lastHash = bucket.Get([]byte(persistence.LastBlockFileNumber))

		blockData := bucket.Get(lastHash)
		block := DeserializeBlock(blockData)

		lastHeight = block.Height

		return nil
	})

	newBlock := NewBlock(transactions, lastHash, lastHeight+1)

	// TODO: Error handling
	blockchain.db.Update(func(transaction *bolt.Tx) error {
		bucket := transaction.Bucket([]byte(persistence.BlocksBucket))
		// TODO: Error handling
		bucket.Put(newBlock.Hash, newBlock.Serialize())

		// TODO: Error handling
		bucket.Put([]byte(persistence.LastBlockFileNumber), newBlock.Hash)

		blockchain.tip = newBlock.Hash

		return nil
	})

	return newBlock
}

func (blockchain *Blockchain) GetBestHeight() int {
	var lastBlock Block

	// TODO: Error handling
	blockchain.db.View(func(transaction *bolt.Tx) error {
		bucket := transaction.Bucket([]byte(persistence.BlocksBucket))
		lastHash := bucket.Get([]byte(persistence.LastBlockFileNumber))
		blockData := bucket.Get(lastHash)
		lastBlock = *DeserializeBlock(blockData)

		return nil
	})

	return lastBlock.Height
}

func (blockchain *Blockchain) GetBlockHashes() [][]byte {
	var blockHashes [][]byte
	blockchainIterator := blockchain.Iterator()

	for {
		block := blockchainIterator.Next()
		blockHashes = append(blockHashes, block.Hash)

		if len(block.PreviousHash) == 0 {
			break
		}
	}

	return blockHashes
}

func CreateBlockchain(address, nodeID string) *Blockchain {
	dbFile := fmt.Sprintf(dbFile, nodeID)
	if dbExists(dbFile) {
		// TODO: Maybe use log.panic
		fmt.Println("Blockchain already exists.")
		os.Exit(1)
	}

	var tip []byte

	coinbaseTransaction := NewCoinbaseTransaction([]byte(address))
	genesisBlock := NewGenesisBlock(coinbaseTransaction)

	// TODO: Error handling
	db, _ := bolt.Open(dbFile, 0600, nil)

	// TODO: Error handling
	db.Update(func(transaction *bolt.Tx) error {
		// TODO: Error handling
		bucket, _ := transaction.CreateBucket([]byte(persistence.BlocksBucket))

		// TODO: Error handling
		bucket.Put(genesisBlock.Hash, genesisBlock.Serialize())

		// TODO: Error handling
		bucket.Put([]byte(persistence.LastBlockFileNumber), genesisBlock.Hash)
		tip = genesisBlock.Hash

		return nil
	})

	return &Blockchain{tip: tip, db: db}
}

func NewBlockchain(nodeID string) *Blockchain {
	dbFile := fmt.Sprintf(dbFile, nodeID)
	if dbExists(dbFile) == false {
		// TODO: Maybe use log.panic
		fmt.Println("No existing Blockchain found. Create one first.")
		os.Exit(1)
	}

	var tip []byte
	// TODO: Error handling
	db, _ := bolt.Open(dbFile, 0600, nil)

	// TODO: Error handling
	db.Update(func(transaction *bolt.Tx) error {
		bucket := transaction.Bucket([]byte(persistence.BlocksBucket))
		tip = bucket.Get([]byte(persistence.LastBlockFileNumber))

		return nil
	})

	return &Blockchain{tip, db}
}

func dbExists(dbFile string) bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}

	return true
}

func (blockchain *Blockchain) Iterator() *BlockchainIterator {
	blockchainIterator := &BlockchainIterator{currentHash: blockchain.tip, db: blockchain.db}

	return blockchainIterator
}
