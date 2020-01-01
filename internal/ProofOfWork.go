package internal

/*
 * Copyright 2019 Information Wants To Be Free
 * Visit: https://github.com/iwtbf
 *
 * This project is licensed under the terms of the Apache 2.0 License.
 */

import (
	"../pkg"
	"bytes"
	"crypto/sha256"
	"math"
	"math/big"
)

const targetBits = 24
const maxNonce = math.MaxInt64

// TODO: Followed the tutorial, but proof of stake might be more applicable
type ProofOfWork struct {
	block  *Block
	target *big.Int
}

func (proofOfWork *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			proofOfWork.block.PreviousHash,
			proofOfWork.block.Data,
			pkg.IntToHex(proofOfWork.block.Timestamp),
			pkg.IntToHex(int64(targetBits)),
			pkg.IntToHex(int64(nonce)),
		},
		[]byte{},
	)

	return data
}

func (proofOfWork *ProofOfWork) Run() (int, []byte) {
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

func (proofOfWork *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := proofOfWork.prepareData(proofOfWork.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	isValid := hashInt.Cmp(proofOfWork.target) == -1

	return isValid
}

func NewProofOfWork(block *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))

	proofOfWork := &ProofOfWork{block, target}

	return proofOfWork
}
