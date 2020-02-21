package command

/*
 * Copyright 2020 Information Wants To Be Free
 * Visit: https://github.com/iwtbf
 *
 * This project is licensed under the terms of the Apache 2.0 License.
 */

type Command string

const (
	Address     Command = "address"
	Block       Command = "block"
	GetBlocks   Command = "getblocks"
	GetData     Command = "getdata"
	Inventory   Command = "inventory"
	Transaction Command = "transaction"
	Version     Command = "version"
)
