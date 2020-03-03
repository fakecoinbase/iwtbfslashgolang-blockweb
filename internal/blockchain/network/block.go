package network

/*
 * Copyright 2020 Information Wants To Be Free
 * Visit: https://github.com/iwtbf
 *
 * This project is licensed under the terms of the Apache 2.0 License.
 */

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/iwtbf/golang-blockweb/internal/blockchain"
	"github.com/iwtbf/golang-blockweb/internal/blockchain/network/command"
	inv "github.com/iwtbf/golang-blockweb/internal/blockchain/network/inventory"
)

type block struct {
	AddressFrom string
	Block       []byte
}

func requestBlocks() {
	for _, knownNode := range knownNodes {
		sendGetBlocks(knownNode)
	}
}

func sendBlock(address string, blockchainBlock *blockchain.Block) {
	data := block{AddressFrom: nodeAddress, Block: blockchainBlock.Serialize()}
	payload := gobEncode(data)
	request := append(commandToBytes(command.Block), payload...)

	sendData(address, request)
}

func handleBlock(request []byte, chain *blockchain.Blockchain) {
	var buffer bytes.Buffer
	var payload block

	buffer.Write(request[commandLength:])
	decoder := gob.NewDecoder(&buffer)
	// TODO: Error handling
	decoder.Decode(&payload)

	blockData := payload.Block
	block := blockchain.DeserializeBlock(blockData)

	chain.AddBlock(block)

	fmt.Printf("Added block %x\n", block.Hash)

	if len(blocksInTransit) > 0 {
		blockHash := blocksInTransit[0]
		sendGetData(payload.AddressFrom, inv.Block, blockHash)

		blocksInTransit = blocksInTransit[1:]
	} else {
		unspentTransactionOutputSet := blockchain.UnspentTransactionOutputSet{Blockchain: chain}
		unspentTransactionOutputSet.Reindex()
	}
}
