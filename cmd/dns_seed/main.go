package main

/*
 * Copyright 2020 Information Wants To Be Free
 * Visit: https://github.com/iwtbf
 *
 * This project is licensed under the terms of the Apache 2.0 License.
 */

import dns "github.com/iwtbf/golang-blockweb/internal/dns_seed"

func main() {
	(*dns.NewDNSSeed()).Run()
}
