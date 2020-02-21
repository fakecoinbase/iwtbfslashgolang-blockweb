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
	"io"
	"io/ioutil"
	"net"
)

const protocol = "tcp"
const nodeVersion = 1
const commandLength = 12

var nodeAddress string
var miningAddress string
var knownNodes = []string{"localhost:3000"}
var blocksInTransit = [][]byte{}
var mempool = make(map[string]blockchain.Transaction)

func handleConnection(conn net.Conn, chain *blockchain.Blockchain) {
	// TODO: Error handling
	request, _ := ioutil.ReadAll(conn)

	command := bytesToCommand(request[:commandLength])
	fmt.Printf("Received %s command\n", command)

	switch command {
	case "addr":
		handleAddress(request)
	case "block":
		handleBlock(request, chain)
	case "inv":
		handleInventory(request, chain)
	case "getblocks":
		handleGetBlocks(request, chain)
	case "getdata":
		handleGetData(request, chain)
	case "tx":
		handleTransaction(request, chain)
	case "version":
		handleVersion(request, chain)
	default:
		fmt.Println("Unknown command!")
	}

	conn.Close()
}

func sendData(address string, data []byte) {
	conn, err := net.Dial(protocol, address)
	if err != nil {
		fmt.Printf("%s is not available\n", address)
		var updatedNodes []string

		for _, node := range knownNodes {
			if node != address {
				updatedNodes = append(updatedNodes, node)
			}
		}

		knownNodes = updatedNodes

		return
	}
	defer conn.Close()

	// TODO: Error handling
	io.Copy(conn, bytes.NewReader(data))
}

func commandToBytes(command string) []byte {
	var bytes [commandLength]byte

	for iterator, char := range command {
		bytes[iterator] = byte(char)
	}

	return bytes[:]
}

func bytesToCommand(bytes []byte) string {
	var command []byte

	for _, b := range bytes {
		if b != 0x0 {
			command = append(command, b)
		}
	}

	return string(command[:])
}

func gobEncode(data interface{}) []byte {
	var buff bytes.Buffer

	encoder := gob.NewEncoder(&buff)
	// TODO: Error handling
	encoder.Encode(data)

	return buff.Bytes()
}

func bootNode(nodeID, minerAddress string) {
	// TODO: Use DNS seed
	nodeAddress = fmt.Sprintf("localhost:%s", nodeID)
	miningAddress = minerAddress

	// TODO: Error handling
	listener, _ := net.Listen(protocol, nodeAddress)
	defer listener.Close()

	chain := blockchain.NewBlockchain(nodeID)

	if nodeAddress != knownNodes[0] {
		sendVersion(knownNodes[0], chain)
	}

	for {
		// TODO: Error handling
		conn, _ := listener.Accept()
		go handleConnection(conn, chain)
	}
}

func isKnownNode(address string) bool {
	for _, knownNode := range knownNodes {
		if knownNode == address {
			return true
		}
	}

	return false
}
