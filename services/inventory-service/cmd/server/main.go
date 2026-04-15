package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	// 1. Correct import for pgx v5
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"

	pb "github.com/durgaprasadraju/AeroStore/api/proto"

	"github.com/durgaprasadraju/AeroStore/services/inventory-service/config"
	"github.com/durgaprasadraju/AeroStore/services/inventory-service/internal/repository"
	"github.com/durgaprasadraju/AeroStore/services/inventory-service/internal/usecase"

	// ADD THE NICKNAME "grpcDelivery" HERE!
	grpcDelivery "github.com/durgaprasadraju/AeroStore/services/inventory-service/internal/delivery/grpc"
)

func main() {
	cfg := config.LoadConfig()

	// 1. Connect to Database
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbPool, err := pgxpool.New(ctx, cfg.PostgresURL)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
	}
	defer dbPool.Close()

	if err := dbPool.Ping(ctx); err != nil {
		log.Fatalf("Could not ping DB: %v\n", err)
	}
	fmt.Println("✅ Successfully connected to PostgreSQL!")

	// 2. Initialize Clean Architecture Layers
	repo := repository.NewPostgresRepository(dbPool)
	inventoryUC := usecase.NewInventoryUseCase(repo)
	handler := grpcDelivery.NewInventoryHandler(inventoryUC)

	// 3. Start the gRPC TCP Listener
	lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", cfg.GRPCPort, err)
	}

	// 4. Create the gRPC Server and Register our Handler
	grpcServer := grpc.NewServer()
	pb.RegisterInventoryServiceServer(grpcServer, handler)

	fmt.Printf("🚀 gRPC Server is actively listening on port %s...\n", cfg.GRPCPort)

	// 5. Serve requests (This blocks forever and keeps the app running)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
