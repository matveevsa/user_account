package main

import (
	_ "account/docs"
	"account/internal/config"
	"account/internal/logger"
	"account/internal/repository"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

// @title Account Service
// @version 1.0
// @description Account Service
// @host localhost:8080
// @BasePath /
func main() {
	cfg, err := config.Load()

	if err != nil {
		log.Fatal("Failed to load config: %v", err)
	}

	lg := logger.New(cfg)

	db, err := gorm.Open(postgres.Open(cfg.DbDsn))

	if err != nil {
		lg.Error().Msgf("Failed to connect to databse: %v", err)
		return
	}
	lg.Info().Msg("database connected")

	repo := repository.NewRepository(db, lg)
	_ = repo

	listenAddr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	lis, err := net.Listen("tcp", listenAddr)

	if err != nil {
		lg.Error().Msgf("Failed ti listen on %s: %v", listenAddr, err)
		return
	}

	grpcServer := grpc.NewServer()

	healthSrv := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthSrv)
	reflection.Register(grpcServer)

	lg.Info().Msgf("gRPC server listening on %s", listenAddr)
	if err := grpcServer.Serve(lis); err != nil {
		lg.Error().Msgf("Failed to serve gRPC: %v", err)
		return
	}

	lg.Info().Msg("service start up")
}

// PingExample godoc
// @Summary Healthcheck service
// @Description Return pong
// @Tags health
// @Success 200 {string} string "pong"
// @Router /ping [get]
func PingExample(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}
