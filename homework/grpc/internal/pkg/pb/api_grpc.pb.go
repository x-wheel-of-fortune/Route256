// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.12.4
// source: api.proto

package pb

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
	PickupPoints_AddPickupPoint_FullMethodName    = "/pickup_point.PickupPoints/AddPickupPoint"
	PickupPoints_UpdatePickupPoint_FullMethodName = "/pickup_point.PickupPoints/UpdatePickupPoint"
	PickupPoints_GetPickupPoint_FullMethodName    = "/pickup_point.PickupPoints/GetPickupPoint"
	PickupPoints_DeletePickupPoint_FullMethodName = "/pickup_point.PickupPoints/DeletePickupPoint"
	PickupPoints_ListPickupPoint_FullMethodName   = "/pickup_point.PickupPoints/ListPickupPoint"
)

// PickupPointsClient is the client API for PickupPoints service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PickupPointsClient interface {
	AddPickupPoint(ctx context.Context, in *PickupPointRequest, opts ...grpc.CallOption) (*PickupPointResponse, error)
	UpdatePickupPoint(ctx context.Context, in *PickupPointRequest, opts ...grpc.CallOption) (*PickupPointResponse, error)
	GetPickupPoint(ctx context.Context, in *IdRequest, opts ...grpc.CallOption) (*PickupPointResponse, error)
	DeletePickupPoint(ctx context.Context, in *IdRequest, opts ...grpc.CallOption) (*Empty, error)
	ListPickupPoint(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*ListPickupPointResponse, error)
}

type pickupPointsClient struct {
	cc grpc.ClientConnInterface
}

func NewPickupPointsClient(cc grpc.ClientConnInterface) PickupPointsClient {
	return &pickupPointsClient{cc}
}

func (c *pickupPointsClient) AddPickupPoint(ctx context.Context, in *PickupPointRequest, opts ...grpc.CallOption) (*PickupPointResponse, error) {
	out := new(PickupPointResponse)
	err := c.cc.Invoke(ctx, PickupPoints_AddPickupPoint_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pickupPointsClient) UpdatePickupPoint(ctx context.Context, in *PickupPointRequest, opts ...grpc.CallOption) (*PickupPointResponse, error) {
	out := new(PickupPointResponse)
	err := c.cc.Invoke(ctx, PickupPoints_UpdatePickupPoint_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pickupPointsClient) GetPickupPoint(ctx context.Context, in *IdRequest, opts ...grpc.CallOption) (*PickupPointResponse, error) {
	out := new(PickupPointResponse)
	err := c.cc.Invoke(ctx, PickupPoints_GetPickupPoint_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pickupPointsClient) DeletePickupPoint(ctx context.Context, in *IdRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, PickupPoints_DeletePickupPoint_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pickupPointsClient) ListPickupPoint(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*ListPickupPointResponse, error) {
	out := new(ListPickupPointResponse)
	err := c.cc.Invoke(ctx, PickupPoints_ListPickupPoint_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PickupPointsServer is the server API for PickupPoints service.
// All implementations must embed UnimplementedPickupPointsServer
// for forward compatibility
type PickupPointsServer interface {
	AddPickupPoint(context.Context, *PickupPointRequest) (*PickupPointResponse, error)
	UpdatePickupPoint(context.Context, *PickupPointRequest) (*PickupPointResponse, error)
	GetPickupPoint(context.Context, *IdRequest) (*PickupPointResponse, error)
	DeletePickupPoint(context.Context, *IdRequest) (*Empty, error)
	ListPickupPoint(context.Context, *Empty) (*ListPickupPointResponse, error)
	mustEmbedUnimplementedPickupPointsServer()
}

// UnimplementedPickupPointsServer must be embedded to have forward compatible implementations.
type UnimplementedPickupPointsServer struct {
}

func (UnimplementedPickupPointsServer) AddPickupPoint(context.Context, *PickupPointRequest) (*PickupPointResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddPickupPoint not implemented")
}
func (UnimplementedPickupPointsServer) UpdatePickupPoint(context.Context, *PickupPointRequest) (*PickupPointResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdatePickupPoint not implemented")
}
func (UnimplementedPickupPointsServer) GetPickupPoint(context.Context, *IdRequest) (*PickupPointResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPickupPoint not implemented")
}
func (UnimplementedPickupPointsServer) DeletePickupPoint(context.Context, *IdRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeletePickupPoint not implemented")
}
func (UnimplementedPickupPointsServer) ListPickupPoint(context.Context, *Empty) (*ListPickupPointResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListPickupPoint not implemented")
}
func (UnimplementedPickupPointsServer) mustEmbedUnimplementedPickupPointsServer() {}

// UnsafePickupPointsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PickupPointsServer will
// result in compilation errors.
type UnsafePickupPointsServer interface {
	mustEmbedUnimplementedPickupPointsServer()
}

func RegisterPickupPointsServer(s grpc.ServiceRegistrar, srv PickupPointsServer) {
	s.RegisterService(&PickupPoints_ServiceDesc, srv)
}

func _PickupPoints_AddPickupPoint_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PickupPointRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PickupPointsServer).AddPickupPoint(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PickupPoints_AddPickupPoint_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PickupPointsServer).AddPickupPoint(ctx, req.(*PickupPointRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PickupPoints_UpdatePickupPoint_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PickupPointRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PickupPointsServer).UpdatePickupPoint(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PickupPoints_UpdatePickupPoint_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PickupPointsServer).UpdatePickupPoint(ctx, req.(*PickupPointRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PickupPoints_GetPickupPoint_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PickupPointsServer).GetPickupPoint(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PickupPoints_GetPickupPoint_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PickupPointsServer).GetPickupPoint(ctx, req.(*IdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PickupPoints_DeletePickupPoint_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PickupPointsServer).DeletePickupPoint(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PickupPoints_DeletePickupPoint_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PickupPointsServer).DeletePickupPoint(ctx, req.(*IdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PickupPoints_ListPickupPoint_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PickupPointsServer).ListPickupPoint(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PickupPoints_ListPickupPoint_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PickupPointsServer).ListPickupPoint(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// PickupPoints_ServiceDesc is the grpc.ServiceDesc for PickupPoints service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PickupPoints_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pickup_point.PickupPoints",
	HandlerType: (*PickupPointsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddPickupPoint",
			Handler:    _PickupPoints_AddPickupPoint_Handler,
		},
		{
			MethodName: "UpdatePickupPoint",
			Handler:    _PickupPoints_UpdatePickupPoint_Handler,
		},
		{
			MethodName: "GetPickupPoint",
			Handler:    _PickupPoints_GetPickupPoint_Handler,
		},
		{
			MethodName: "DeletePickupPoint",
			Handler:    _PickupPoints_DeletePickupPoint_Handler,
		},
		{
			MethodName: "ListPickupPoint",
			Handler:    _PickupPoints_ListPickupPoint_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api.proto",
}
