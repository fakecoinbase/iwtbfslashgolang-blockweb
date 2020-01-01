package pkg

/*
 * Copyright 2020 Information Wants To Be Free
 * Visit: https://github.com/iwtbf
 *
 * This project is licensed under the terms of the Apache 2.0 License.
 */

import "strconv"

/*
 * Convert an int64 to its hexadecimal []byte representation.
 */
func IntToHex(int int64) []byte {
	return []byte(strconv.FormatInt(int, 16))
}
