package node

/*
 * Copyright 2020 Information Wants To Be Free
 * Visit: https://github.com/iwtbf
 *
 * This project is licensed under the terms of the Apache 2.0 License.
 */

import (
	"github.com/ipfs/go-log"
	farmerClient "github.com/iwtbf/golang-blockweb/internal/farmer/client"
	"github.com/iwtbf/golang-blockweb/internal/keygen"
	"os"
	"os/signal"
	"syscall"
)

var (
	logger = log.Logger("node")
)

type relay struct {
	*node
}

func bootRelay(port int16, keyPair keygen.KeyPair) {
	logger.Infof("Booting relay (full node) on port %d..", port)

	host, hostAddress, grpcProtocol := startLibp2pGrcpHost(port, keyPair)
	RegisterRelayServer(grpcProtocol.GetGRPCServer(), &relay{node: &node{PeerID: host.ID()}})

	relayAddress, relayPublicKeyHash := farmerClient.RequestRandomCoreRelayInformation(hostAddress.String())

	if relayAddress != hostAddress.String() {
		connectToCoreRelay(relayAddress, relayPublicKeyHash)
	} else {
		logger.Warning("I am the genesis relay - I won't bow to anyone!")
	}

	logger.Info("Done.")

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	logger.Info("Received signal, shutting down...")

	if err := host.Close(); err != nil {
		panic(err)
	}
}
