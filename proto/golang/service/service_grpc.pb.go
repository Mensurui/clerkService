// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.12.4
// source: service/service.proto

package service

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
	Service_GetList_FullMethodName = "/service.Service/GetList"
)

// ServiceClient is the client API for Service service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ServiceClient interface {
	GetList(ctx context.Context, in *GetListRequest, opts ...grpc.CallOption) (*GetListResponse, error)
}

type serviceClient struct {
	cc grpc.ClientConnInterface
}

func NewServiceClient(cc grpc.ClientConnInterface) ServiceClient {
	return &serviceClient{cc}
}

func (c *serviceClient) GetList(ctx context.Context, in *GetListRequest, opts ...grpc.CallOption) (*GetListResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetListResponse)
	err := c.cc.Invoke(ctx, Service_GetList_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ServiceServer is the server API for Service service.
// All implementations must embed UnimplementedServiceServer
// for forward compatibility.
type ServiceServer interface {
	GetList(context.Context, *GetListRequest) (*GetListResponse, error)
	mustEmbedUnimplementedServiceServer()
}

// UnimplementedServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedServiceServer struct{}

func (UnimplementedServiceServer) GetList(context.Context, *GetListRequest) (*GetListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetList not implemented")
}
func (UnimplementedServiceServer) mustEmbedUnimplementedServiceServer() {}
func (UnimplementedServiceServer) testEmbeddedByValue()                 {}

// UnsafeServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ServiceServer will
// result in compilation errors.
type UnsafeServiceServer interface {
	mustEmbedUnimplementedServiceServer()
}

func RegisterServiceServer(s grpc.ServiceRegistrar, srv ServiceServer) {
	// If the following call pancis, it indicates UnimplementedServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Service_ServiceDesc, srv)
}

func _Service_GetList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).GetList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Service_GetList_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).GetList(ctx, req.(*GetListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Service_ServiceDesc is the grpc.ServiceDesc for Service service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Service_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "service.Service",
	HandlerType: (*ServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetList",
			Handler:    _Service_GetList_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "service/service.proto",
}
