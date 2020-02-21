package blockchain

/*
 * Copyright 2020 Information Wants To Be Free
 * Visit: https://github.com/iwtbf
 *
 * This project is licensed under the terms of the Apache 2.0 License.
 */

import (
	"bytes"
	"crypto/sha256"
	"github.com/iwtbf/golang-blockweb/pkg/Int2Hex"
	"math"
	"math/big"
)

const targetBits = 2
const maxNonce = math.MaxInt64

// TODO: Followed the tutorial, but proof of stake might be more applicable
type proofOfWork struct {
	block  *Block
	target *big.Int
}

func (proofOfWork *proofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			proofOfWork.block.PreviousHash,
			proofOfWork.block.HashTransactions(),
			Int2Hex.Convert(proofOfWork.block.Timestamp),
			Int2Hex.Convert(int64(targetBits)),
			Int2Hex.Convert(int64(nonce)),
		},
		[]byte{},
	)

	return data
}

func (proofOfWork *proofOfWork) run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	for nonce < maxNonce {
		data := proofOfWork.prepareData(nonce)
		hash = sha256.Sum256(data)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(proofOfWork.target) == -1 {
			break
		} else {
			nonce++
		}
	}

	return nonce, hash[:]
}

func (proofOfWork *proofOfWork) Validate() bool {
	var hashInt big.Int

	data := proofOfWork.prepareData(proofOfWork.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	isValid := hashInt.Cmp(proofOfWork.target) == -1

	return isValid
}

func NewProofOfWork(block *Block) *proofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))

	proofOfWork := &proofOfWork{block: block, target: target}

	return proofOfWork
}
