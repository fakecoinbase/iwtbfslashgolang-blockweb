package node

/*
 * Copyright 2020 Information Wants To Be Free
 * Visit: https://github.com/iwtbf
 *
 * This project is licensed under the terms of the Apache 2.0 License.
 */

import (
	"context"
	discovery "github.com/libp2p/go-libp2p-discovery"
	"google.golang.org/grpc"
)

func (relay relay) bootstrapBlockchain() {
	logger.Info("Bootstrapping blockchain..")

	peerInfos, err := discovery.FindPeers(context.Background(), relay.hashTable.routingDiscovery, relayRendezvous)
	if err != nil {
		panic(err)
	}

	// TODO: Detect unsynced blockchain
	for _, peer := range peerInfos {
		if peer.ID == relay.host.ID() {
			continue
		}

		logger.Debugf("Trying to dial peer '%s'", peer.ID)

		// TODO: TLS
		_, err := relay.grpcProtocol.Dial(context.Background(), peer.ID, grpc.WithBlock())
		if err != nil {
			logger.Warningf("Connecting to peer '%s' failed: %v", peer.ID, err)
			continue
		}

		logger.Debug("Connection successful, requesting version..")
	}

	logger.Debug("Blockchain successfully synchronized.")
}
