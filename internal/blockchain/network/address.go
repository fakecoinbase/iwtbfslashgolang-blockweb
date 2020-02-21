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
)

type address struct {
	AddressList []string
}

func sendAddress(newAddress string) {
	nodes := address{AddressList: knownNodes}
	nodes.AddressList = append(nodes.AddressList, nodeAddress)
	payload := gobEncode(nodes)
	request := append(commandToBytes("addr"), payload...)

	sendData(newAddress, request)
}

func handleAddress(request []byte) {
	var buffer bytes.Buffer
	var payload address

	buffer.Write(request[commandLength:])
	decoder := gob.NewDecoder(&buffer)
	// TODO: Error handling
	decoder.Decode(&payload)

	knownNodes = append(knownNodes, payload.AddressList...)
	fmt.Printf("There are %d known nodes now!\n", len(knownNodes))
	requestBlocks()
}
