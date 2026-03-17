package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Nap20192/shipment/internal/config"
	"github.com/Nap20192/shipment/internal/core/app"
	"github.com/Nap20192/shipment/internal/deps"
	"github.com/Nap20192/shipment/internal/infra"
	"github.com/Nap20192/shipment/internal/pkg/sqlc"
	"github.com/Nap20192/shipment/internal/presentation/grpc"
	"github.com/Nap20192/shipment/pkg/logger"
)

func main() {
	// 1. Load Config
	cfg := config.LoadConfig()

	// 2. Initialize Logger
	logInstance, err := logger.InitLogger(cfg.LogLevel, true, cfg.LogDir)
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	slog.SetDefault(logInstance)
	slog.Info("Starting Shipment Tracking Service", "port", cfg.GRPCPort)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 3. Initialize Database Pool
	pool, err := infra.NewPgxPool(ctx, cfg.DBConnString())
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	// 4. Initialize Dependencies
	queries := sqlc.New(pool)
	eventBus := app.NewEventBus()

	d, err := deps.NewDeps(ctx, deps.WithShipmentService(queries, eventBus))
	if err != nil {
		slog.Error("Failed to initialize dependencies", "error", err)
		os.Exit(1)
	}

	// 5. Initialize gRPC Server
	addr := fmt.Sprintf(":%d", cfg.GRPCPort)
	server, err := grpc.NewServer(addr, d.AppService)
	if err != nil {
		slog.Error("Failed to create gRPC server", "error", err)
		os.Exit(1)
	}

	// 6. Start Serving
	go func() {
		if err := server.Serve(); err != nil {
			slog.Error("gRPC server failed", "error", err)
		}
	}()

	slog.Info("Service is running", "addr", addr)

	// 7. Graceful Shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	slog.Info("Shutting down service...")
	server.GracefulStop()
	slog.Info("Service stopped gracefully")
}
