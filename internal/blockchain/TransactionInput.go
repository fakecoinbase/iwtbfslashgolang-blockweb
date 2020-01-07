package blockchain

/*
 * Copyright 2020 Information Wants To Be Free
 * Visit: https://github.com/iwtbf
 *
 * This project is licensed under the terms of the Apache 2.0 License.
 */

type TransactionInput struct {
	TransactionID   []byte
	Vout            int
	ScriptSignature string
}

func (transactionInput *TransactionInput) CanUnlockUsing(unlockingData string) bool {
	return transactionInput.ScriptSignature == unlockingData
}
