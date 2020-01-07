package blockchain

import (
	"encoding/hex"
	"fmt"
	"os"
)

/*
 * Copyright 2020 Information Wants To Be Free
 * Visit: https://github.com/iwtbf
 *
 * This project is licensed under the terms of the Apache 2.0 License.
 */

var coinbaseTransactionValue = "I am a dreamer. Seriously, I'm living on another planet."

type Transaction struct {
	ID   []byte
	Vin  []TransactionInput
	Vout []TransactionOutput
}

func (transaction Transaction) setID() {
	// TODO: Generate transaction id
}

func (transaction Transaction) IsCoinbase() bool {
	return len(transaction.Vin) == 1 && transaction.Vin[0].ScriptSignature == coinbaseTransactionValue
}

func NewCoinbaseTransaction(to string) *Transaction {
	transactionInput := TransactionInput{[]byte{}, -1, coinbaseTransactionValue}
	transactionOutput := TransactionOutput{50000, to}
	transaction := Transaction{nil, []TransactionInput{transactionInput}, []TransactionOutput{transactionOutput}}
	transaction.setID()

	return &transaction
}

func NewTransaction(from, to string, amount int, blockchain *Blockchain) *Transaction {
	var transactionInputs []TransactionInput
	var transactionOutputs []TransactionOutput

	balance, unspentOutputs := blockchain.FindSpendableOutputs(from, amount)

	if balance < amount {
		// TODO: Maybe use log.Panic
		fmt.Printf("Not enough balance on account!")
		os.Exit(1)
	}

	for outputIndex, unspentOutput := range unspentOutputs {
		// TODO: Error handling
		transactionID, _ := hex.DecodeString(outputIndex)

		for _, outputTransaction := range unspentOutput {
			input := TransactionInput{transactionID, outputTransaction, from}
			transactionInputs = append(transactionInputs, input)
		}
	}

	transactionOutputs = append(transactionOutputs, TransactionOutput{amount, to})
	if balance > amount {
		transactionOutputs = append(transactionOutputs, TransactionOutput{balance - amount, from}) // a change
	}

	transaction := Transaction{nil, transactionInputs, transactionOutputs}
	transaction.setID()

	return &transaction
}
