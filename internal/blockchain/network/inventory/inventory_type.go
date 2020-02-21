package inventory

/*
 * Copyright 2020 Information Wants To Be Free
 * Visit: https://github.com/iwtbf
 *
 * This project is licensed under the terms of the Apache 2.0 License.
 */

type InventoryType string

const (
	Block       InventoryType = "block"
	Transaction InventoryType = "transaction"
)
