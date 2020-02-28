package keygen

/*
 * Copyright 2020 Information Wants To Be Free
 * Visit: https://github.com/iwtbf
 *
 * This project is licensed under the terms of the Apache 2.0 License.
 */

import (
	"bufio"
	"os"
)

func acceptedDisclaimer(answer string) bool {
	return answer == "y\n" || answer == "Y\n"
}

func acceptDisclaimer(disclaimer, retry string) bool {
	println(disclaimer)

	reader := bufio.NewReader(os.Stdin)
	for answer, err := reader.ReadString('\n'); err != nil || answer == "\n" || !acceptedDisclaimer(answer); answer, err = reader.ReadString('\n') {
		if err != nil {
			panic(err)
		}

		if answer == "\n" || answer == "n\n" || answer == "N\n" {
			return false
		}

		println(retry)
	}

	return true
}
