// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.19.4
// source: grpc-cloud.proto

//option go_package = "./;cloudservice"; //dir of create proto-file

package __

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	CloudExchange_ProcessCloud_FullMethodName = "/cloudservice.CloudExchange/processCloud"
)

// CloudExchangeClient is the client API for CloudExchange service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CloudExchangeClient interface {
	ProcessCloud(ctx context.Context, opts ...grpc.CallOption) (CloudExchange_ProcessCloudClient, error)
}

type cloudExchangeClient struct {
	cc grpc.ClientConnInterface
}

func NewCloudExchangeClient(cc grpc.ClientConnInterface) CloudExchangeClient {
	return &cloudExchangeClient{cc}
}

func (c *cloudExchangeClient) ProcessCloud(ctx context.Context, opts ...grpc.CallOption) (CloudExchange_ProcessCloudClient, error) {
	stream, err := c.cc.NewStream(ctx, &CloudExchange_ServiceDesc.Streams[0], CloudExchange_ProcessCloud_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &cloudExchangeProcessCloudClient{stream}
	return x, nil
}

type CloudExchange_ProcessCloudClient interface {
	Send(*RequestIO) error
	Recv() (*StatusIO, error)
	grpc.ClientStream
}

type cloudExchangeProcessCloudClient struct {
	grpc.ClientStream
}

func (x *cloudExchangeProcessCloudClient) Send(m *RequestIO) error {
	return x.ClientStream.SendMsg(m)
}

func (x *cloudExchangeProcessCloudClient) Recv() (*StatusIO, error) {
	m := new(StatusIO)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// CloudExchangeServer is the server API for CloudExchange service.
// All implementations should embed UnimplementedCloudExchangeServer
// for forward compatibility
type CloudExchangeServer interface {
	ProcessCloud(CloudExchange_ProcessCloudServer) error
}

// UnimplementedCloudExchangeServer should be embedded to have forward compatible implementations.
type UnimplementedCloudExchangeServer struct {
}

func (UnimplementedCloudExchangeServer) ProcessCloud(CloudExchange_ProcessCloudServer) error {
	return status.Errorf(codes.Unimplemented, "method ProcessCloud not implemented")
}

// UnsafeCloudExchangeServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CloudExchangeServer will
// result in compilation errors.
type UnsafeCloudExchangeServer interface {
	mustEmbedUnimplementedCloudExchangeServer()
}

func RegisterCloudExchangeServer(s grpc.ServiceRegistrar, srv CloudExchangeServer) {
	s.RegisterService(&CloudExchange_ServiceDesc, srv)
}

func _CloudExchange_ProcessCloud_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(CloudExchangeServer).ProcessCloud(&cloudExchangeProcessCloudServer{stream})
}

type CloudExchange_ProcessCloudServer interface {
	Send(*StatusIO) error
	Recv() (*RequestIO, error)
	grpc.ServerStream
}

type cloudExchangeProcessCloudServer struct {
	grpc.ServerStream
}

func (x *cloudExchangeProcessCloudServer) Send(m *StatusIO) error {
	return x.ServerStream.SendMsg(m)
}

func (x *cloudExchangeProcessCloudServer) Recv() (*RequestIO, error) {
	m := new(RequestIO)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// CloudExchange_ServiceDesc is the grpc.ServiceDesc for CloudExchange service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CloudExchange_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "cloudservice.CloudExchange",
	HandlerType: (*CloudExchangeServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "processCloud",
			Handler:       _CloudExchange_ProcessCloud_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "grpc-cloud.proto",
}
