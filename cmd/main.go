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
	"github.com/Nap20192/shipment/pkg/logger"
	"golang.org/x/sync/errgroup"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	logInstance, err := logger.InitLogger(cfg.LogLevel, true, cfg.LogDir)
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	slog.SetDefault(logInstance)
	slog.Info("Starting Shipment Tracking Service", "port", cfg.GRPCPort)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dependencies, err := deps.NewDeps(
		ctx, deps.WithRepository(cfg.DBConnString()),
		deps.WithEventBus(),
		deps.WithShipmentService(),
		deps.WithGrpcServer(cfg.GrpcPortString()),
	)
	if err != nil {
		slog.Error("Failed to initialize dependencies", "error", err)
		os.Exit(1)
	}
	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return dependencies.Server.Serve()
	})
	g.Go(func() error {
		logSubscriber := app.NewLogSubscriber()
		slog.Info("Subscribing to events with LogSubscriber")
		dependencies.EventBus.Subscribe(app.EventBusKey, logSubscriber)
		return nil
	})
	g.Go(func() error {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		defer signal.Stop(c)
		select {
		case <-gCtx.Done():
			return nil
		case sig := <-c:
			slog.Info("Received signal, shutting down", "signal", sig)
			cancel()
			return nil
		}
	})

	g.Go(func() error {
		<-gCtx.Done()
		dependencies.Server.GracefulStop()
		dependencies.Pool.Close()
		return nil
	})

	if err := g.Wait(); err != nil {
		slog.Error("Server error", "error", err)
	}
	slog.Info("Shutting down Shipment Tracking Service")
}
