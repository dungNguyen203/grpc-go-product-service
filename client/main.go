package main

import (
	"context"
	"io"
	"log"
	"product-service/auth"
	"product-service/models"
	"time"

	pb "product-service/proto/gen"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type ProductClient struct {
	client pb.ProductServiceClient
}

func NewProductClient(cc *grpc.ClientConn) *ProductClient {
	return &ProductClient{client: pb.NewProductServiceClient(cc)}
}

// REST API handler functions
func (c *ProductClient) createProduct(ctx *gin.Context) {
	var product models.Product
	if err := ctx.ShouldBindJSON(&product); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid product data"})
		return
	}

	protoProduct := product.ToProto()
	req := &pb.ProductRequest{Product: protoProduct}
	res, err := c.client.CreateProduct(ctx, req)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(201, models.ProductFromProto(res.Product))
}

func (c *ProductClient) getAllProducts(ctx *gin.Context) {
	deadlineCtx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	res, err := c.client.GetAllProducts(deadlineCtx, &emptypb.Empty{})
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	var products []models.Product
	for _, protoProduct := range res.Products {
		products = append(products, *models.ProductFromProto(protoProduct))
	}

	ctx.JSON(200, products)
}

// Background job for streaming new products
func (c *ProductClient) StreamNewProducts() {
	ctx := context.Background()
	go func() {
		for {
			stream, err := c.client.ListProducts(ctx, &emptypb.Empty{})
			if err != nil {
				log.Printf("Error connecting to ListProducts: %v", err)
				time.Sleep(5 * time.Second) // Retry delay
				continue
			}

			for {
				product, err := stream.Recv()
				if err == io.EOF {
					// The EOF error in your StreamNewProducts function likely indicates that the server has closed the stream, often because there are no new products to send, and the stream reaches the end
					log.Println("Completed!, Stream closed by server.")
					break // Break inner loop to reconnect
				}
				if err != nil {
					log.Printf("Error receiving product: %v", err)
					break
				}
				log.Printf("New Product: %v", product)
			}

			// Optional reconnect delay
			time.Sleep(1 * time.Minute)
		}
	}()
}

func (p *ProductClient) UpdateProduct(ctx *gin.Context) {
    // Get product ID from path parameter
    id := ctx.Param("id")
    var product models.Product
    if err := ctx.ShouldBindJSON(&product); err != nil {
        ctx.JSON(400, gin.H{"error": "Invalid product data"})
        return
    }
	product.ID = id // Ensure the ID is set from the request
	// Update product in database
	protoProduct := product.ToProto()
	req := &pb.ProductUdpateRequest{Product: protoProduct}
	res, err := p.client.UpdateProduct(ctx, req)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, models.ProductFromProto(res.Product))

}

func setupRouter(pc *ProductClient) *gin.Engine {
	r := gin.Default()
	r.POST("/products", pc.createProduct)
	r.GET("/products", pc.getAllProducts)
    r.PUT("/products/:id", pc.UpdateProduct)
	return r
}

func main() {
	unaryToken := "unary-token"
	streamToken := "stream-token"
	// This approach keeps the authorization token consistent across all requests without manually adding it each time.
	conn, err := grpc.NewClient(":50052", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithUnaryInterceptor(auth.AuthInterceptor(unaryToken)), grpc.WithStreamInterceptor(auth.AuthStreamInterceptor(streamToken)))

	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer conn.Close()

	productClient := NewProductClient(conn)

	// Start background streaming of new products
	productClient.StreamNewProducts()

	// Setup Gin REST API server
	r := setupRouter(productClient)
	r.Run(":8081")
}
