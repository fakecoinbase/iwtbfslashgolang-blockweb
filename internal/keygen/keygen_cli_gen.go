package keygen

/*
 * Copyright 2020 Information Wants To Be Free
 * Visit: https://github.com/iwtbf
 *
 * This project is licensed under the terms of the Apache 2.0 License.
 */

type genCmd struct {
}

func (cmd *genCmd) Run() error {
	newKeyPair()

	return nil
}
