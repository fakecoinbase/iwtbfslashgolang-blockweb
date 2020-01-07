package cmd

/*
 * Copyright 2020 Information Wants To Be Free
 * Visit: https://github.com/iwtbf
 *
 * This project is licensed under the terms of the Apache 2.0 License.
 */

import (
	"../internal/blockchain"
	"flag"
	"fmt"
	"os"
	"strconv"
)

type CLI struct {
	Blockchain *blockchain.Blockchain
}

func (cli *CLI) printChain() {
	blockchainIterator := cli.Blockchain.Iterator()

	for {
		block := blockchainIterator.Next()

		fmt.Printf("Previous hash: %x\n", block.PreviousHash)
		fmt.Printf("Data: %s\n", block.Transactions)
		fmt.Printf("Hash: %x\n", block.Hash)

		proofOfWork := blockchain.NewProofOfWork(block)

		fmt.Printf("Valid?: %s\n", strconv.FormatBool(proofOfWork.Validate()))
		fmt.Println()

		if len(block.PreviousHash) == 0 {
			break
		}
	}
}

func (cli *CLI) getBalance(address string) {
	defer cli.Blockchain.DB.Close()

	balance := 0
	unspentTransactionOutputs := cli.Blockchain.FindUnspentTransactionOutputs(address)

	for _, transactionOutput := range unspentTransactionOutputs {
		balance += transactionOutput.Value
	}

	fmt.Printf("Balance of '%s': %d\n", address, balance)
}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

func (cli *CLI) printUsage() {
	fmt.Printf("Consider learning usage - noob!\n")
}

func (cli *CLI) Run() {
	cli.validateArgs()

	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	getBalanceData := getBalanceCmd.String("address", "", "Coinbase Address")

	switch os.Args[1] {
	case "getbalance":
		// TODO: Error handling
		getBalanceCmd.Parse(os.Args[2:])
	case "printchain":
		// TODO: Error handling
		printChainCmd.Parse(os.Args[2:])
	default:
		cli.printUsage()
		os.Exit(1)
	}

	if getBalanceCmd.Parsed() {
		if *getBalanceData == "" {
			getBalanceCmd.Usage()
			os.Exit(1)
		}

		cli.getBalance(*getBalanceData)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}
}
