package internal

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
	PreviousHash []byte
	Timestamp    int64
	Data         []byte
	Hash         []byte
	Nonce        int
}

func (block *Block) SetHash() {
	timestamp := []byte(strconv.FormatInt(block.Timestamp, 10))
	headers := bytes.Join([][]byte{block.PreviousHash, block.Data, timestamp}, []byte{})
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

// TODO: Validate block (e.g. verify ppk validity)

// TODO: data [string] might not be applicable
func NewBlock(data string, previousHash []byte) *Block {
	block := &Block{previousHash, time.Now().Unix(), []byte(data), []byte{}, 0}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

func DeserializeBlock(data []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(data))
	// TODO: Error handling
	decoder.Decode(&block)

	return &block
}
