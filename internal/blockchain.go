package internal

/*
 * Copyright 2020 Information Wants To Be Free
 * Visit: https://github.com/iwtbf
 *
 * This project is licensed under the terms of the Apache 2.0 License.
 */

import (
	"./persistence/BlocksBucketKeys"
	"fmt"
	"github.com/boltdb/bolt"
	"os"
)

const blocksBucket = "blocks"

type Blockchain struct {
	tip []byte
	DB  *bolt.DB
}

func (blockchain *Blockchain) AddBlock(data string) {
	var lastHash []byte

	// TODO: Error handling
	blockchain.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		lastHash = bucket.Get([]byte(BlocksBucketKeys.LastBlockFileNumber))

		return nil
	})

	newBlock := NewBlock(data, lastHash)

	// TODO: Error handling
	blockchain.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		// TODO: Error handling
		b.Put(newBlock.Hash, newBlock.Serialize())
		b.Put([]byte(BlocksBucketKeys.LastBlockFileNumber), newBlock.Hash)
		blockchain.tip = newBlock.Hash

		return nil
	})
}

func newGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}

func NewBlockchain(dbFilePath string) *Blockchain {
	var tip []byte
	db, err := bolt.Open(dbFilePath, os.FileMode(0600), nil)
	if err != nil {
		// TODO: Use logger
		fmt.Printf("Error opening Database: %s\n", err)
		os.Exit(1)
	}

	// TODO: Error handling
	db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))

		if bucket == nil {
			genesis := newGenesisBlock()
			// TODO: Error handling
			newBucket, _ := tx.CreateBucket([]byte(blocksBucket))
			newBucket.Put(genesis.Hash, genesis.Serialize())
			newBucket.Put([]byte(BlocksBucketKeys.LastBlockFileNumber), genesis.Hash)
			tip = genesis.Hash
		} else {
			tip = bucket.Get([]byte(BlocksBucketKeys.LastBlockFileNumber))
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
