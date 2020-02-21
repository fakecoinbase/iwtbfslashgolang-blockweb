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
	"log"
	"math/big"
)

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

// TODO: Use proto buffers
func DeserializeTransaction(data []byte) *Transaction {
	var transaction Transaction

	decoder := gob.NewDecoder(bytes.NewReader(data))
	// TODO: Error handling
	decoder.Decode(&transaction)

	return &transaction
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

		rawPubKey := ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}
		if ecdsa.Verify(&rawPubKey, transactionCopy.ID, &r, &s) == false {
			return false
		}
	}

	return true
}

func NewCoinbaseTransaction(to []byte) *Transaction {
	transactionInput := TransactionInput{TransactionID: []byte{}, transactionOutputID: -1, Signature: nil, PublicKey: []byte(genesisCoinbaseData)}
	transactionOutput := NewTransactionOutput(50000, to)
	transaction := Transaction{ID: nil, TransactionInputs: []TransactionInput{transactionInput}, TransactionOutputs: []TransactionOutput{*transactionOutput}}
	transaction.ID = transaction.hash()

	return &transaction
}

func NewTransaction(wallet *Wallet, to []byte, amount int, unspentTransactionOutputSet UnspentTransactionOutputSet) *Transaction {
	var transactionInputs []TransactionInput
	var transactionOutputs []TransactionOutput

	publicKeyHash := HashPublicKey(wallet.PublicKey)
	accumulated, spendableOutputs := unspentTransactionOutputSet.FindSpendableOutputs(publicKeyHash, amount)

	if accumulated < amount {
		// TODO: Not enough balance
		log.Panic("Not enough funds")
	}

	for outputIDIterator, spendableOutputIDs := range spendableOutputs {
		// TODO: Error handling
		transactionID, _ := hex.DecodeString(outputIDIterator)

		for _, spendableOutputID := range spendableOutputIDs {
			input := TransactionInput{TransactionID: transactionID, transactionOutputID: spendableOutputID, Signature: nil, PublicKey: wallet.PublicKey}
			transactionInputs = append(transactionInputs, input)
		}
	}

	from := string(wallet.GetAddress()[:])
	transactionOutputs = append(transactionOutputs, *NewTransactionOutput(amount, to))
	if accumulated > amount {
		transactionOutputs = append(transactionOutputs, *NewTransactionOutput(accumulated-amount, []byte(from))) // a change
	}

	transaction := Transaction{ID: nil, TransactionInputs: transactionInputs, TransactionOutputs: transactionOutputs}
	transaction.ID = transaction.hash()
	unspentTransactionOutputSet.Blockchain.SignTransaction(&transaction, wallet.PrivateKey)

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

	transactionCopy := Transaction{ID: transaction.ID, TransactionInputs: inputs, TransactionOutputs: outputs}

	return transactionCopy
}
