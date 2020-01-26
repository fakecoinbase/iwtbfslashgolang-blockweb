package blockchain

/*
 * Copyright 2020 Information Wants To Be Free
 * Visit: https://github.com/iwtbf
 *
 * This project is licensed under the terms of the Apache 2.0 License.
 */

import (
	"encoding/hex"
	"github.com/boltdb/bolt"
	"github.com/iwtbf/golang-blockweb/internal/persistence"
)

type UnspentTransactionOutputSet struct {
	Blockchain *Blockchain
}

func (unspentTransactionOutputSet UnspentTransactionOutputSet) FindSpendableOutputs(publicKeyHash []byte, amount int) (int, map[string][]int) {
	spendableOutputs := make(map[string][]int)
	accumulated := 0
	db := unspentTransactionOutputSet.Blockchain.db

	// TODO: Error handling
	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(persistence.ChainstateBucket))
		cursor := bucket.Cursor()

		for key, value := cursor.First(); key != nil; key, value = cursor.Next() {
			transactionID := hex.EncodeToString(key)
			transactionOutputSet := NewTransactionOutputSet(value)

			for outputIterator, transactionOutput := range transactionOutputSet.TransactionOutputs {
				if transactionOutput.IsLockedWithKey(publicKeyHash) && accumulated < amount {
					accumulated += transactionOutput.Value
					spendableOutputs[transactionID] = append(spendableOutputs[transactionID], outputIterator)
				}
			}
		}

		return nil
	})

	return accumulated, spendableOutputs
}

func (unspentTransactionOutputSet UnspentTransactionOutputSet) FindUnspentTransactionOutputs(publicKeyHash []byte) []TransactionOutput {
	var transactionOutputs []TransactionOutput
	db := unspentTransactionOutputSet.Blockchain.db

	// TODO: Error handling
	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(persistence.ChainstateBucket))
		cursor := bucket.Cursor()

		for key, value := cursor.First(); key != nil; key, value = cursor.Next() {
			transactionOutputSet := NewTransactionOutputSet(value)

			for _, transactionOutput := range transactionOutputSet.TransactionOutputs {
				if transactionOutput.IsLockedWithKey(publicKeyHash) {
					transactionOutputs = append(transactionOutputs, transactionOutput)
				}
			}
		}

		return nil
	})

	return transactionOutputs
}

func (unspentTransactionOutputSet UnspentTransactionOutputSet) CountTransactions() int {
	db := unspentTransactionOutputSet.Blockchain.db
	counter := 0

	// TODO: Error handling
	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(persistence.ChainstateBucket))
		cursor := bucket.Cursor()

		for key, _ := cursor.First(); key != nil; key, _ = cursor.Next() {
			counter++
		}

		return nil
	})

	return counter
}

func (unspentTransactionOutputSet UnspentTransactionOutputSet) Update(block *Block) {
	db := unspentTransactionOutputSet.Blockchain.db

	// TODO: Error handling
	db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(persistence.ChainstateBucket))

		for _, transaction := range block.Transactions {
			if transaction.IsCoinbase() == false {
				for _, transactionInput := range transaction.TransactionInputs {
					transactionOutputSet := TransactionOutputSet{}
					serializedTransactionOutput := bucket.Get(transactionInput.TransactionID)
					loadedTransactionOutputSet := NewTransactionOutputSet(serializedTransactionOutput)

					for outputIterator, transactionOutput := range loadedTransactionOutputSet.TransactionOutputs {
						if outputIterator != transactionInput.transactionOutputID {
							transactionOutputSet.TransactionOutputs = append(transactionOutputSet.TransactionOutputs, transactionOutput)
						}
					}

					if len(transactionOutputSet.TransactionOutputs) == 0 {
						// TODO: Error handling
						bucket.Delete(transactionInput.TransactionID)
					} else {
						// TODO: Error handling
						bucket.Put(transactionInput.TransactionID, transactionOutputSet.Serialize())
					}

				}
			}

			transactionOutputSet := TransactionOutputSet{}
			for _, transactionOutput := range transaction.TransactionOutputs {
				transactionOutputSet.TransactionOutputs = append(transactionOutputSet.TransactionOutputs, transactionOutput)
			}

			// TODO: Error handling
			bucket.Put(transaction.ID, transactionOutputSet.Serialize())
		}

		return nil
	})
}

func (unspentTransactionOutputSet UnspentTransactionOutputSet) Reindex() {
	db := unspentTransactionOutputSet.Blockchain.db
	chainstateBucketName := []byte(persistence.ChainstateBucket)

	// TODO: Error handling
	db.Update(func(tx *bolt.Tx) error {
		// TODO: Error handling
		tx.DeleteBucket(chainstateBucketName)

		// TODO: Error handling
		tx.CreateBucket(chainstateBucketName)

		return nil
	})

	unspentTransactionOutputs := unspentTransactionOutputSet.Blockchain.FindUnspentTransactionOutputs()

	// TODO: Error handling
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(chainstateBucketName)

		for transactionIterator, transactionOutputSet := range unspentTransactionOutputs {
			// TODO: Error handling
			key, _ := hex.DecodeString(transactionIterator)

			// TODO: Error handling
			b.Put(key, transactionOutputSet.Serialize())
		}

		return nil
	})
}
