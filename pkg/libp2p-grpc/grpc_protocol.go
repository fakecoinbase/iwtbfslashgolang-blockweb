package libp2p_grpc

/*
 * Copyright 2020 Information Wants To Be Free
 * Visit: https://github.com/iwtbf
 *
 * This project is licensed under the terms of the Apache 2.0 License.
 */

import (
	"context"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/protocol"
	"google.golang.org/grpc"
)

const Protocol protocol.ID = "/grpc/1.0.0"

/*
 * GRPC-transported protocol handler for libp2p hosts (github.com/libp2p/go-libp2p-core/host#Host).
 */
type GRPCProtocol struct {
	ctx        context.Context
	host       host.Host
	grpcServer *grpc.Server
	streamCh   chan network.Stream
}

func (grpcProtocol *GRPCProtocol) GetGRPCServer() *grpc.Server {
	return grpcProtocol.grpcServer
}

func (grpcProtocol *GRPCProtocol) HandleStream(stream network.Stream) {
	select {
	case <-grpcProtocol.ctx.Done():
		return
	case grpcProtocol.streamCh <- stream:
	}
}

func NewGRPCProtocol(ctx context.Context, host host.Host) *GRPCProtocol {
	grpcServer := grpc.NewServer()
	grpcProtocol := &GRPCProtocol{
		ctx:        ctx,
		host:       host,
		grpcServer: grpcServer,
		streamCh:   make(chan network.Stream),
	}
	host.SetStreamHandler(Protocol, grpcProtocol.HandleStream)
	go grpcServer.Serve(newGrpcListener(grpcProtocol))
	return grpcProtocol
}
