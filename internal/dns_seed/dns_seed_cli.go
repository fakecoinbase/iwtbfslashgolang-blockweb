package dns_seed

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
	start -port PORT - Start a DNS seed listening on PORT (default 10000)`

	startCmd = flag.NewFlagSet("start", flag.ExitOnError)

	port = startCmd.Int("port", 10000, "The server port")
)

type DNSSeedCLI struct {
}

func (dnsSeedCLI *DNSSeedCLI) validateArgs() {
	if len(os.Args) < 2 {
		dnsSeedCLI.printUsage()
		os.Exit(1)
	}
}

// TODO: Use a good github package
func (dnsSeedCLI *DNSSeedCLI) printUsage() {
	fmt.Println(usage)
}

func (dnsSeedCLI *DNSSeedCLI) parseArguments() {
	switch os.Args[1] {
	case "start":
		err := startCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		dnsSeedCLI.printUsage()
		os.Exit(1)
	}
}

func (dnsSeedCLI *DNSSeedCLI) executeCommand() {
	if startCmd.Parsed() {
		dnsSeedCLI.startDNSSeed(*port)
	}
}

func (dnsSeedCLI *DNSSeedCLI) Run() {
	dnsSeedCLI.validateArgs()

	dnsSeedCLI.parseArguments()
	dnsSeedCLI.executeCommand()
}

func NewDNSSeedCLI() *DNSSeedCLI {
	return &DNSSeedCLI{}
}
