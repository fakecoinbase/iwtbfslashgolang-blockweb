package blockchain

/*
 * Copyright 2020 Information Wants To Be Free
 * Visit: https://github.com/iwtbf
 *
 * This project is licensed under the terms of the Apache 2.0 License.
 */

type merkleTree struct {
	RootNode *merkleNode
}

func newMerkleTree(leafes [][]byte) *merkleTree {
	var merkleNodes []merkleNode

	if len(leafes)%2 != 0 {
		leafes = append(leafes, leafes[len(leafes)-1])
	}

	for _, hashes := range leafes {
		merkleNode := newMerkleNode(nil, nil, hashes)
		merkleNodes = append(merkleNodes, *merkleNode)
	}

	for i := 0; i < len(leafes)/2; i++ {
		var nextMerkleNodesLevel []merkleNode

		for j := 0; j < len(merkleNodes); j += 2 {
			merkleNode := newMerkleNode(&merkleNodes[j], &merkleNodes[j+1], nil)
			nextMerkleNodesLevel = append(nextMerkleNodesLevel, *merkleNode)
		}

		merkleNodes = nextMerkleNodesLevel
	}

	merkleTree := merkleTree{RootNode: &merkleNodes[0]}

	return &merkleTree
}
