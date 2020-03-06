package node

/*
 * Copyright 2020 Information Wants To Be Free
 * Visit: https://github.com/iwtbf
 *
 * This project is licensed under the terms of the Apache 2.0 License.
 */

import (
	"github.com/alecthomas/kong"
)

var nodeCli struct {
	Relay relayNodeCmd `cmd optional help:"Start a blockweb relay."`
}

func Run() {
	context := kong.Parse(&nodeCli)

	err := context.Run()
	context.FatalIfErrorf(err)
}
