package farmer

/*
 * Copyright 2020 Information Wants To Be Free
 * Visit: https://github.com/iwtbf
 *
 * This project is licensed under the terms of the Apache 2.0 License.
 */

var (
	knownCoreRelays = []coreRelay{
		{
			address:       "/ip4/127.0.0.1/tcp/2000/p2p/QmNxGGyb9CJUnpzsFZwk1vUz6PEEQxz6xA6DmGc2tcDBud",
			publicKeyHash: "1LvLm5G4aJcyKMHGStRDVDNzt97TNjM4Vu",
		},
	}
)

type coreRelay struct {
	address       string
	publicKeyHash string
}

func prepareCoreRelays(hostAddress string) []coreRelay {
	var coreRelays = make([]coreRelay, len(knownCoreRelays))
	copy(coreRelays, knownCoreRelays)

	for iterator := 0; iterator < len(coreRelays); iterator++ {
		if coreRelays[iterator].address == hostAddress {
			coreRelays[len(coreRelays)-1], coreRelays[iterator] = coreRelays[iterator], coreRelays[len(coreRelays)-1]
			coreRelays = coreRelays[:len(coreRelays)-1]
		}
	}

	return coreRelays
}
