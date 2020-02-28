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

var dnsSeedCli struct {
	Start startDnsSeedCmd `cmd optional help:"Start a DNS seed server."`
}

func Run() {
	context := kong.Parse(&dnsSeedCli)

	err := context.Run()
	context.FatalIfErrorf(err)
}
