package main

/*
 * Copyright 2020 Information Wants To Be Free
 * Visit: https://github.com/iwtbf
 *
 * This project is licensed under the terms of the Apache 2.0 License.
 */

import (
	"./cmd"
	"./internal"
)

func main() {
	blockchain := internal.NewBlockchain("")
	defer blockchain.DB.Close()

	cli := cmd.CLI{blockchain}
	cli.Run()
}
