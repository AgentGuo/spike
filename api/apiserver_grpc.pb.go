// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.20.3
// source: api/apiserver.proto

package api

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
	SpikeService_CallFunction_FullMethodName         = "/spike.SpikeService/CallFunction"
	SpikeService_CreateFunction_FullMethodName       = "/spike.SpikeService/CreateFunction"
	SpikeService_DeleteFunction_FullMethodName       = "/spike.SpikeService/DeleteFunction"
	SpikeService_GetAllFunctions_FullMethodName      = "/spike.SpikeService/GetAllFunctions"
	SpikeService_GetFunctionResources_FullMethodName = "/spike.SpikeService/GetFunctionResources"
)

// SpikeServiceClient is the client API for SpikeService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// The spike service definition.
type SpikeServiceClient interface {
	// Calls a function.
	CallFunction(ctx context.Context, in *CallFunctionRequest, opts ...grpc.CallOption) (*CallFunctionResponse, error)
	// Adds a new function.
	CreateFunction(ctx context.Context, in *CreateFunctionRequest, opts ...grpc.CallOption) (*CreateFunctionResponse, error)
	// Deletes an existing function.
	DeleteFunction(ctx context.Context, in *DeleteFunctionRequest, opts ...grpc.CallOption) (*DeleteFunctionResponse, error)
	GetAllFunctions(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*GetAllFunctionsResponse, error)
	GetFunctionResources(ctx context.Context, in *GetFunctionResourcesRequest, opts ...grpc.CallOption) (*GetFunctionResourcesResponse, error)
}

type spikeServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewSpikeServiceClient(cc grpc.ClientConnInterface) SpikeServiceClient {
	return &spikeServiceClient{cc}
}

func (c *spikeServiceClient) CallFunction(ctx context.Context, in *CallFunctionRequest, opts ...grpc.CallOption) (*CallFunctionResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CallFunctionResponse)
	err := c.cc.Invoke(ctx, SpikeService_CallFunction_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *spikeServiceClient) CreateFunction(ctx context.Context, in *CreateFunctionRequest, opts ...grpc.CallOption) (*CreateFunctionResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateFunctionResponse)
	err := c.cc.Invoke(ctx, SpikeService_CreateFunction_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *spikeServiceClient) DeleteFunction(ctx context.Context, in *DeleteFunctionRequest, opts ...grpc.CallOption) (*DeleteFunctionResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteFunctionResponse)
	err := c.cc.Invoke(ctx, SpikeService_DeleteFunction_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *spikeServiceClient) GetAllFunctions(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*GetAllFunctionsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetAllFunctionsResponse)
	err := c.cc.Invoke(ctx, SpikeService_GetAllFunctions_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *spikeServiceClient) GetFunctionResources(ctx context.Context, in *GetFunctionResourcesRequest, opts ...grpc.CallOption) (*GetFunctionResourcesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetFunctionResourcesResponse)
	err := c.cc.Invoke(ctx, SpikeService_GetFunctionResources_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SpikeServiceServer is the server API for SpikeService service.
// All implementations must embed UnimplementedSpikeServiceServer
// for forward compatibility.
//
// The spike service definition.
type SpikeServiceServer interface {
	// Calls a function.
	CallFunction(context.Context, *CallFunctionRequest) (*CallFunctionResponse, error)
	// Adds a new function.
	CreateFunction(context.Context, *CreateFunctionRequest) (*CreateFunctionResponse, error)
	// Deletes an existing function.
	DeleteFunction(context.Context, *DeleteFunctionRequest) (*DeleteFunctionResponse, error)
	GetAllFunctions(context.Context, *Empty) (*GetAllFunctionsResponse, error)
	GetFunctionResources(context.Context, *GetFunctionResourcesRequest) (*GetFunctionResourcesResponse, error)
	mustEmbedUnimplementedSpikeServiceServer()
}

// UnimplementedSpikeServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedSpikeServiceServer struct{}

func (UnimplementedSpikeServiceServer) CallFunction(context.Context, *CallFunctionRequest) (*CallFunctionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CallFunction not implemented")
}
func (UnimplementedSpikeServiceServer) CreateFunction(context.Context, *CreateFunctionRequest) (*CreateFunctionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateFunction not implemented")
}
func (UnimplementedSpikeServiceServer) DeleteFunction(context.Context, *DeleteFunctionRequest) (*DeleteFunctionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteFunction not implemented")
}
func (UnimplementedSpikeServiceServer) GetAllFunctions(context.Context, *Empty) (*GetAllFunctionsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllFunctions not implemented")
}
func (UnimplementedSpikeServiceServer) GetFunctionResources(context.Context, *GetFunctionResourcesRequest) (*GetFunctionResourcesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFunctionResources not implemented")
}
func (UnimplementedSpikeServiceServer) mustEmbedUnimplementedSpikeServiceServer() {}
func (UnimplementedSpikeServiceServer) testEmbeddedByValue()                      {}

// UnsafeSpikeServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SpikeServiceServer will
// result in compilation errors.
type UnsafeSpikeServiceServer interface {
	mustEmbedUnimplementedSpikeServiceServer()
}

func RegisterSpikeServiceServer(s grpc.ServiceRegistrar, srv SpikeServiceServer) {
	// If the following call pancis, it indicates UnimplementedSpikeServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&SpikeService_ServiceDesc, srv)
}

func _SpikeService_CallFunction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CallFunctionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SpikeServiceServer).CallFunction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SpikeService_CallFunction_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SpikeServiceServer).CallFunction(ctx, req.(*CallFunctionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SpikeService_CreateFunction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateFunctionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SpikeServiceServer).CreateFunction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SpikeService_CreateFunction_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SpikeServiceServer).CreateFunction(ctx, req.(*CreateFunctionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SpikeService_DeleteFunction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteFunctionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SpikeServiceServer).DeleteFunction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SpikeService_DeleteFunction_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SpikeServiceServer).DeleteFunction(ctx, req.(*DeleteFunctionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SpikeService_GetAllFunctions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SpikeServiceServer).GetAllFunctions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SpikeService_GetAllFunctions_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SpikeServiceServer).GetAllFunctions(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _SpikeService_GetFunctionResources_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetFunctionResourcesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SpikeServiceServer).GetFunctionResources(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SpikeService_GetFunctionResources_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SpikeServiceServer).GetFunctionResources(ctx, req.(*GetFunctionResourcesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// SpikeService_ServiceDesc is the grpc.ServiceDesc for SpikeService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SpikeService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "spike.SpikeService",
	HandlerType: (*SpikeServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CallFunction",
			Handler:    _SpikeService_CallFunction_Handler,
		},
		{
			MethodName: "CreateFunction",
			Handler:    _SpikeService_CreateFunction_Handler,
		},
		{
			MethodName: "DeleteFunction",
			Handler:    _SpikeService_DeleteFunction_Handler,
		},
		{
			MethodName: "GetAllFunctions",
			Handler:    _SpikeService_GetAllFunctions_Handler,
		},
		{
			MethodName: "GetFunctionResources",
			Handler:    _SpikeService_GetFunctionResources_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/apiserver.proto",
}
