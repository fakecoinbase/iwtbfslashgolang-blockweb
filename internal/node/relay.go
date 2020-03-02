package node

/*
 * Copyright 2020 Information Wants To Be Free
 * Visit: https://github.com/iwtbf
 *
 * This project is licensed under the terms of the Apache 2.0 License.
 */

import (
	"context"
	"fmt"
	"github.com/ipfs/go-log"
	"github.com/iwtbf/golang-blockweb/internal/keygen"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	peerstore "github.com/libp2p/go-libp2p-core/peer"
	"os"
	"os/signal"
	"syscall"
)

var (
	logger = log.Logger("node")
)

func bootRelay(port int16, keyPair keygen.KeyPair) {
	logger.Infof("Booting full node (relay) on port %d..", port)

	ctx := context.Background()

	ecdsaPrivateKey, _, err := crypto.ECDSAKeyPairFromKey(&keyPair.PrivateKey)
	if err != nil {
		panic(err)
	}

	node, err := libp2p.New(
		ctx,
		libp2p.Identity(ecdsaPrivateKey),
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", port)),
	)
	if err != nil {
		panic(err)
	}

	peerInfo := peerstore.AddrInfo{ID: node.ID(), Addrs: node.Addrs(),}
	addresses, err := peerstore.AddrInfoToP2pAddrs(&peerInfo)
	if err != nil {
		panic(err)
	}
	logger.Debugf("Node address: %v", addresses[0])
	logger.Info("Done.")

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	logger.Info("Received signal, shutting down...")

	if err := node.Close(); err != nil {
		panic(err)
	}
}
