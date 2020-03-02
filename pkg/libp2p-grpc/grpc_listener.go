package libp2p_grpc

/*
 * Copyright 2020 Information Wants To Be Free
 * Visit: https://github.com/iwtbf
 *
 * This project is licensed under the terms of the Apache 2.0 License.
 */

import (
	"context"
	manet "github.com/multiformats/go-multiaddr-net"
	"io"
	"net"
)

/*
 * Type implementing net.Listener.
 */
type grpcListener struct {
	*GRPCProtocol
	ctx        context.Context
	cancelFunc context.CancelFunc
}

func newGrpcListener(grpcProtocol *GRPCProtocol) net.Listener {
	grpcListener := &grpcListener{
		GRPCProtocol: grpcProtocol,
	}
	grpcListener.ctx, grpcListener.cancelFunc = context.WithCancel(grpcProtocol.ctx)
	return grpcListener
}

func (grpcListener *grpcListener) Accept() (net.Conn, error) {
	select {
	case <-grpcListener.ctx.Done():
		return nil, io.EOF
	case stream := <-grpcListener.streamCh:
		return &streamConnection{Stream: stream}, nil
	}
}

func (grpcListener *grpcListener) Addr() net.Addr {
	listenAddresses := grpcListener.host.Network().ListenAddresses()
	if len(listenAddresses) > 0 {
		for _, listenAddress := range listenAddresses {
			if netAddr, err := manet.ToNetAddr(listenAddress); err == nil {
				return netAddr
			}
		}
	}

	panic("No listening address found! Make sure your peer started correctly.")
}

func (grpcListener *grpcListener) Close() error {
	grpcListener.cancelFunc()
	return nil
}
