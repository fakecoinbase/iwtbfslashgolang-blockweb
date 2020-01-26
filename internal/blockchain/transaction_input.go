package blockchain

/*
 * Copyright 2020 Information Wants To Be Free
 * Visit: https://github.com/iwtbf
 *
 * This project is licensed under the terms of the Apache 2.0 License.
 */

import "bytes"

type TransactionInput struct {
	TransactionID       []byte
	transactionOutputID int
	Signature           []byte
	PublicKey           []byte
}

func (transactionInput *TransactionInput) UsesKey(pubKeyHash []byte) bool {
	lockingHash := HashPublicKey(transactionInput.PublicKey)

	return bytes.Compare(lockingHash, pubKeyHash) == 0
}
