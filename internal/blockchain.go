package internal

/**
 * Copyright 2019 Information Wants To Be Free
 * Visit: https://github.com/iwtbf
 *
 * This project is licensed under the terms of the Apache 2.0 License.
 */

type Blockchain struct {
	blocks []*Block
}

func (blockchain *Blockchain) AddBlock(data string) {
	prevBlock := blockchain.blocks[len(blockchain.blocks)-1]
	newBlock := NewBlock(data, prevBlock.Hash)
	blockchain.blocks = append(blockchain.blocks, newBlock)
}

func newGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}

func NewBlockchain() *Blockchain {
	return &Blockchain{[]*Block{newGenesisBlock()}}
}
