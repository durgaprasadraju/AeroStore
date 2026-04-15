package main

import (
	"context"
	"fmt"
	"log"
	"time"

	// 1. Correct import for pgx v5
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/durgaprasadraju/AeroStore/services/inventory-service/config"
	"github.com/durgaprasadraju/AeroStore/services/inventory-service/internal/repository"
	"github.com/durgaprasadraju/AeroStore/services/inventory-service/internal/usecase"
)

func main() {
	cfg := config.LoadConfig()

	// 2. Setup a context with timeout for the initial connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 3. Connect using pgxpool
	// Note: cfg.PostgresURL should look like: "postgres://user:pass@localhost:5432/dbname"
	dbPool, err := pgxpool.New(ctx, cfg.PostgresURL)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
	}
	defer dbPool.Close()

	// 4. Ping the database to ensure connection is actually alive
	err = dbPool.Ping(ctx)
	if err != nil {
		log.Fatalf("Could not ping DB: %v\n", err)
	}

	fmt.Println("✅ Successfully connected to PostgreSQL!")

	// 5. Initialize Layers
	repo := repository.NewPostgresRepository(dbPool)
	inventoryUC := usecase.NewInventoryUseCase(repo)

	fmt.Println("🚀 Inventory Service Layers Ready...", inventoryUC)

	// Keep the service running (we will add gRPC server here later)
	select {}
}
