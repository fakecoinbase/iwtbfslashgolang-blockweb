package blockchain

/*
 * Copyright 2020 Information Wants To Be Free
 * Visit: https://github.com/iwtbf
 *
 * This project is licensed under the terms of the Apache 2.0 License.
 */

import (
	"bytes"
	"encoding/gob"
)

type TransactionOutputSet struct {
	TransactionOutputs []TransactionOutput
}

// TODO: Maybe use proto buffs
func (transactionOutputSet TransactionOutputSet) Serialize() []byte {
	var buff bytes.Buffer

	encoder := gob.NewEncoder(&buff)
	// TODO: Error handling
	encoder.Encode(transactionOutputSet)

	return buff.Bytes()
}

func NewTransactionOutputSet(data []byte) TransactionOutputSet {
	var transactionOutputSet TransactionOutputSet

	decoder := gob.NewDecoder(bytes.NewReader(data))
	// TODO: Error handling
	decoder.Decode(&transactionOutputSet)

	return transactionOutputSet
}
