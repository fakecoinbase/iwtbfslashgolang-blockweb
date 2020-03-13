package keygen

/*
 * Copyright 2020 Information Wants To Be Free
 * Visit: https://github.com/iwtbf
 *
 * This project is licensed under the terms of the Apache 2.0 License.
 */

import (
	"bufio"
	"fmt"
	"os"
)

const (
	disclaimer = `DISCLAIMER:
This program * can * be used to generate key pairs and certificates for 'golang-blockweb'.
It is open source software, but there is no liability of the contributors or guarantee for the security of the generated keys. Please report any problems or vulnerabilities at https://github.com/iwtbf/golang-blockweb/issues/new.
You could use any other key tool with elliptic P256 curve cryptography (ECDSA) + x509 + PEM. Use the command 'read-key' to verify the generated key.
By continuing, you confirm that you have understood the above points (y / N):`
	retry = `Sorry, I didn't understand this. Please try again (y / N):`
)

func acceptedDisclaimer(answer string) bool {
	return answer == "y\n" || answer == "Y\n"
}

func acceptDisclaimer(disclaimer, retry string) bool {
	fmt.Println(disclaimer)

	reader := bufio.NewReader(os.Stdin)
	for answer, err := reader.ReadString('\n'); err != nil || answer == "\n" || !acceptedDisclaimer(answer); answer, err = reader.ReadString('\n') {
		if err != nil {
			panic(err)
		}

		if answer == "\n" || answer == "n\n" || answer == "N\n" {
			return false
		}

		fmt.Println(retry)
	}

	return true
}
