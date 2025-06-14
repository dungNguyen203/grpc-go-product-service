// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.3
// source: product.proto

// Defines the protocol buffer's package as productpb. This helps organize and prevent naming conflicts in large projects.

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	ProductService_CreateProduct_FullMethodName  = "/productpb.ProductService/CreateProduct"
	ProductService_GetProduct_FullMethodName     = "/productpb.ProductService/GetProduct"
	ProductService_GetAllProducts_FullMethodName = "/productpb.ProductService/GetAllProducts"
	ProductService_ListProducts_FullMethodName   = "/productpb.ProductService/ListProducts"
	ProductService_UpdateProduct_FullMethodName  = "/productpb.ProductService/UpdateProduct"
)

// ProductServiceClient is the client API for ProductService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ProductServiceClient interface {
	// This case we don't want to expose CreateProduct to gRPC-Gateway, so it can only be called by gRPC common method
	CreateProduct(ctx context.Context, in *ProductRequest, opts ...grpc.CallOption) (*ProductResponse, error)
	GetProduct(ctx context.Context, in *ProductID, opts ...grpc.CallOption) (*ProductResponse, error)
	GetAllProducts(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ProductList, error)
	ListProducts(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (grpc.ServerStreamingClient[Product], error)
	UpdateProduct(ctx context.Context, in *ProductUdpateRequest, opts ...grpc.CallOption) (*ProductResponse, error)
}

type productServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewProductServiceClient(cc grpc.ClientConnInterface) ProductServiceClient {
	return &productServiceClient{cc}
}

func (c *productServiceClient) CreateProduct(ctx context.Context, in *ProductRequest, opts ...grpc.CallOption) (*ProductResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ProductResponse)
	err := c.cc.Invoke(ctx, ProductService_CreateProduct_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *productServiceClient) GetProduct(ctx context.Context, in *ProductID, opts ...grpc.CallOption) (*ProductResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ProductResponse)
	err := c.cc.Invoke(ctx, ProductService_GetProduct_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *productServiceClient) GetAllProducts(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ProductList, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ProductList)
	err := c.cc.Invoke(ctx, ProductService_GetAllProducts_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *productServiceClient) ListProducts(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (grpc.ServerStreamingClient[Product], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &ProductService_ServiceDesc.Streams[0], ProductService_ListProducts_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[emptypb.Empty, Product]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type ProductService_ListProductsClient = grpc.ServerStreamingClient[Product]

func (c *productServiceClient) UpdateProduct(ctx context.Context, in *ProductUdpateRequest, opts ...grpc.CallOption) (*ProductResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ProductResponse)
	err := c.cc.Invoke(ctx, ProductService_UpdateProduct_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ProductServiceServer is the server API for ProductService service.
// All implementations should embed UnimplementedProductServiceServer
// for forward compatibility.
type ProductServiceServer interface {
	// This case we don't want to expose CreateProduct to gRPC-Gateway, so it can only be called by gRPC common method
	CreateProduct(context.Context, *ProductRequest) (*ProductResponse, error)
	GetProduct(context.Context, *ProductID) (*ProductResponse, error)
	GetAllProducts(context.Context, *emptypb.Empty) (*ProductList, error)
	ListProducts(*emptypb.Empty, grpc.ServerStreamingServer[Product]) error
	UpdateProduct(context.Context, *ProductUdpateRequest) (*ProductResponse, error)
}

// UnimplementedProductServiceServer should be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedProductServiceServer struct{}

func (UnimplementedProductServiceServer) CreateProduct(context.Context, *ProductRequest) (*ProductResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateProduct not implemented")
}
func (UnimplementedProductServiceServer) GetProduct(context.Context, *ProductID) (*ProductResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetProduct not implemented")
}
func (UnimplementedProductServiceServer) GetAllProducts(context.Context, *emptypb.Empty) (*ProductList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllProducts not implemented")
}
func (UnimplementedProductServiceServer) ListProducts(*emptypb.Empty, grpc.ServerStreamingServer[Product]) error {
	return status.Errorf(codes.Unimplemented, "method ListProducts not implemented")
}
func (UnimplementedProductServiceServer) UpdateProduct(context.Context, *ProductUdpateRequest) (*ProductResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateProduct not implemented")
}
func (UnimplementedProductServiceServer) testEmbeddedByValue() {}

// UnsafeProductServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ProductServiceServer will
// result in compilation errors.
type UnsafeProductServiceServer interface {
	mustEmbedUnimplementedProductServiceServer()
}

func RegisterProductServiceServer(s grpc.ServiceRegistrar, srv ProductServiceServer) {
	// If the following call pancis, it indicates UnimplementedProductServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&ProductService_ServiceDesc, srv)
}

func _ProductService_CreateProduct_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProductRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProductServiceServer).CreateProduct(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProductService_CreateProduct_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProductServiceServer).CreateProduct(ctx, req.(*ProductRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProductService_GetProduct_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProductID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProductServiceServer).GetProduct(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProductService_GetProduct_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProductServiceServer).GetProduct(ctx, req.(*ProductID))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProductService_GetAllProducts_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProductServiceServer).GetAllProducts(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProductService_GetAllProducts_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProductServiceServer).GetAllProducts(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProductService_ListProducts_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(emptypb.Empty)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ProductServiceServer).ListProducts(m, &grpc.GenericServerStream[emptypb.Empty, Product]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type ProductService_ListProductsServer = grpc.ServerStreamingServer[Product]

func _ProductService_UpdateProduct_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProductUdpateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProductServiceServer).UpdateProduct(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProductService_UpdateProduct_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProductServiceServer).UpdateProduct(ctx, req.(*ProductUdpateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ProductService_ServiceDesc is the grpc.ServiceDesc for ProductService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ProductService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "productpb.ProductService",
	HandlerType: (*ProductServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateProduct",
			Handler:    _ProductService_CreateProduct_Handler,
		},
		{
			MethodName: "GetProduct",
			Handler:    _ProductService_GetProduct_Handler,
		},
		{
			MethodName: "GetAllProducts",
			Handler:    _ProductService_GetAllProducts_Handler,
		},
		{
			MethodName: "UpdateProduct",
			Handler:    _ProductService_UpdateProduct_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ListProducts",
			Handler:       _ProductService_ListProducts_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "product.proto",
}
