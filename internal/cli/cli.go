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
	"log"
	"os"
)

// TODO: Use a good github package
var (
	usage = `Usage:
	createwallet - Generates a new key-pair and saves it into the wallet file
	getbalance -address ADDRESS - Get balance of ADDRESS
	listaddresses - Lists all addresses from the wallet file
	reindex - Rebuilds the UnspentTransactionOutputSet
	printchain - Print all the blocks of the blockchain
	send -from FROM -to TO -amount AMOUNT - Send AMOUNT of coins from FROM address to TO`

	getBalanceCmd       = flag.NewFlagSet("getbalance", flag.ExitOnError)
	createBlockchainCmd = flag.NewFlagSet("createblockchain", flag.ExitOnError)
	createWalletCmd     = flag.NewFlagSet("createwallet", flag.ExitOnError)
	listAddressesCmd    = flag.NewFlagSet("listaddresses", flag.ExitOnError)
	reindexCmd          = flag.NewFlagSet("reindex", flag.ExitOnError)
	sendCmd             = flag.NewFlagSet("send", flag.ExitOnError)
	printChainCmd       = flag.NewFlagSet("printchain", flag.ExitOnError)

	getBalanceAddress = getBalanceCmd.String("address", "", "The address to get balance for")
	sendFrom          = sendCmd.String("from", "", "Source wallet address")
	sendTo            = sendCmd.String("to", "", "Destination wallet address")
	sendAmount        = sendCmd.Int("amount", 0, "Amount to send")
)

type CLI struct {
}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

// TODO: Use a good github package
func (cli *CLI) printUsage() {
	fmt.Println(usage)
}

func (cli *CLI) parseArguments() {
	switch os.Args[1] {
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createblockchain":
		err := createBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createwallet":
		err := createWalletCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "listaddresses":
		err := listAddressesCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "reindex":
		err := reindexCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.printUsage()
		os.Exit(1)
	}
}

func (cli *CLI) executeCommand(nodeID string) {
	if getBalanceCmd.Parsed() {
		if *getBalanceAddress == "" {
			getBalanceCmd.Usage()
			os.Exit(1)
		}

		cli.getBalance(*getBalanceAddress, nodeID)
	}

	if createBlockchainCmd.Parsed() {
		cli.createBlockchain(nodeID)
	}

	if createWalletCmd.Parsed() {
		cli.createWallet(nodeID)
	}

	if listAddressesCmd.Parsed() {
		cli.listAddresses(nodeID)
	}

	if printChainCmd.Parsed() {
		cli.printChain(nodeID)
	}

	if reindexCmd.Parsed() {
		cli.reindexUnspentTransactionOutputSet(nodeID)
	}

	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCmd.Usage()
			os.Exit(1)
		}

		cli.send(*sendFrom, *sendTo, *sendAmount, nodeID)
	}
}

func (cli *CLI) Run() {
	cli.validateArgs()

	nodeID := os.Getenv("NODE_ID")
	if nodeID == "" {
		fmt.Println("NODE_ID env var is not set!")
		os.Exit(1)
	}

	cli.parseArguments()
	cli.executeCommand(nodeID)
}

func NewCli() *CLI {
	return &CLI{}
}
