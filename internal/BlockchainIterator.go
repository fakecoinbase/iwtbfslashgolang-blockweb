package internal

/*
 * Copyright 2020 Information Wants To Be Free
 * Visit: https://github.com/iwtbf
 *
 * This project is licensed under the terms of the Apache 2.0 License.
 */

import "github.com/boltdb/bolt"

type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

func (blockchainIterator *BlockchainIterator) Next() *Block {
	var block *Block

	// TODO: Error handling
	blockchainIterator.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		encodedBlock := bucket.Get(blockchainIterator.currentHash)
		block = DeserializeBlock(encodedBlock)

		return nil
	})

	blockchainIterator.currentHash = block.PreviousHash

	return block
}
