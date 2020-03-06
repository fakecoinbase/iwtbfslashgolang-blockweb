package node

/*
 * Copyright 2020 Information Wants To Be Free
 * Visit: https://github.com/iwtbf
 *
 * This project is licensed under the terms of the Apache 2.0 License.
 */

import (
	"github.com/ipfs/go-log"
	"github.com/libp2p/go-libp2p-core/crypto"
	"os"
	"os/signal"
	"syscall"
)

var (
	logger = log.Logger("node")
)

type relay struct {
	*node
	hashTable *DistributedHashTable
}

func newGrpcRelay(port int16, privKey crypto.PrivKey) *relay {
	logger.Debug("Starting GRPC service..")
	host, grpcProtocol := startLibp2pGrcpHost(port, privKey)
	relay := &relay{node: &node{host: host, grpcProtocol: grpcProtocol}}

	RegisterRelayServer(grpcProtocol.GetGRPCServer(), relay)
	logger.Debug("GRPC service started.")

	return relay
}

func bootRelay(port int16, privKey crypto.PrivKey) {
	logger.Infof("Booting relay (full node) on port %d..", port)

	relay := newGrpcRelay(port, privKey)
	relay.hashTable = newDistributedHashTable(relay)
	relay.hashTable.synchronize(relay.host)

	// TODO: Bootstrap blockchain

	relay.hashTable.announce()

	logger.Info("Boot successful. relay is ready to improve the world.")

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	logger.Info("Received signal, shutting down...")

	if err := relay.host.Close(); err != nil {
		panic(err)
	}
}
