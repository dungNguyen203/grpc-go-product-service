package api

import (
	"context"
	db "product-service/db/sqlc"
	"product-service/models"
	pb "product-service/proto/gen"
	"product-service/utils"

	// "go-grpc-product-service/models"
	// pb "go-grpc-product-service/protocol/gen"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

type server struct {
	db *gorm.DB
}

type Server struct {
	config     utils.Config
	store      db.Store
	// tokenMaker token.Maker
	router     *gin.Engine
    db *gorm.DB
}

// func NewServer() *Server {
// 	db, err := gorm.Open(sqlite.Open("products.db"), &gorm.Config{})
// 	if err != nil {
// 		log.Fatalf("Failed to connect database: %v", err)
// 	}
// 	db.AutoMigrate(&models.Product{})
// 	return &Server{db: db}
// }

func NewServer(config utils.Config, store db.Store) (*Server, error) {
	// tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	// if err != nil {
	// 	return nil, fmt.Errorf("cannot create token maker: %w", err)
	// }
	server := &Server{config: config, store: store}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.setupRouter()

	return server, nil
}

func (s *server) CreateProduct(ctx context.Context, req *pb.ProductRequest) (*pb.ProductResponse, error) {

	product := models.ProductFromProto(req.Product)
	product.ID = uuid.New().String()
	if err := s.db.Create(&product).Error; err != nil {
		return nil, err
	}
	return &pb.ProductResponse{Product: product.ToProto()}, nil
}

func (s *server) GetProduct(ctx context.Context, req *pb.ProductID) (*pb.ProductResponse, error) {

	var product models.Product
	if err := s.db.First(&product, "id = ?", req.Id).Error; err != nil {
		return nil, err
	}
	return &pb.ProductResponse{Product: product.ToProto()}, nil
}

func (s *server) GetAllProducts(ctx context.Context, req *emptypb.Empty) (*pb.ProductList, error) {
	var products []models.Product
	if err := s.db.Find(&products).Error; err != nil {
		return nil, err
	}

	var productList []*pb.Product
	for _, product := range products {
		productList = append(productList, product.ToProto())
	}

	return &pb.ProductList{Products: productList}, nil
}

// Streaming method to list products
func (s *server) ListProducts(req *emptypb.Empty, stream pb.ProductService_ListProductsServer) error {
	var products []models.Product
	if err := s.db.Find(&products).Error; err != nil {
		return err
	}
	for _, product := range products {
		if err := stream.Send(product.ToProto()); err != nil {
			return err
		}
	}
	return nil
}

func (s *server) UpdateProduct(ctx context.Context, req *pb.ProductUdpateRequest) (*pb.ProductResponse, error) {
	product := models.ProductFromProto(req.Product)
	product.ID = req.Product.Id // Ensure the ID is set from the request
	if err := s.db.Save(&product).Error; err != nil {
		return nil, err
	}
	return &pb.ProductResponse{Product: product.ToProto()}, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	// router.POST("/auth/refresh_token", server.renewAccessToken)

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	// authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	// authRoutes.POST("/accounts", server.createAccount)
	// authRoutes.GET("/accounts/:id", server.getAccount)
	// authRoutes.GET("/accounts", server.listAccount)

	// authRoutes.POST("/transfers", server.createTransfer)

	server.router = router
}

// Additional CRUD methods...

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
