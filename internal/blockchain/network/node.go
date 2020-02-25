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
	cmd "github.com/iwtbf/golang-blockweb/internal/blockchain/network/command"
	"io"
	"io/ioutil"
	"net"
)

const (
	protocol      = "tcp"
	nodeVersion   = 1
	commandLength = 12
)

var (
	nodeAddress     string
	miningAddress   string
	knownNodes      = []string{"localhost:3000"}
	blocksInTransit = [][]byte{}
	mempool         = make(map[string]blockchain.Transaction)
)

func bootNode(nodeID, minerAddress string) {
	fmt.Printf("Starting node %s\n", nodeID)

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

	// TODO: Dispatch to background
	for {
		// TODO: Error handling
		conn, _ := listener.Accept()
		go handleConnection(conn, chain)
	}

	fmt.Printf("Node %s ready to accept connections\n", nodeID)
}

func handleConnection(conn net.Conn, chain *blockchain.Blockchain) {
	// TODO: Error handling
	request, _ := ioutil.ReadAll(conn)

	command := bytesToCommand(request[:commandLength])
	fmt.Printf("Received %s command\n", command)

	switch command {
	case cmd.Address:
		handleAddress(request)
	case cmd.Block:
		handleBlock(request, chain)
	case cmd.Inventory:
		handleInventory(request, chain)
	case cmd.GetBlocks:
		handleGetBlocks(request, chain)
	case cmd.GetData:
		handleGetData(request, chain)
	case cmd.Transaction:
		handleTransaction(request, chain)
	case cmd.Version:
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

		// TODO: Check if last node was removed

		return
	}

	defer conn.Close()

	// TODO: Error handling
	io.Copy(conn, bytes.NewReader(data))
}

func commandToBytes(command cmd.Command) []byte {
	var bytes [commandLength]byte

	for iterator, char := range command {
		bytes[iterator] = byte(char)
	}

	return bytes[:]
}

func bytesToCommand(bytes []byte) cmd.Command {
	var command []byte

	for _, b := range bytes {
		if b != 0x0 {
			command = append(command, b)
		}
	}

	return cmd.Command(string(command[:]))
}

func gobEncode(data interface{}) []byte {
	var buffer bytes.Buffer

	encoder := gob.NewEncoder(&buffer)
	// TODO: Error handling
	encoder.Encode(data)

	return buffer.Bytes()
}

func isKnownNode(address string) bool {
	for _, knownNode := range knownNodes {
		if knownNode == address {
			return true
		}
	}

	return false
}
