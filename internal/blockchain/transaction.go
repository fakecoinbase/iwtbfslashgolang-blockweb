package blockchain

/*
 * Copyright 2020 Information Wants To Be Free
 * Visit: https://github.com/iwtbf
 *
 * This project is licensed under the terms of the Apache 2.0 License.
 */

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
)

var coinbaseTransactionValue = "I am a dreamer. Seriously, I'm living on another planet."

type Transaction struct {
	ID                 []byte
	TransactionInputs  []TransactionInput
	TransactionOutputs []TransactionOutput
}

// TODO: Use proto buffers
func (transaction *Transaction) Serialize() []byte {
	var encoded bytes.Buffer

	encoder := gob.NewEncoder(&encoded)
	// TODO: Error handling
	encoder.Encode(transaction)

	return encoded.Bytes()
}

func (transaction *Transaction) hash() []byte {
	var hash [32]byte

	transactionCopy := *transaction
	transactionCopy.ID = []byte{}

	hash = sha256.Sum256(transactionCopy.Serialize())

	return hash[:]
}

func (transaction Transaction) IsCoinbase() bool {
	return len(transaction.TransactionInputs) == 1 && len(transaction.TransactionInputs[0].TransactionID) == 0 && transaction.TransactionInputs[0].transactionOutputID == -1
}

func (transaction *Transaction) Sign(privateKey ecdsa.PrivateKey, previousTransactions map[string]Transaction) {
	if transaction.IsCoinbase() {
		return
	}

	transactionCopy := transaction.TrimmedCopy()

	for inputIterator, transactionInput := range transactionCopy.TransactionInputs {
		previousTransaction := previousTransactions[hex.EncodeToString(transactionInput.TransactionID)]
		transactionCopy.TransactionInputs[inputIterator].Signature = nil
		transactionCopy.TransactionInputs[inputIterator].PublicKey = previousTransaction.TransactionOutputs[transactionInput.transactionOutputID].PublicKeyHash
		transactionCopy.ID = transactionCopy.hash()
		transactionCopy.TransactionInputs[inputIterator].PublicKey = nil

		// TODO: Error handling
		r, s, _ := ecdsa.Sign(rand.Reader, &privateKey, transactionCopy.ID)
		signature := append(r.Bytes(), s.Bytes()...)

		transaction.TransactionInputs[inputIterator].Signature = signature
	}
}

func (transaction *Transaction) Verify(previousTransactions map[string]Transaction) bool {
	transactionCopy := transaction.TrimmedCopy()
	curve := elliptic.P256()

	for inputIterator, transactionInput := range transaction.TransactionInputs {
		previousTransaction := previousTransactions[hex.EncodeToString(transactionInput.TransactionID)]
		transactionCopy.TransactionInputs[inputIterator].Signature = nil
		transactionCopy.TransactionInputs[inputIterator].PublicKey = previousTransaction.TransactionOutputs[transactionInput.transactionOutputID].PublicKeyHash
		transactionCopy.ID = transactionCopy.hash()
		transactionCopy.TransactionInputs[inputIterator].PublicKey = nil

		r := big.Int{}
		s := big.Int{}
		sigLen := len(transactionInput.Signature)
		r.SetBytes(transactionInput.Signature[:(sigLen / 2)])
		s.SetBytes(transactionInput.Signature[(sigLen / 2):])

		x := big.Int{}
		y := big.Int{}
		keyLen := len(transactionInput.PublicKey)
		x.SetBytes(transactionInput.PublicKey[:(keyLen / 2)])
		y.SetBytes(transactionInput.PublicKey[(keyLen / 2):])

		rawPubKey := ecdsa.PublicKey{curve, &x, &y}
		if ecdsa.Verify(&rawPubKey, transactionCopy.ID, &r, &s) == false {
			return false
		}
	}

	return true
}

func NewCoinbaseTransaction(to, data string) *Transaction {
	transactionInput := TransactionInput{TransactionID: []byte{}, transactionOutputID: -1, Signature: nil, PublicKey: []byte(coinbaseTransactionValue)}
	transactionOutput := NewTransactionOutput(50000, to)
	transaction := Transaction{nil, []TransactionInput{transactionInput}, []TransactionOutput{*transactionOutput}}
	transaction.hash()

	return &transaction
}

func NewTransaction(from, to []byte, amount int, blockchain *Blockchain) *Transaction {
	var transactionInputs []TransactionInput
	var transactionOutputs []TransactionOutput

	if amount <= 0 {
		// TODO: This is not a transaction
	}

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
			input := TransactionInput{transactionID, outputTransaction, from, nil}
			transactionInputs = append(transactionInputs, input)
		}
	}

	transactionOutputs = append(transactionOutputs, TransactionOutput{amount, to})
	if balance > amount {
		transactionOutputs = append(transactionOutputs, TransactionOutput{balance - amount, from}) // a change
	}

	transaction := Transaction{nil, transactionInputs, transactionOutputs}
	transaction.hash()

	return &transaction
}

func (transaction *Transaction) TrimmedCopy() Transaction {
	var inputs []TransactionInput
	var outputs []TransactionOutput

	for _, transactionInput := range transaction.TransactionInputs {
		inputs = append(inputs, TransactionInput{TransactionID: transactionInput.TransactionID, transactionOutputID: transactionInput.transactionOutputID, Signature: nil, PublicKey: nil})
	}

	for _, transactionOutput := range transaction.TransactionOutputs {
		outputs = append(outputs, TransactionOutput{Value: transactionOutput.Value, PublicKeyHash: transactionOutput.PublicKeyHash})
	}

	transactionCopy := Transaction{transaction.ID, inputs, outputs}

	return transactionCopy
}
