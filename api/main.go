package main

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	handler "Assignment_Golang/user/handler"
	pb "Assignment_Golang/user/proto/user_service/v1"
	repository "Assignment_Golang/user/repository"
	service "Assignment_Golang/user/service"
	handlerwallet "Assignment_Golang/wallet/handler"
	pbwallet "Assignment_Golang/wallet/proto/wallet_service/v1"
	repositoryWallet "Assignment_Golang/wallet/repository"
	serviceWallet "Assignment_Golang/wallet/service"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	dsn := "postgresql://postgres:postgres@localhost:5432/Training_Golang"
	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{SkipDefaultTransaction: true})
	if err != nil {
		log.Fatalln(err)
	}

	//setup service
	userRepo := repository.NewUserRepository(gormDB)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	//setup routernya
	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, userHandler)

	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	go func() {
		log.Println("Server is running on port 50051")
		if err = grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	//run the grpc gateway
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln("failed to dial server:", err)
	}

	//setup service wallet
	userServiceClient := pb.NewUserServiceClient(conn)
	walletRepo := repositoryWallet.NewWalletRepository(gormDB)
	walletService := serviceWallet.NewWalletService(walletRepo)
	walletHandler := handlerwallet.NewWalletHandler(walletService, userServiceClient)
	//setup router Waletnya
	grpcServerWallet := grpc.NewServer()
	pbwallet.RegisterWalletServiceServer(grpcServerWallet, walletHandler)

	liswallet, err := net.Listen("tcp", "localhost:50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	go func() {
		log.Println("Server is running on port 50052")
		if err = grpcServerWallet.Serve(liswallet); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	//run the grpc gateway
	connwalet, err := grpc.NewClient("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln("failed to dial server:", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	gwmux := runtime.NewServeMux()
	if err = pb.RegisterUserServiceHandler(ctx, gwmux, conn); err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	if err = pbwallet.RegisterWalletServiceHandler(ctx, gwmux, connwalet); err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	ginserver := gin.Default()

	ginserver.Group("/v1/*{grpc_gateway}").Any("", gin.WrapH(gwmux))
	log.Println("running grpc gateway server in port: 8080")
	ginserver.Run("localhost:8080")
}
