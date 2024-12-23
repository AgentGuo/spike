// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.20.3
// source: pkg/worker/spikeworker.proto

package worker

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	SpikeWorkerService_CallWorkerFunction_FullMethodName = "/SpikeWorkerService/CallWorkerFunction"
)

// SpikeWorkerServiceClient is the client API for SpikeWorkerService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// spike worker暴露给spike server的服务
type SpikeWorkerServiceClient interface {
	CallWorkerFunction(ctx context.Context, in *CallWorkerFunctionReq, opts ...grpc.CallOption) (*CallWorkerFunctionResp, error)
}

type spikeWorkerServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewSpikeWorkerServiceClient(cc grpc.ClientConnInterface) SpikeWorkerServiceClient {
	return &spikeWorkerServiceClient{cc}
}

func (c *spikeWorkerServiceClient) CallWorkerFunction(ctx context.Context, in *CallWorkerFunctionReq, opts ...grpc.CallOption) (*CallWorkerFunctionResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CallWorkerFunctionResp)
	err := c.cc.Invoke(ctx, SpikeWorkerService_CallWorkerFunction_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SpikeWorkerServiceServer is the server API for SpikeWorkerService service.
// All implementations must embed UnimplementedSpikeWorkerServiceServer
// for forward compatibility.
//
// spike worker暴露给spike server的服务
type SpikeWorkerServiceServer interface {
	CallWorkerFunction(context.Context, *CallWorkerFunctionReq) (*CallWorkerFunctionResp, error)
	mustEmbedUnimplementedSpikeWorkerServiceServer()
}

// UnimplementedSpikeWorkerServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedSpikeWorkerServiceServer struct{}

func (UnimplementedSpikeWorkerServiceServer) CallWorkerFunction(context.Context, *CallWorkerFunctionReq) (*CallWorkerFunctionResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CallWorkerFunction not implemented")
}
func (UnimplementedSpikeWorkerServiceServer) mustEmbedUnimplementedSpikeWorkerServiceServer() {}
func (UnimplementedSpikeWorkerServiceServer) testEmbeddedByValue()                            {}

// UnsafeSpikeWorkerServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SpikeWorkerServiceServer will
// result in compilation errors.
type UnsafeSpikeWorkerServiceServer interface {
	mustEmbedUnimplementedSpikeWorkerServiceServer()
}

func RegisterSpikeWorkerServiceServer(s grpc.ServiceRegistrar, srv SpikeWorkerServiceServer) {
	// If the following call pancis, it indicates UnimplementedSpikeWorkerServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&SpikeWorkerService_ServiceDesc, srv)
}

func _SpikeWorkerService_CallWorkerFunction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CallWorkerFunctionReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SpikeWorkerServiceServer).CallWorkerFunction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SpikeWorkerService_CallWorkerFunction_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SpikeWorkerServiceServer).CallWorkerFunction(ctx, req.(*CallWorkerFunctionReq))
	}
	return interceptor(ctx, in, info, handler)
}

// SpikeWorkerService_ServiceDesc is the grpc.ServiceDesc for SpikeWorkerService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SpikeWorkerService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "SpikeWorkerService",
	HandlerType: (*SpikeWorkerServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CallWorkerFunction",
			Handler:    _SpikeWorkerService_CallWorkerFunction_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pkg/worker/spikeworker.proto",
}
