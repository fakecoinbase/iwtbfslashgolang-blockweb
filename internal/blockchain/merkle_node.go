package blockchain

/*
 * Copyright 2020 Information Wants To Be Free
 * Visit: https://github.com/iwtbf
 *
 * This project is licensed under the terms of the Apache 2.0 License.
 */

import "crypto/sha256"

type merkleNode struct {
	Left  *merkleNode
	Right *merkleNode
	Hash  []byte
}

func newMerkleNode(left, right *merkleNode, hash []byte) *merkleNode {
	merkleNode := merkleNode{}

	if left == nil && right == nil {
		hash := sha256.Sum256(hash)
		merkleNode.Hash = hash[:]
	} else {
		prevHashes := append(left.Hash, right.Hash...)
		hash := sha256.Sum256(prevHashes)
		merkleNode.Hash = hash[:]
	}

	merkleNode.Left = left
	merkleNode.Right = right

	return &merkleNode
}
