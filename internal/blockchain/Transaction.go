package blockchain

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
