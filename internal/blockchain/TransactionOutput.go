package blockchain

/*
 * Copyright 2020 Information Wants To Be Free
 * Visit: https://github.com/iwtbf
 *
 * This project is licensed under the terms of the Apache 2.0 License.
 */

type TransactionOutput struct {
	Value           int
	ScriptPublicKey string
}

func (transactionOutput *TransactionOutput) CanUnlockUsing(unlockingData string) bool {
	return transactionOutput.ScriptPublicKey == unlockingData
}
