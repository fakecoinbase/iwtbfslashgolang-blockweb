package dns_seed

/*
 * Copyright 2020 Information Wants To Be Free
 * Visit: https://github.com/iwtbf
 *
 * This project is licensed under the terms of the Apache 2.0 License.
 */

import (
	"github.com/alecthomas/kong"
)

type startDnsSeedCmd struct {
	Port  int16  `help:"The servers listening Port." default:"10000"`
	Level string `help:"One of github.com/ipfs/go-log#LogLevel." default:"INFO"`
}

var dnsSeedCli struct {
	Start startDnsSeedCmd `cmd help:"Start a DNS seed server."`
}

func Run() {
	context := kong.Parse(&dnsSeedCli)

	err := context.Run()
	context.FatalIfErrorf(err)
}
