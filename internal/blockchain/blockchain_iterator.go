package blockchain

/*
 * Copyright 2020 Information Wants To Be Free
 * Visit: https://github.com/iwtbf
 *
 * This project is licensed under the terms of the Apache 2.0 License.
 */

import (
	"github.com/boltdb/bolt"
	"github.com/iwtbf/golang-blockweb/internal/persistence"
)

type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

func (blockchainIterator *BlockchainIterator) Next() *Block {
	var block *Block

	// TODO: Error handling
	blockchainIterator.db.View(func(transaction *bolt.Tx) error {
		bucket := transaction.Bucket([]byte(persistence.BlocksBucket))
		encodedBlock := bucket.Get(blockchainIterator.currentHash)
		block = DeserializeBlock(encodedBlock)

		return nil
	})

	blockchainIterator.currentHash = block.PreviousHash

	return block
}
