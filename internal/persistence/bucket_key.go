package persistence

/*
 * Copyright 2020 Information Wants To Be Free
 * Visit: https://github.com/iwtbf
 *
 * This project is licensed under the terms of the Apache 2.0 License.
 */

type BlocksBucketKey string

const (
	BlockIndex          BlocksBucketKey = "b"
	FileInformation     BlocksBucketKey = "f"
	LastBlockFileNumber BlocksBucketKey = "l"
	Reindexing          BlocksBucketKey = "R"
	Flags               BlocksBucketKey = "F"
	TransactionHash     BlocksBucketKey = "t"
)

type ChainstateBucketKey string

const (
	UnspentTransactionOutput ChainstateBucketKey = "c"
	BlockHash                ChainstateBucketKey = "B"
)
