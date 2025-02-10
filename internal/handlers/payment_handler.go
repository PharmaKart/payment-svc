package handlers

import (
	"context"
	"errors"

	"github.com/PharmaKart/payment-svc/internal/models"
	"github.com/PharmaKart/payment-svc/internal/proto"
	"github.com/PharmaKart/payment-svc/internal/repositories"
	"github.com/PharmaKart/payment-svc/internal/services"
	"github.com/PharmaKart/payment-svc/pkg/config"
	"github.com/google/uuid"
)

type PaymentHandler interface {
	GeneratePaymentURL(ctx context.Context, req *proto.GeneratePaymentURLRequest) (*proto.GeneratePaymentURLResponse, error)
	StorePayment(ctx context.Context, req *proto.StorePaymentRequest) (*proto.StorePaymentResponse, error)
	RefundPayment(ctx context.Context, req *proto.RefundPaymentRequest) (*proto.RefundPaymentResponse, error)
	GetPaymentByTransactionID(ctx context.Context, req *proto.GetPaymentByTransactionIDRequest) (*proto.GetPaymentResponse, error)
	GetPayment(ctx context.Context, req *proto.GetPaymentRequest) (*proto.GetPaymentResponse, error)
	GetPaymentByOrderID(ctx context.Context, req *proto.GetPaymentByOrderIDRequest) (*proto.GetPaymentResponse, error)
}

type paymentHandler struct {
	proto.UnimplementedPaymentServiceServer
	paymentService services.PaymentService
}

func NewPaymentHandler(paymentRepo repositories.PaymentRepository, orderClient *proto.OrderServiceClient, cfg *config.Config) *paymentHandler {
	return &paymentHandler{
		paymentService: services.NewPaymentService(paymentRepo, orderClient, cfg),
	}
}

func (h *paymentHandler) GeneratePaymentURL(ctx context.Context, req *proto.GeneratePaymentURLRequest) (*proto.GeneratePaymentURLResponse, error) {
	resp, err := h.paymentService.GeneratePaymentURL(req.OrderId, req.CustomerId)
	if err != nil {
		return nil, err
	}

	return &proto.GeneratePaymentURLResponse{
		Url: resp.URL,
	}, nil
}

func (h *paymentHandler) StorePayment(ctx context.Context, req *proto.StorePaymentRequest) (*proto.StorePaymentResponse, error) {
	orderId, err := uuid.Parse(req.OrderId)
	customerId, err := uuid.Parse(req.CustomerId)
	payment := &models.Payment{
		OrderID:       orderId,
		CustomerID:    customerId,
		TransactionID: req.TransactionId,
		Amount:        req.Amount,
		Status:        req.Status,
	}
	message, err := h.paymentService.StorePayment(payment)
	if err != nil {
		return nil, err
	}

	return &proto.StorePaymentResponse{Message: message}, nil
}

func (h *paymentHandler) RefundPayment(ctx context.Context, req *proto.RefundPaymentRequest) (*proto.RefundPaymentResponse, error) {
	err := h.paymentService.RefundPayment(req.TransactionId)
	if err != nil {
		return nil, err
	}

	return &proto.RefundPaymentResponse{}, nil
}

func (h *paymentHandler) GetPaymentByTransactionID(ctx context.Context, req *proto.GetPaymentByTransactionIDRequest) (*proto.GetPaymentResponse, error) {
	payment, err := h.paymentService.GetPaymentByTransactionID(req.TransactionId)
	if err != nil {
		return nil, err
	}

	customerId := req.CustomerId

	if customerId != "admin" && payment.CustomerID.String() != customerId {
		return nil, errors.New("Access denied")
	}

	return &proto.GetPaymentResponse{
		PaymentId:     payment.ID.String(),
		OrderId:       payment.OrderID.String(),
		CustomerId:    payment.CustomerID.String(),
		TransactionId: payment.TransactionID,
		Amount:        payment.Amount,
		Status:        payment.Status,
	}, nil
}

func (h *paymentHandler) GetPayment(ctx context.Context, req *proto.GetPaymentRequest) (*proto.GetPaymentResponse, error) {
	payment, err := h.paymentService.GetPayment(req.PaymentId)
	if err != nil {
		return nil, err
	}

	customerId := req.CustomerId

	if customerId != "admin" && payment.CustomerID.String() != customerId {
		return nil, errors.New("Access denied")
	}

	return &proto.GetPaymentResponse{
		PaymentId:     payment.ID.String(),
		OrderId:       payment.OrderID.String(),
		CustomerId:    payment.CustomerID.String(),
		TransactionId: payment.TransactionID,
		Amount:        payment.Amount,
		Status:        payment.Status,
	}, nil
}

func (h *paymentHandler) GetPaymentByOrderID(ctx context.Context, req *proto.GetPaymentByOrderIDRequest) (*proto.GetPaymentResponse, error) {
	payment, err := h.paymentService.GetPaymentByOrderID(req.OrderId)
	if err != nil {
		return nil, err
	}

	customerId := req.CustomerId

	if customerId != "admin" && payment.CustomerID.String() != customerId {
		return nil, errors.New("Access denied")
	}

	return &proto.GetPaymentResponse{
		PaymentId:     payment.ID.String(),
		OrderId:       payment.OrderID.String(),
		CustomerId:    payment.CustomerID.String(),
		TransactionId: payment.TransactionID,
		Amount:        payment.Amount,
		Status:        payment.Status,
	}, nil
}
