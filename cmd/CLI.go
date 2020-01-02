package cmd

/*
 * Copyright 2020 Information Wants To Be Free
 * Visit: https://github.com/iwtbf
 *
 * This project is licensed under the terms of the Apache 2.0 License.
 */

import (
	"../internal"
	"flag"
	"fmt"
	"os"
	"strconv"
)

type CLI struct {
	Blockchain *internal.Blockchain
}

func (cli *CLI) addBlock(data string) {
	cli.Blockchain.AddBlock(data)
	fmt.Println("Success!")
}

func (cli *CLI) printChain() {
	blockchainIterator := cli.Blockchain.Iterator()

	for {
		block := blockchainIterator.Next()

		fmt.Printf("Prev. hash: %x\n", block.PreviousHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)

		proofOfWork := internal.NewProofOfWork(block)

		fmt.Printf("PoW: %s\n", strconv.FormatBool(proofOfWork.Validate()))
		fmt.Println()

		if len(block.PreviousHash) == 0 {
			break
		}
	}
}

func (cli *CLI) Run() {
	cli.validateArgs()

	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	addBlockData := addBlockCmd.String("data", "", "Block data")

	switch os.Args[1] {
	case "addblock":
		// TODO: Error handling
		err := addBlockCmd.Parse(os.Args[2:])
	case "printchain":
		// TODO: Error handling
		err := printChainCmd.Parse(os.Args[2:])
	default:
		cli.printUsage()
		os.Exit(1)
	}

	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			addBlockCmd.Usage()
			os.Exit(1)
		}

		cli.addBlock(*addBlockData)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}
}
