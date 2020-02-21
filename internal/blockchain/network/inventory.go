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
	"fmt"
	"github.com/iwtbf/golang-blockweb/internal/blockchain"
	inv "github.com/iwtbf/golang-blockweb/internal/blockchain/network/inventory"
)

type inventory struct {
	AddressFrom string
	Type        inv.InventoryType
	Items       [][]byte
}

func sendInventory(address string, kind inv.InventoryType, items [][]byte) {
	inventory := inventory{AddressFrom: nodeAddress, Type: kind, Items: items}
	payload := gobEncode(inventory)
	request := append(commandToBytes("inv"), payload...)

	sendData(address, request)
}

func handleInventory(request []byte, chain *blockchain.Blockchain) {
	var buffer bytes.Buffer
	var payload inventory

	buffer.Write(request[commandLength:])
	decoder := gob.NewDecoder(&buffer)
	// TODO: Error handling
	decoder.Decode(&payload)

	fmt.Printf("Recevied inventory with %d %s\n", len(payload.Items), payload.Type)

	if payload.Type == inv.Block {
		blocksInTransit = payload.Items

		blockHash := payload.Items[0]
		sendGetData(payload.AddressFrom, "block", blockHash)

		newBlocksInTransit := [][]byte{}
		for _, block := range blocksInTransit {
			if bytes.Compare(block, blockHash) != 0 {
				newBlocksInTransit = append(newBlocksInTransit, block)
			}
		}
		blocksInTransit = newBlocksInTransit
	}

	if payload.Type == inv.Transaction {
		txID := payload.Items[0]

		if mempool[hex.EncodeToString(txID)].ID == nil {
			sendGetData(payload.AddressFrom, "tx", txID)
		}
	}
}
