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

func newMerkleTree(data [][]byte) *merkleTree {
	var merkleNodes []merkleNode

	if len(data)%2 != 0 {
		data = append(data, data[len(data)-1])
	}

	for _, datum := range data {
		merkleNode := newMerkleNode(nil, nil, datum)
		merkleNodes = append(merkleNodes, *merkleNode)
	}

TreeNode:
	for {
		var nextMerkleNodesLevel []merkleNode

		for j := 0; j < len(merkleNodes); j += 2 {
			if len(merkleNodes) == 1 {
				break TreeNode
			}
			node := newMerkleNode(&merkleNodes[j], &merkleNodes[j+1], nil)
			nextMerkleNodesLevel = append(nextMerkleNodesLevel, *node)
		}

		merkleNodes = nextMerkleNodesLevel
	}

	merkleTree := merkleTree{RootNode: &merkleNodes[0]}

	return &merkleTree
}
