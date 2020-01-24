package blockchain

/*
 * Copyright 2020 Information Wants To Be Free
 * Visit: https://github.com/iwtbf
 *
 * This project is licensed under the terms of the Apache 2.0 License.
 */

import (
	"bytes"
	"github.com/btcsuite/btcutil/base58"
)

type TransactionOutput struct {
	Value         int
	PublicKeyHash []byte
}

func (transactionOutput *TransactionOutput) Lock(address []byte) {
	publicKeyHash := base58.Decode(string(address[:]))
	publicKeyHash = publicKeyHash[1 : len(publicKeyHash)-4]
	transactionOutput.PublicKeyHash = publicKeyHash
}

func (transactionOutput *TransactionOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(transactionOutput.PublicKeyHash, pubKeyHash) == 0
}

func NewTransactionOutput(value int, address string) *TransactionOutput {
	transactionOutput := &TransactionOutput{Value: value, PublicKeyHash: nil}
	transactionOutput.Lock([]byte(address))

	return transactionOutput
}
