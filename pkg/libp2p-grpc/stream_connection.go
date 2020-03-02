package libp2p_grpc

import (
	"github.com/libp2p/go-libp2p-core/network"
	manet "github.com/multiformats/go-multiaddr-net"
	"net"
)

/*
 * Type implementing network.Stream
 */
type streamConnection struct {
	network.Stream
}

func (streamConnection *streamConnection) LocalAddr() net.Addr {
	localAddr, err := manet.ToNetAddr(streamConnection.Stream.Conn().LocalMultiaddr())
	if err != nil {
		panic(err)
	}

	return localAddr
}

func (streamConnection *streamConnection) RemoteAddr() net.Addr {
	remoteAddr, err := manet.ToNetAddr(streamConnection.Stream.Conn().RemoteMultiaddr())
	if err != nil {
		panic(err)
	}

	return remoteAddr
}
