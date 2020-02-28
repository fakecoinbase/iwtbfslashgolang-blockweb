package keygen

/*
 * Copyright 2020 Information Wants To Be Free
 * Visit: https://github.com/iwtbf
 *
 * This project is licensed under the terms of the Apache 2.0 License.
 */

import (
	"github.com/alecthomas/kong"
	"os"
)

const (
	disclaimer = `DISCLAIMER:
This program * can * be used to generate key pairs for 'golang-blockweb'.
It is open source software, but there is no liability of the contributors or guarantee for the security of the generated keys. Please report any problems or vulnerabilities at https://github.com/iwtbf/golang-blockweb/issues/new.
You could use any other key tool with elliptic P256 curve cryptography (ECDSA) + x509 + PEM. Use the command 'read' to verify the generated key.
By continuing, you confirm that you have understood the above points (y / N):`
	retry = `Sorry, I didn't understand this. Please try again (y / N):`
)

var dnsSeedCli struct {
	Gen  newCmd  `cmd optional help:"Generate a new ECDSA key pair."`
	Read readCmd `cmd optional help:"Validate an existing ECDSA private key."`
}

func Run() {
	context := kong.Parse(&dnsSeedCli)

	if !acceptDisclaimer(disclaimer, retry) {
		os.Exit(0)
	}

	err := context.Run()
	context.FatalIfErrorf(err)
}
