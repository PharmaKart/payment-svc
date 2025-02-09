package handlers

import (
	"context"

	"github.com/PharmaKart/payment-svc/internal/models"
	"github.com/PharmaKart/payment-svc/internal/proto"
	"github.com/PharmaKart/payment-svc/internal/repositories"
	"github.com/PharmaKart/payment-svc/internal/services"
	"github.com/google/uuid"
)

type PaymentHandler interface {
	ProcessPayment(ctx context.Context, req *proto.ProcessPaymentRequest) (*proto.ProcessPaymentResponse, error)
	RefundPayment(ctx context.Context, req *proto.RefundPaymentRequest) (*proto.RefundPaymentResponse, error)
	GetPaymentByTransactionID(ctx context.Context, req *proto.GetPaymentByTransactionIDRequest) (*proto.GetPaymentResponse, error)
	GetPayment(ctx context.Context, req *proto.GetPaymentRequest) (*proto.GetPaymentResponse, error)
	GetPaymentByOrderID(ctx context.Context, req *proto.GetPaymentByOrderIDRequest) (*proto.GetPaymentResponse, error)
}

type paymentHandler struct {
	proto.UnimplementedPaymentServiceServer
	paymentService services.PaymentService
}

func NewPaymentHandler(paymentRepo repositories.PaymentRepository, orderClient *proto.OrderServiceClient) *paymentHandler {
	return &paymentHandler{
		paymentService: services.NewPaymentService(paymentRepo, orderClient),
	}
}

func (h *paymentHandler) ProcessPayment(ctx context.Context, req *proto.ProcessPaymentRequest) (*proto.ProcessPaymentResponse, error) {
	orderId, err := uuid.Parse(req.OrderId)
	customerId, err := uuid.Parse(req.CustomerId)
	payment := &models.Payment{
		OrderID:       orderId,
		CustomerID:    customerId,
		TransactionID: req.TransactionId,
		Amount:        req.Amount,
		Status:        req.Status,
	}
	transactionID, err := h.paymentService.ProcessPayment(payment)
	if err != nil {
		return nil, err
	}

	return &proto.ProcessPaymentResponse{TransactionId: transactionID}, nil
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

	return &proto.GetPaymentResponse{
		PaymentId:     payment.ID.String(),
		OrderId:       payment.OrderID.String(),
		CustomerId:    payment.CustomerID.String(),
		TransactionId: payment.TransactionID,
		Amount:        payment.Amount,
		Status:        payment.Status,
	}, nil
}
