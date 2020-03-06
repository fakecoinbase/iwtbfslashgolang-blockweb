// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/node/relay.proto

package node

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

func init() { proto.RegisterFile("proto/node/relay.proto", fileDescriptor_33615706ef1962af) }

var fileDescriptor_33615706ef1962af = []byte{
	// 122 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x2b, 0x28, 0xca, 0x2f,
	0xc9, 0xd7, 0xcf, 0xcb, 0x4f, 0x49, 0xd5, 0x2f, 0x4a, 0xcd, 0x49, 0xac, 0xd4, 0x03, 0x0b, 0x08,
	0x71, 0xa4, 0xe4, 0x15, 0xc7, 0x17, 0xa7, 0xa6, 0xa6, 0x48, 0x49, 0x20, 0xa9, 0x28, 0x4b, 0x2d,
	0x2a, 0xce, 0xcc, 0xcf, 0x83, 0xa8, 0x31, 0x72, 0xe2, 0x62, 0x0d, 0x02, 0x69, 0x11, 0xb2, 0xe4,
	0xe2, 0x77, 0xad, 0x48, 0xce, 0x48, 0xcc, 0x4b, 0x4f, 0x0d, 0x83, 0xa8, 0x10, 0x12, 0xd4, 0x83,
	0x19, 0xa0, 0x07, 0x15, 0x92, 0xc2, 0x14, 0x52, 0x62, 0x70, 0xe2, 0x8f, 0xe2, 0xcd, 0xcc, 0x2b,
	0x49, 0x2d, 0xca, 0x4b, 0xcc, 0x01, 0x5b, 0x91, 0xc4, 0x06, 0x36, 0xdb, 0x18, 0x10, 0x00, 0x00,
	0xff, 0xff, 0x4f, 0xc8, 0x62, 0x7a, 0x99, 0x00, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// RelayClient is the client API for relay service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type RelayClient interface {
	ExchangeVersion(ctx context.Context, in *Version, opts ...grpc.CallOption) (*Version, error)
}

type relayClient struct {
	cc grpc.ClientConnInterface
}

func NewRelayClient(cc grpc.ClientConnInterface) RelayClient {
	return &relayClient{cc}
}

func (c *relayClient) ExchangeVersion(ctx context.Context, in *Version, opts ...grpc.CallOption) (*Version, error) {
	out := new(Version)
	err := c.cc.Invoke(ctx, "/farmer.relay/ExchangeVersion", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RelayServer is the server API for relay service.
type RelayServer interface {
	ExchangeVersion(context.Context, *Version) (*Version, error)
}

// UnimplementedRelayServer can be embedded to have forward compatible implementations.
type UnimplementedRelayServer struct {
}

func (*UnimplementedRelayServer) ExchangeVersion(ctx context.Context, req *Version) (*Version, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ExchangeVersion not implemented")
}

func RegisterRelayServer(s *grpc.Server, srv RelayServer) {
	s.RegisterService(&_Relay_serviceDesc, srv)
}

func _Relay_ExchangeVersion_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Version)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RelayServer).ExchangeVersion(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/farmer.relay/ExchangeVersion",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RelayServer).ExchangeVersion(ctx, req.(*Version))
	}
	return interceptor(ctx, in, info, handler)
}

var _Relay_serviceDesc = grpc.ServiceDesc{
	ServiceName: "farmer.relay",
	HandlerType: (*RelayServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ExchangeVersion",
			Handler:    _Relay_ExchangeVersion_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/node/relay.proto",
}