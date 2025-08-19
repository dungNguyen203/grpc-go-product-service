package main

import (
	"context"
	"net"
	"net/http"
	"os"
	db "product-service/db/sqlc"
	"product-service/mail"
	"product-service/proto/pb"
	api "product-service/server/apis"
	"product-service/server/gapi"
	"product-service/utils"
	"product-service/worker"

	"github.com/golang-migrate/migrate/v4"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
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
	config, err := utils.LoadConfig("../")
	if err != nil {
		log.Fatal().Msg("cannot load config: " + err.Error())
	}

	if config.Env == "dev" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	connPool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal().Msg("cannot connect to db")
	}

	// runDBMigration(config.MigrationURL, config.DBSource)

	store := db.NewStore(connPool)

	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}

	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)

	// go runTaskProcessor(config, redisOpt, store)
	go runGatewayServer(config, store, taskDistributor)
	go runGinServer(config, store)
	runGrpcServer(config, store, taskDistributor)
}

func runTaskProcessor(config utils.Config, redisOpt asynq.RedisClientOpt, store db.Store) {
	mailer := mail.NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, store, mailer)
	log.Info().Msg("start task processor")
	err := taskProcessor.Start()
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start task processor")
	}
}

func runGrpcServer(config utils.Config, store db.Store, taskDistributor worker.TaskDistributor) {
	server, err := gapi.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Fatal().Msg("cannot create server: " + err.Error())
	}

	grpcLogger := grpc.UnaryInterceptor(gapi.GrpcLogger)
	grpcServer := grpc.NewServer(grpcLogger)
	pb.RegisterTinyBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal().Msg("cannot listen to port")
	}

	log.Printf("start gRPC server is listening at %s", config.GRPCServerAddress)

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal().Msg("cannot start gRPC server")
	}

}

func runGatewayServer(config utils.Config, store db.Store, taskDistributor worker.TaskDistributor) {
	server, err := gapi.NewServer(config, store, taskDistributor);
	if err != nil {
		log.Fatal().Msg("cannot create server: " + err.Error())
	}

	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})
	grpcMux := runtime.NewServeMux(jsonOption)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = pb.RegisterTinyBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal().Msg("cannot register handler server")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	// statikFS, err := fs.New()
	// if err != nil {
	// 	log.Fatal().Msg("cannot create statik fs : " + err.Error())
	// }

	// swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFS))
	// mux.Handle("/swagger/", swaggerHandler)

	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatal().Msg("cannot listen to port")
	}

	log.Printf("start HTTP server is listening at %s", config.HTTPServerAddress)

	handler := gapi.HttpLogger(mux)
	err = http.Serve(listener, handler)
	if err != nil {
		log.Fatal().Msg("cannot start HTTP gateway server")
	}

}


func runGinServer(config utils.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal().Msg("cannot create server : " + err.Error())
	}

	log.Printf("start Gin server is listening at %s", config.GINServerAddress)
	err = server.Start(config.GINServerAddress)
	if err != nil {
		log.Fatal().Msg("cannot start server")
	}
}

func runDBMigration(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)

	if err != nil {
		log.Fatal().Msg("cannot create migration")
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal().Msg("cannot run migration")
	}

	log.Info().Msg("migration is done")
}