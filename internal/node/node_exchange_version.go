package node

/*
 * Copyright 2020 Information Wants To Be Free
 * Visit: https://github.com/iwtbf
 *
 * This project is licensed under the terms of the Apache 2.0 License.
 */

import "context"

// TODO: Really exchange (blockchain-)version
func (node *node) ExchangeVersion(context.Context, *Version) (*Version, error) {
	if node == nil {
		panic("node#ExchangeVersion called on nil reference!")
	}

	return &Version{
		Version:     0,
		BestHeight:  0,
		AddressFrom: "",
	}, nil
}
