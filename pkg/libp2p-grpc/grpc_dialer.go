package libp2p_grpc

/*
 * Copyright 2020 Information Wants To Be Free
 * Visit: https://github.com/iwtbf
 *
 * This project is licensed under the terms of the Apache 2.0 License.
 */

import (
	"context"
	"github.com/libp2p/go-libp2p-core/peer"
	"google.golang.org/grpc"
	"net"
)

func (grpcProtocol *GRPCProtocol) GetDialOption(ctx context.Context) grpc.DialOption {
	return grpc.WithContextDialer(func(subCtx context.Context, peerID string) (net.Conn, error) {
		id, err := peer.Decode(peerID)
		if err != nil {
			return nil, err
		}

		if err = grpcProtocol.host.Connect(subCtx, peer.AddrInfo{
			ID: id,
		}); err != nil {
			return nil, err
		}

		stream, err := grpcProtocol.host.NewStream(ctx, id, Protocol)
		if err != nil {
			return nil, err
		}

		return &streamConnection{Stream: stream}, nil
	})
}

func (grpcProtocol *GRPCProtocol) Dial(ctx context.Context, peerID peer.ID, dialOpts ...grpc.DialOption) (*grpc.ClientConn, error) {
	dialOpsPrepended := append([]grpc.DialOption{grpcProtocol.GetDialOption(ctx)}, dialOpts...)
	return grpc.DialContext(ctx, peerID.Pretty(), dialOpsPrepended...)
}
