package network

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
	start -miner ADDRESS - Start a node named according to environment variable NODE_ID -miner enables mining`

	startCmd = flag.NewFlagSet("start", flag.ExitOnError)

	startMiner = startCmd.String("miner", "", "Enable mining mode and send reward to ADDRESS")
)

type NodeCLI struct {
}

func (nodeCLI *NodeCLI) validateArgs() {
	if len(os.Args) < 2 {
		nodeCLI.printUsage()
		os.Exit(1)
	}
}

// TODO: Use a good github package
func (nodeCLI *NodeCLI) printUsage() {
	fmt.Println(usage)
}

func (nodeCLI *NodeCLI) parseArguments() {
	switch os.Args[1] {
	case "start":
		err := startCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		nodeCLI.printUsage()
		os.Exit(1)
	}
}

func (nodeCLI *NodeCLI) executeCommand(nodeID string) {
	if startCmd.Parsed() {
		nodeCLI.startNode(nodeID, *startMiner)
	}
}

func (nodeCLI *NodeCLI) Run() {
	nodeCLI.validateArgs()

	nodeID := os.Getenv("NODE_ID")
	if nodeID == "" {
		fmt.Println("NODE_ID env var is not set!")
		os.Exit(1)
	}

	nodeCLI.parseArguments()
	nodeCLI.executeCommand(nodeID)
}

func NewNodeCLI() *NodeCLI {
	return &NodeCLI{}
}
