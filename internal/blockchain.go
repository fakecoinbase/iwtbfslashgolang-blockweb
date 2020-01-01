package internal

/*
 * Copyright 2019 Information Wants To Be Free
 * Visit: https://github.com/iwtbf
 *
 * This project is licensed under the terms of the Apache 2.0 License.
 */

import (
	"./persistence/BlocksBucketKeys"
	"github.com/boltdb/bolt"
	"os"
)

const blocksBucket = "blocks"

type Blockchain struct {
	tip []byte
	db  *bolt.DB
}

func (blockchain *Blockchain) AddBlock(data string) {
	var lastHash []byte

	// TODO: Error handling
	err := blockchain.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		lastHash = bucket.Get([]byte(BlocksBucketKeys.LastBlockFileNumber))

		return nil
	})

	newBlock := NewBlock(data, lastHash)

	// TODO: Error handling
	err = blockchain.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		// TODO: Error handling
		err := b.Put(newBlock.Hash, newBlock.Serialize())
		err = b.Put([]byte(BlocksBucketKeys.LastBlockFileNumber), newBlock.Hash)
		blockchain.tip = newBlock.Hash

		return nil
	})
}

func newGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}

func NewBlockchain() *Blockchain {
	var tip []byte
	// TODO: Error handling
	db, err := bolt.Open(dbFile, os.FileMode(0600), nil)

	// TODO: Error handling
	err := db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))

		if bucket == nil {
			genesis := newGenesisBlock()
			// TODO: Error handling
			newBucket, err := tx.CreateBucket([]byte(blocksBucket))
			err = newBucket.Put(genesis.Hash, genesis.Serialize())
			err = newBucket.Put([]byte(BlocksBucketKeys.LastBlockFileNumber), genesis.Hash)
			tip = genesis.Hash
		} else {
			tip = bucket.Get([]byte(BlocksBucketKeys.LastBlockFileNumber))
		}

		return nil
	})

	bc := Blockchain{tip, db}

	return &bc
}
