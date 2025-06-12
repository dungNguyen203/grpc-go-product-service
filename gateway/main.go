package main

import (
	"context"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"product-service/auth"
	pb "product-service/proto/gen"
)

func main() {
	mux := runtime.NewServeMux()
	err := pb.RegisterProductServiceHandlerFromEndpoint(context.Background(), mux, ":50052", []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithUnaryInterceptor(auth.AuthInterceptor("unary-token"))})
	if err != nil {
		log.Fatalf("Failed to start HTTP gateway: %v", err)
	}

	log.Println("HTTP Gateway running on :8080")
	http.ListenAndServe(":8080", mux)
    
}
