package blockchain

/*
 * Copyright 2020 Information Wants To Be Free
 * Visit: https://github.com/iwtbf
 *
 * This project is licensed under the terms of the Apache 2.0 License.
 */

type MerkleTree struct {
	RootNode *MerkleNode
}

func NewMerkleTree(leafes [][]byte) *MerkleTree {
	var merkleNodes []MerkleNode

	if len(leafes)%2 != 0 {
		leafes = append(leafes, leafes[len(leafes)-1])
	}

	for _, hashes := range leafes {
		merkleNode := NewMerkleNode(nil, nil, hashes)
		merkleNodes = append(merkleNodes, *merkleNode)
	}

	for i := 0; i < len(leafes)/2; i++ {
		var nextMerkleNodesLevel []MerkleNode

		for j := 0; j < len(merkleNodes); j += 2 {
			merkleNode := NewMerkleNode(&merkleNodes[j], &merkleNodes[j+1], nil)
			nextMerkleNodesLevel = append(nextMerkleNodesLevel, *merkleNode)
		}

		merkleNodes = nextMerkleNodesLevel
	}

	merkleTree := MerkleTree{RootNode: &merkleNodes[0]}

	return &merkleTree
}
