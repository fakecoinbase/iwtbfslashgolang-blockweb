package cli

/*
 * Copyright 2020 Information Wants To Be Free
 * Visit: https://github.com/iwtbf
 *
 * This project is licensed under the terms of the Apache 2.0 License.
 */

import (
	"flag"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"github.com/iwtbf/golang-blockweb/internal/blockchain"
	"os"
	"strconv"
)

type CLI struct {
	blockchain *blockchain.Blockchain
}

func (cli *CLI) printChain() {
	defer cli.blockchain.CloseDB()

	blockchainIterator := cli.blockchain.Iterator()

	for {
		block := blockchainIterator.Next()

		fmt.Printf("Previous hash: %x\n", block.PreviousHash)
		fmt.Printf("Data: %s\n", block.Transactions)
		fmt.Printf("hash: %x\n", block.Hash)

		proofOfWork := blockchain.NewProofOfWork(block)

		fmt.Printf("Valid?: %s\n", strconv.FormatBool(proofOfWork.Validate()))
		fmt.Println()

		if len(block.PreviousHash) == 0 {
			break
		}
	}
}

func (cli *CLI) getBalance(address string) {
	defer cli.blockchain.CloseDB()

	balance := 0
	publicKeyHash := base58.Decode(address)
	publicKeyHash = publicKeyHash[1 : len(publicKeyHash)-4]
	unspentTransactionOutputs := cli.blockchain.FindUnspentOutputs(publicKeyHash)

	for _, unspentTransactionOutput := range unspentTransactionOutputs {
		balance += unspentTransactionOutput.Value
	}

	fmt.Printf("Balance of '%s': %d\n", address, balance)
}

func (cli *CLI) send(from, to string, amount int) {
	defer cli.blockchain.CloseDB()

	tx := blockchain.NewTransaction([]byte(from), []byte(to), amount, cli.blockchain)
	cli.blockchain.MineBlock([]*blockchain.Transaction{tx})

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
	getBalanceAddress := getBalanceCmd.String("address", "", "The address to get balance for")

	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	sendFrom := sendCmd.String("from", "", "Source wallet address")
	sendTo := sendCmd.String("to", "", "Destination wallet address")
	sendAmount := sendCmd.Int("amount", 0, "Amount to send")

	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

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

func NewCli() *CLI {
	return &CLI{blockchain.NewBlockchain("", "")}
}
