package keygen

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
	Gen genCmd `cmd help:"Generate a new ECDSA key pair."`
}

func Run() {
	context := kong.Parse(&dnsSeedCli)

	err := context.Run()
	context.FatalIfErrorf(err)
}
