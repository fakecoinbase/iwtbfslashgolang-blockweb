package farmer

/*
 * Copyright 2020 Information Wants To Be Free
 * Visit: https://github.com/iwtbf
 *
 * This project is licensed under the terms of the Apache 2.0 License.
 */

import (
	"context"
	"github.com/ipfs/go-log"
	"math/rand"
)

var (
	farmerLogger    = log.Logger("farmer")
)

type farmer struct {
	UnimplementedFarmerServer
}

func (farmer *farmer) RequestSeed(ctx context.Context, seedRequest *SeedRequest) (*SeedReply, error) {
	coreRelay := knownCoreRelays[0]

	if len(knownCoreRelays) > 1 {
		coreRelays := prepareCoreRelays(seedRequest.HostAddress)
		coreRelay = coreRelays[rand.Intn(len(coreRelays))]
	}

	farmerLogger.Debugf("RequestSeed -> SeedReply{Address: %s, PublicKeyHash: %s}", coreRelay.address, coreRelay.publicKeyHash)

	return &SeedReply{Address: coreRelay.address, PublicKeyHash: []byte(coreRelay.publicKeyHash)}, nil
}
