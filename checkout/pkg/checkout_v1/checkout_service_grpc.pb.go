// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.12.4
// source: checkout_service.proto

package checkout_v1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	Checkout_AddToCart_FullMethodName      = "/checkout.Checkout/AddToCart"
	Checkout_DeleteFromCart_FullMethodName = "/checkout.Checkout/DeleteFromCart"
	Checkout_ListCart_FullMethodName       = "/checkout.Checkout/ListCart"
	Checkout_Purchase_FullMethodName       = "/checkout.Checkout/Purchase"
)

// CheckoutClient is the client API for Checkout service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CheckoutClient interface {
	AddToCart(ctx context.Context, in *AddToCartRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	DeleteFromCart(ctx context.Context, in *DeleteFromCartRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	ListCart(ctx context.Context, in *ListCartRequest, opts ...grpc.CallOption) (*ListCartResponse, error)
	Purchase(ctx context.Context, in *PurchaseRequest, opts ...grpc.CallOption) (*PurchaseResponse, error)
}

type checkoutClient struct {
	cc grpc.ClientConnInterface
}

func NewCheckoutClient(cc grpc.ClientConnInterface) CheckoutClient {
	return &checkoutClient{cc}
}

func (c *checkoutClient) AddToCart(ctx context.Context, in *AddToCartRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, Checkout_AddToCart_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *checkoutClient) DeleteFromCart(ctx context.Context, in *DeleteFromCartRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, Checkout_DeleteFromCart_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *checkoutClient) ListCart(ctx context.Context, in *ListCartRequest, opts ...grpc.CallOption) (*ListCartResponse, error) {
	out := new(ListCartResponse)
	err := c.cc.Invoke(ctx, Checkout_ListCart_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *checkoutClient) Purchase(ctx context.Context, in *PurchaseRequest, opts ...grpc.CallOption) (*PurchaseResponse, error) {
	out := new(PurchaseResponse)
	err := c.cc.Invoke(ctx, Checkout_Purchase_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CheckoutServer is the server API for Checkout service.
// All implementations must embed UnimplementedCheckoutServer
// for forward compatibility
type CheckoutServer interface {
	AddToCart(context.Context, *AddToCartRequest) (*emptypb.Empty, error)
	DeleteFromCart(context.Context, *DeleteFromCartRequest) (*emptypb.Empty, error)
	ListCart(context.Context, *ListCartRequest) (*ListCartResponse, error)
	Purchase(context.Context, *PurchaseRequest) (*PurchaseResponse, error)
	mustEmbedUnimplementedCheckoutServer()
}

// UnimplementedCheckoutServer must be embedded to have forward compatible implementations.
type UnimplementedCheckoutServer struct {
}

func (UnimplementedCheckoutServer) AddToCart(context.Context, *AddToCartRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddToCart not implemented")
}
func (UnimplementedCheckoutServer) DeleteFromCart(context.Context, *DeleteFromCartRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteFromCart not implemented")
}
func (UnimplementedCheckoutServer) ListCart(context.Context, *ListCartRequest) (*ListCartResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListCart not implemented")
}
func (UnimplementedCheckoutServer) Purchase(context.Context, *PurchaseRequest) (*PurchaseResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Purchase not implemented")
}
func (UnimplementedCheckoutServer) mustEmbedUnimplementedCheckoutServer() {}

// UnsafeCheckoutServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CheckoutServer will
// result in compilation errors.
type UnsafeCheckoutServer interface {
	mustEmbedUnimplementedCheckoutServer()
}

func RegisterCheckoutServer(s grpc.ServiceRegistrar, srv CheckoutServer) {
	s.RegisterService(&Checkout_ServiceDesc, srv)
}

func _Checkout_AddToCart_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddToCartRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CheckoutServer).AddToCart(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Checkout_AddToCart_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CheckoutServer).AddToCart(ctx, req.(*AddToCartRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Checkout_DeleteFromCart_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteFromCartRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CheckoutServer).DeleteFromCart(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Checkout_DeleteFromCart_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CheckoutServer).DeleteFromCart(ctx, req.(*DeleteFromCartRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Checkout_ListCart_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListCartRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CheckoutServer).ListCart(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Checkout_ListCart_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CheckoutServer).ListCart(ctx, req.(*ListCartRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Checkout_Purchase_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PurchaseRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CheckoutServer).Purchase(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Checkout_Purchase_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CheckoutServer).Purchase(ctx, req.(*PurchaseRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Checkout_ServiceDesc is the grpc.ServiceDesc for Checkout service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Checkout_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "checkout.Checkout",
	HandlerType: (*CheckoutServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddToCart",
			Handler:    _Checkout_AddToCart_Handler,
		},
		{
			MethodName: "DeleteFromCart",
			Handler:    _Checkout_DeleteFromCart_Handler,
		},
		{
			MethodName: "ListCart",
			Handler:    _Checkout_ListCart_Handler,
		},
		{
			MethodName: "Purchase",
			Handler:    _Checkout_Purchase_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "checkout_service.proto",
}
