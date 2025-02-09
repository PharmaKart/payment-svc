package main

import (
	"net"

	"github.com/PharmaKart/payment-svc/internal/handlers"
	"github.com/PharmaKart/payment-svc/internal/proto"
	"github.com/PharmaKart/payment-svc/internal/repositories"
	"github.com/PharmaKart/payment-svc/pkg/config"
	"github.com/PharmaKart/payment-svc/pkg/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/gorm"
)

func main() {
	// Initialize logger
	utils.InitLogger()

	// Load configurations
	cfg := config.LoadConfig()

	// Initialize repositories
	paymentRepo := repositories.NewPaymentRepository(&gorm.DB{})

	// Initialize order client
	conn, err := grpc.NewClient(cfg.OrderServiceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		utils.Logger.Fatal("Failed to connect to order service", map[string]interface{}{
			"error": err,
		})
	}

	orderClient := proto.NewOrderServiceClient(conn)
	defer conn.Close()

	// Initialize handlers
	paymentHandler := handlers.NewPaymentHandler(paymentRepo, &orderClient)

	// Initialize grpc server
	lis, err := net.Listen("tcp", ":"+cfg.Port)

	if err != nil {
		utils.Logger.Fatal("Failed to listen", map[string]interface{}{
			"error": err,
		})
	}

	grpcServer := grpc.NewServer()
	proto.RegisterPaymentServiceServer(grpcServer, paymentHandler)

	utils.Info("Starting payment service", map[string]interface{}{
		"port": cfg.Port,
	})

	if err := grpcServer.Serve(lis); err != nil {
		utils.Logger.Fatal("Failed to serve", map[string]interface{}{
			"error": err,
		})
	}
}
