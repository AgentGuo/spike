// protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative api/apiserver.proto

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
	FaaS_CallFunction_FullMethodName         = "/faas.FaaS/CallFunction"
	FaaS_CreateFunction_FullMethodName       = "/faas.FaaS/CreateFunction"
	FaaS_DeleteFunction_FullMethodName       = "/faas.FaaS/DeleteFunction"
	FaaS_GetAllFunctions_FullMethodName      = "/faas.FaaS/GetAllFunctions"
	FaaS_GetFunctionResources_FullMethodName = "/faas.FaaS/GetFunctionResources"
)

// FaaSClient is the client API for FaaS service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// The faas service definition.
type FaaSClient interface {
	// Calls a function.
	CallFunction(ctx context.Context, in *CallFunctionRequest, opts ...grpc.CallOption) (*CallFunctionResponse, error)
	// Adds a new function.
	CreateFunction(ctx context.Context, in *CreateFunctionRequest, opts ...grpc.CallOption) (*CreateFunctionResponse, error)
	// Deletes an existing function.
	DeleteFunction(ctx context.Context, in *DeleteFunctionRequest, opts ...grpc.CallOption) (*DeleteFunctionResponse, error)
	GetAllFunctions(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*GetAllFunctionsResponse, error)
	GetFunctionResources(ctx context.Context, in *GetFunctionResourcesRequest, opts ...grpc.CallOption) (*GetFunctionResourcesResponse, error)
}

type faaSClient struct {
	cc grpc.ClientConnInterface
}

func NewFaaSClient(cc grpc.ClientConnInterface) FaaSClient {
	return &faaSClient{cc}
}

func (c *faaSClient) CallFunction(ctx context.Context, in *CallFunctionRequest, opts ...grpc.CallOption) (*CallFunctionResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CallFunctionResponse)
	err := c.cc.Invoke(ctx, FaaS_CallFunction_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *faaSClient) CreateFunction(ctx context.Context, in *CreateFunctionRequest, opts ...grpc.CallOption) (*CreateFunctionResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateFunctionResponse)
	err := c.cc.Invoke(ctx, FaaS_CreateFunction_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *faaSClient) DeleteFunction(ctx context.Context, in *DeleteFunctionRequest, opts ...grpc.CallOption) (*DeleteFunctionResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteFunctionResponse)
	err := c.cc.Invoke(ctx, FaaS_DeleteFunction_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *faaSClient) GetAllFunctions(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*GetAllFunctionsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetAllFunctionsResponse)
	err := c.cc.Invoke(ctx, FaaS_GetAllFunctions_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *faaSClient) GetFunctionResources(ctx context.Context, in *GetFunctionResourcesRequest, opts ...grpc.CallOption) (*GetFunctionResourcesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetFunctionResourcesResponse)
	err := c.cc.Invoke(ctx, FaaS_GetFunctionResources_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FaaSServer is the server API for FaaS service.
// All implementations must embed UnimplementedFaaSServer
// for forward compatibility.
//
// The faas service definition.
type FaaSServer interface {
	// Calls a function.
	CallFunction(context.Context, *CallFunctionRequest) (*CallFunctionResponse, error)
	// Adds a new function.
	CreateFunction(context.Context, *CreateFunctionRequest) (*CreateFunctionResponse, error)
	// Deletes an existing function.
	DeleteFunction(context.Context, *DeleteFunctionRequest) (*DeleteFunctionResponse, error)
	GetAllFunctions(context.Context, *Empty) (*GetAllFunctionsResponse, error)
	GetFunctionResources(context.Context, *GetFunctionResourcesRequest) (*GetFunctionResourcesResponse, error)
	mustEmbedUnimplementedFaaSServer()
}

// UnimplementedFaaSServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedFaaSServer struct{}

func (UnimplementedFaaSServer) CallFunction(context.Context, *CallFunctionRequest) (*CallFunctionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CallFunction not implemented")
}
func (UnimplementedFaaSServer) CreateFunction(context.Context, *CreateFunctionRequest) (*CreateFunctionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateFunction not implemented")
}
func (UnimplementedFaaSServer) DeleteFunction(context.Context, *DeleteFunctionRequest) (*DeleteFunctionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteFunction not implemented")
}
func (UnimplementedFaaSServer) GetAllFunctions(context.Context, *Empty) (*GetAllFunctionsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllFunctions not implemented")
}
func (UnimplementedFaaSServer) GetFunctionResources(context.Context, *GetFunctionResourcesRequest) (*GetFunctionResourcesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFunctionResources not implemented")
}
func (UnimplementedFaaSServer) mustEmbedUnimplementedFaaSServer() {}
func (UnimplementedFaaSServer) testEmbeddedByValue()              {}

// UnsafeFaaSServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FaaSServer will
// result in compilation errors.
type UnsafeFaaSServer interface {
	mustEmbedUnimplementedFaaSServer()
}

func RegisterFaaSServer(s grpc.ServiceRegistrar, srv FaaSServer) {
	// If the following call pancis, it indicates UnimplementedFaaSServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&FaaS_ServiceDesc, srv)
}

func _FaaS_CallFunction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CallFunctionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FaaSServer).CallFunction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FaaS_CallFunction_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FaaSServer).CallFunction(ctx, req.(*CallFunctionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FaaS_CreateFunction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateFunctionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FaaSServer).CreateFunction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FaaS_CreateFunction_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FaaSServer).CreateFunction(ctx, req.(*CreateFunctionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FaaS_DeleteFunction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteFunctionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FaaSServer).DeleteFunction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FaaS_DeleteFunction_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FaaSServer).DeleteFunction(ctx, req.(*DeleteFunctionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FaaS_GetAllFunctions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FaaSServer).GetAllFunctions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FaaS_GetAllFunctions_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FaaSServer).GetAllFunctions(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _FaaS_GetFunctionResources_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetFunctionResourcesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FaaSServer).GetFunctionResources(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FaaS_GetFunctionResources_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FaaSServer).GetFunctionResources(ctx, req.(*GetFunctionResourcesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// FaaS_ServiceDesc is the grpc.ServiceDesc for FaaS service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var FaaS_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "faas.FaaS",
	HandlerType: (*FaaSServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CallFunction",
			Handler:    _FaaS_CallFunction_Handler,
		},
		{
			MethodName: "CreateFunction",
			Handler:    _FaaS_CreateFunction_Handler,
		},
		{
			MethodName: "DeleteFunction",
			Handler:    _FaaS_DeleteFunction_Handler,
		},
		{
			MethodName: "GetAllFunctions",
			Handler:    _FaaS_GetAllFunctions_Handler,
		},
		{
			MethodName: "GetFunctionResources",
			Handler:    _FaaS_GetFunctionResources_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/apiserver.proto",
}