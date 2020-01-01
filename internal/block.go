package internal

/**
 * Copyright 2019 Information Wants To Be Free
 * Visit: https://github.com/iwtbf
 *
 * This project is licensed under the terms of the Apache 2.0 License.
 */

import (
	"bytes"
	"crypto/sha256"
	"strconv"
	"time"
)

type Block struct {
	PreviousHash []byte
	Timestamp    int64
	Data         []byte
	Hash         []byte
}

func (block *Block) SetHash() {
	timestamp := []byte(strconv.FormatInt(block.Timestamp, 10))
	headers := bytes.Join([][]byte{block.PreviousHash, block.Data, timestamp}, []byte{})
	hash := sha256.Sum256(headers)

	block.Hash = hash[:]
}

// TODO: data [string] might not be applicable
func NewBlock(data string, previousHash []byte) *Block {
	block := &Block{previousHash, time.Now().Unix(), []byte(data), []byte{}}
	block.SetHash()
	return block
}
