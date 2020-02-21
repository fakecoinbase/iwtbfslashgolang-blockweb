package network

import (
	"bytes"
	"encoding/gob"
	"github.com/iwtbf/golang-blockweb/internal/blockchain"
)

/*
 * Copyright 2020 Information Wants To Be Free
 * Visit: https://github.com/iwtbf
 *
 * This project is licensed under the terms of the Apache 2.0 License.
 */

type getBlocks struct {
	AddressFrom string
}

func sendGetBlocks(address string) {
	payload := gobEncode(getBlocks{AddressFrom: nodeAddress})
	request := append(commandToBytes("getblocks"), payload...)

	sendData(address, request)
}

func handleGetBlocks(request []byte, chain *blockchain.Blockchain) {
	var buffer bytes.Buffer
	var payload getBlocks

	buffer.Write(request[commandLength:])
	decoder := gob.NewDecoder(&buffer)
	// TODO: Error handling
	decoder.Decode(&payload)

	blockHashes := chain.GetBlockHashes()
	sendInventory(payload.AddressFrom, "block", blockHashes)
}
