package api

import (
	"fmt"
	db "product-service/db/sqlc"
	"product-service/utils"
	"product-service/utils/token"

	// "go-grpc-product-service/models"
	// pb "go-grpc-product-service/protocol/gen"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type server struct {
	db *gorm.DB
}

type Server struct {
	config     utils.Config
	store      db.Store
	tokenMaker token.Maker
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
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{config: config, store: store, tokenMaker: tokenMaker}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	// router.POST("/auth/refresh_token", server.renewAccessToken)
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)
	router.GET("/users", server.getUsers)

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

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
