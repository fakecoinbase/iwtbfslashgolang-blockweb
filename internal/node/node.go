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
	"github.com/iwtbf/golang-blockweb/internal/keygen"
	libp2pGrpc "github.com/iwtbf/golang-blockweb/pkg/libp2p-grpc"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiformats/go-multiaddr"
)

type node struct {
	PeerID peer.ID
}

func startLibp2pGrcpHost(port int16, keyPair keygen.KeyPair) (host.Host, multiaddr.Multiaddr, *libp2pGrpc.GRPCProtocol) {
	ctx := context.Background()

	ecdsaPrivateKey, _, err := crypto.ECDSAKeyPairFromKey(&keyPair.PrivateKey)
	if err != nil {
		panic(err)
	}

	host, err := libp2p.New(
		ctx,
		libp2p.Identity(ecdsaPrivateKey),
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", port)),
	)
	if err != nil {
		panic(err)
	}

	peerInfo := peer.AddrInfo{ID: host.ID(), Addrs: host.Addrs()}
	addresses, err := peer.AddrInfoToP2pAddrs(&peerInfo)
	if err != nil {
		panic(err)
	}
	logger.Debugf("Relay address: %v", addresses[0])

	return host, addresses[0], libp2pGrpc.NewGRPCProtocol(ctx, host)
}

func connectToCoreRelay(relayAddress string, relayPublicKeyHash []byte) {
	logger.Debugf("Trying to connect to bootstrap peer: {address: %s, publicKeyHash: %s}", relayAddress, string(relayPublicKeyHash))
}
