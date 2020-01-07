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
	// TODO: CLI should connect without arguments
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
	defer cli.Blockchain.CloseDB()

	balance := 0
	unspentTransactionOutputs := cli.Blockchain.FindUnspentTransactionOutputs(address)

	for _, transactionOutput := range unspentTransactionOutputs {
		balance += transactionOutput.Value
	}

	fmt.Printf("Balance of '%s': %d\n", address, balance)
}

func (cli *CLI) send(from, to string, amount int) {
	defer cli.Blockchain.CloseDB()

	tx := blockchain.NewTransaction(from, to, amount, cli.Blockchain)
	cli.Blockchain.MineBlock([]*blockchain.Transaction{tx})

	fmt.Println("Success!")
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
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	getBalanceAddress := getBalanceCmd.String("address", "", "The address to get balance for")
	sendFrom := sendCmd.String("from", "", "Source wallet address")
	sendTo := sendCmd.String("to", "", "Destination wallet address")
	sendAmount := sendCmd.Int("amount", 0, "Amount to send")

	switch os.Args[1] {
	case "getbalance":
		// TODO: Error handling
		getBalanceCmd.Parse(os.Args[2:])
	case "printchain":
		// TODO: Error handling
		printChainCmd.Parse(os.Args[2:])
	case "send":
		// TODO: Error handling
		sendCmd.Parse(os.Args[2:])
	default:
		cli.printUsage()
		os.Exit(1)
	}

	if getBalanceCmd.Parsed() {
		if *getBalanceAddress == "" {
			getBalanceCmd.Usage()
			os.Exit(1)
		}

		cli.getBalance(*getBalanceAddress)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}

	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCmd.Usage()
			os.Exit(1)
		}

		cli.send(*sendFrom, *sendTo, *sendAmount)
	}
}
