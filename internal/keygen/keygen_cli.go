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
	GenKey  genKeyCmd  `cmd optional help:"Generate a new ECDSA key pair."`
	ReadKey readKeyCmd `cmd optional help:"Validate an existing ECDSA private key."`
	GenCert genCertCmd `cmd optional help:"Generate a server certificate from ECDSA key file."`
}

func Run() {
	context := kong.Parse(&dnsSeedCli)

	err := context.Run()
	context.FatalIfErrorf(err)
}
