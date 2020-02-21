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
	"encoding/hex"
	"github.com/iwtbf/golang-blockweb/internal/blockchain"
	inv "github.com/iwtbf/golang-blockweb/internal/blockchain/network/inventory"
)

type getData struct {
	AddressFrom string
	Type        inv.InventoryType
	ID          []byte
}

func sendGetData(address string, kind inv.InventoryType, id []byte) {
	payload := gobEncode(getData{AddressFrom: nodeAddress, Type: kind, ID: id})
	request := append(commandToBytes("getdata"), payload...)

	sendData(address, request)
}

func handleGetData(request []byte, chain *blockchain.Blockchain) {
	var buffer bytes.Buffer
	var payload getData

	buffer.Write(request[commandLength:])
	decoder := gob.NewDecoder(&buffer)
	// TODO: Error handling
	decoder.Decode(&payload)

	if payload.Type == inv.Block {
		// TODO: Error handling
		block, _ := chain.GetBlock([]byte(payload.ID))

		sendBlock(payload.AddressFrom, &block)
	}

	if payload.Type == inv.Transaction {
		transactionID := hex.EncodeToString(payload.ID)
		transaction := mempool[transactionID]

		sendTransaction(payload.AddressFrom, &transaction)
	}
}
