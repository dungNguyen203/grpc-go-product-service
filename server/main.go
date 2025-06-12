package main

import (
	"context"
	"log"
	"net"
	pb "product-service/proto/gen"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func ServerAuthInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok || len(md["authorization"]) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "no auth token")
	}

	authToken := md["authorization"][0]
	if authToken != "unary-token" {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token")
	}

	return handler(ctx, req)
}

func ServerStreamAuthInterceptor(
	srv interface{},
	ss grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	// Extract metadata from stream context
	md, ok := metadata.FromIncomingContext(ss.Context())
	if !ok || len(md["authorization"]) == 0 {
		return status.Errorf(codes.Unauthenticated, "no auth token")
	}

	// Validate the authorization token
	authToken := md["authorization"][0]
	if authToken != "stream-token" {
		return status.Errorf(codes.Unauthenticated, "invalid token")
	}

	// Continue to the handler if authenticated
	return handler(srv, ss)
}

func main() {
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(ServerAuthInterceptor), grpc.StreamInterceptor(ServerStreamAuthInterceptor))
	
    pb.RegisterProductServiceServer(grpcServer, NewServer())

	listener, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	log.Println("Server is running on port :50052")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}