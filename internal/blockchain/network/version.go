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
	"github.com/iwtbf/golang-blockweb/internal/blockchain"
	"github.com/iwtbf/golang-blockweb/internal/blockchain/network/command"
)

type version struct {
	Version     int
	BestHeight  int
	AddressFrom string
}

func sendVersion(address string, chain *blockchain.Blockchain) {
	payload := gobEncode(version{Version: nodeVersion, BestHeight: chain.GetBestHeight(), AddressFrom: nodeAddress})
	request := append(commandToBytes(command.Version), payload...)

	sendData(address, request)
}

func handleVersion(request []byte, chain *blockchain.Blockchain) {
	var buffer bytes.Buffer
	var payload version

	buffer.Write(request[commandLength:])
	decoder := gob.NewDecoder(&buffer)
	// TODO: Error handling
	decoder.Decode(&payload)

	myBestHeight := chain.GetBestHeight()
	foreignerBestHeight := payload.BestHeight

	if myBestHeight < foreignerBestHeight {
		sendGetBlocks(payload.AddressFrom)
	} else if myBestHeight > foreignerBestHeight {
		sendVersion(payload.AddressFrom, chain)
	}

	// sendAddress(payload.AddressFrom)
	if !isKnownNode(payload.AddressFrom) {
		knownNodes = append(knownNodes, payload.AddressFrom)
	}
}
