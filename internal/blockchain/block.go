package blockchain

/*
 * Copyright 2020 Information Wants To Be Free
 * Visit: https://github.com/iwtbf
 *
 * This project is licensed under the terms of the Apache 2.0 License.
 */

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"strconv"
	"time"
)

type Block struct {
	Timestamp    int64
	Transactions []*Transaction
	PreviousHash []byte
	Hash         []byte
	Nonce        int
	Height       int
}

func (b *Block) HashTransactions() []byte {
	var transactionHashes [][]byte
	var transactionHash [32]byte

	for _, transaction := range b.Transactions {
		transactionHashes = append(transactionHashes, transaction.ID)
	}
	transactionHash = sha256.Sum256(bytes.Join(transactionHashes, []byte{}))

	return transactionHash[:]
}

func (block *Block) SetHash() {
	timestamp := []byte(strconv.FormatInt(block.Timestamp, 10))
	headers := bytes.Join([][]byte{block.PreviousHash, block.HashTransactions(), timestamp}, []byte{})
	hash := sha256.Sum256(headers)

	block.Hash = hash[:]
}

// TODO: Maybe use protocol buffers
func (block *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	// TODO: Error handling
	encoder.Encode(block)

	return result.Bytes()
}

func DeserializeBlock(data []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(data))
	// TODO: Error handling
	decoder.Decode(&block)

	return &block
}

// TODO: Validate block (e.g. verify ppk validity)

// TODO: data [string] might not be applicable
func NewBlock(transactions []*Transaction, previousHash []byte, height int) *Block {
	block := &Block{Timestamp: time.Now().Unix(), Transactions: transactions, PreviousHash: previousHash, Hash: []byte{}, Nonce: 0, Height: height}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

func NewGenesisBlock(coinbase *Transaction) *Block {
	return NewBlock([]*Transaction{coinbase}, []byte{}, 0)
}
