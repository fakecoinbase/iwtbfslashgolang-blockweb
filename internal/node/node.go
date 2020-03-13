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
	libp2pGrpc "github.com/iwtbf/golang-blockweb/pkg/libp2p-grpc"
	"github.com/libp2p/go-libp2p"
	circuit "github.com/libp2p/go-libp2p-circuit"
	core "github.com/libp2p/go-libp2p-core"
	"github.com/libp2p/go-libp2p-core/crypto"
	libp2pHost "github.com/libp2p/go-libp2p-core/host"
)

type node struct {
	host         core.Host
	grpcProtocol *libp2pGrpc.GRPCProtocol
}

func startLibp2pGrcpHost(port int16, privKey crypto.PrivKey) (libp2pHost.Host, *libp2pGrpc.GRPCProtocol) {
	ctx := context.Background()

	host, err := libp2p.New(
		ctx,
		libp2p.Identity(privKey),
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", port)),
		// TODO: Maybe not always relay?
		libp2p.EnableRelay(circuit.OptDiscovery),
	)
	if err != nil {
		panic(err)
	}

	return host, libp2pGrpc.NewGRPCProtocol(ctx, host)
}
