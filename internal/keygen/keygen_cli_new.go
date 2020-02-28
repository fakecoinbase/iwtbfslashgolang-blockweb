package keygen

/*
 * Copyright 2020 Information Wants To Be Free
 * Visit: https://github.com/iwtbf
 *
 * This project is licensed under the terms of the Apache 2.0 License.
 */

type newCmd struct {
	Out   string `arg help:"The path to the output file."`
	Print bool   `flag help:"Print the *private key* to System.Out as well."`
}

func (cmd *newCmd) Run() error {
	newKeyPair(cmd.Out, cmd.Print)

	return nil
}
