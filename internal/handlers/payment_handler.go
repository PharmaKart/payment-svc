package handlers

import (
	"context"
	"fmt"

	"github.com/PharmaKart/payment-svc/internal/models"
	"github.com/PharmaKart/payment-svc/internal/proto"
	"github.com/PharmaKart/payment-svc/internal/repositories"
	"github.com/PharmaKart/payment-svc/internal/services"
	"github.com/PharmaKart/payment-svc/pkg/config"
	"github.com/PharmaKart/payment-svc/pkg/errors"
	"github.com/PharmaKart/payment-svc/pkg/utils"
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
		if appErr, ok := errors.IsAppError(err); ok {
			return &proto.GeneratePaymentURLResponse{
				Success: false,
				Error: &proto.Error{
					Type:    string(appErr.Type),
					Message: appErr.Message,
					Details: utils.ConvertMapToKeyValuePairs(appErr.Details),
				},
			}, nil
		}
		return &proto.GeneratePaymentURLResponse{
			Success: false,
			Error: &proto.Error{
				Type:    string(errors.InternalError),
				Message: "An unexpected error occurred",
			},
		}, nil
	}

	return &proto.GeneratePaymentURLResponse{
		Success: true,
		Url:     resp.URL,
	}, nil
}

func (h *paymentHandler) StorePayment(ctx context.Context, req *proto.StorePaymentRequest) (*proto.StorePaymentResponse, error) {
	orderId, err := uuid.Parse(req.OrderId)
	if err != nil {
		return &proto.StorePaymentResponse{
			Success: false,
			Error: &proto.Error{
				Type:    string(errors.ValidationError),
				Message: "Invalid order ID",
				Details: utils.ConvertMapToKeyValuePairs(map[string]string{"order_id": fmt.Sprintf("Invalid UUID: %s", req.OrderId)}),
			},
		}, nil
	}

	customerId, err := uuid.Parse(req.CustomerId)
	if err != nil {
		return &proto.StorePaymentResponse{
			Success: false,
			Error: &proto.Error{
				Type:    string(errors.ValidationError),
				Message: "Invalid customer ID",
				Details: utils.ConvertMapToKeyValuePairs(map[string]string{"customer_id": fmt.Sprintf("Invalid UUID: %s", req.CustomerId)}),
			},
		}, nil
	}

	payment := &models.Payment{
		OrderID:       orderId,
		CustomerID:    customerId,
		TransactionID: req.TransactionId,
		Amount:        req.Amount,
		Status:        req.Status,
	}
	message, err := h.paymentService.StorePayment(payment)
	if err != nil {
		if appErr, ok := errors.IsAppError(err); ok {
			return &proto.StorePaymentResponse{
				Success: false,
				Error: &proto.Error{
					Type:    string(appErr.Type),
					Message: appErr.Message,
					Details: utils.ConvertMapToKeyValuePairs(appErr.Details),
				},
			}, nil
		}
		return &proto.StorePaymentResponse{
			Success: false,
			Error: &proto.Error{
				Type:    string(errors.InternalError),
				Message: "An unexpected error occurred",
			},
		}, nil
	}

	return &proto.StorePaymentResponse{
		Success: true,
		Message: message,
	}, nil
}

func (h *paymentHandler) RefundPayment(ctx context.Context, req *proto.RefundPaymentRequest) (*proto.RefundPaymentResponse, error) {
	err := h.paymentService.RefundPayment(req.TransactionId)
	if err != nil {
		if appErr, ok := errors.IsAppError(err); ok {
			return &proto.RefundPaymentResponse{
				Success: false,
				Error: &proto.Error{
					Type:    string(appErr.Type),
					Message: appErr.Message,
					Details: utils.ConvertMapToKeyValuePairs(appErr.Details),
				},
			}, nil
		}
		return &proto.RefundPaymentResponse{
			Success: false,
			Error: &proto.Error{
				Type:    string(errors.InternalError),
				Message: "An unexpected error occurred",
			},
		}, nil
	}

	return &proto.RefundPaymentResponse{
		Success: true,
	}, nil
}

func (h *paymentHandler) GetPaymentByTransactionID(ctx context.Context, req *proto.GetPaymentByTransactionIDRequest) (*proto.GetPaymentResponse, error) {
	payment, err := h.paymentService.GetPaymentByTransactionID(req.TransactionId)
	if err != nil {
		if appErr, ok := errors.IsAppError(err); ok {
			return &proto.GetPaymentResponse{
				Success: false,
				Error: &proto.Error{
					Type:    string(appErr.Type),
					Message: appErr.Message,
					Details: utils.ConvertMapToKeyValuePairs(appErr.Details),
				},
			}, nil
		}

		return &proto.GetPaymentResponse{
			Success: false,
			Error: &proto.Error{
				Type:    string(errors.InternalError),
				Message: "An unexpected error occurred",
			},
		}, nil
	}

	customerId := req.CustomerId

	if customerId != "admin" && payment.CustomerID.String() != customerId {
		return &proto.GetPaymentResponse{
			Success: false,
			Error: &proto.Error{
				Type:    string(errors.AuthError),
				Message: "You are not authorized to view this payment",
			},
		}, nil
	}

	return &proto.GetPaymentResponse{
		Success:       true,
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
		if appErr, ok := errors.IsAppError(err); ok {
			return &proto.GetPaymentResponse{
				Success: false,
				Error: &proto.Error{
					Type:    string(appErr.Type),
					Message: appErr.Message,
					Details: utils.ConvertMapToKeyValuePairs(appErr.Details),
				},
			}, nil
		}

		return &proto.GetPaymentResponse{
			Success: false,
			Error: &proto.Error{
				Type:    string(errors.InternalError),
				Message: "An unexpected error occurred",
			},
		}, nil
	}

	customerId := req.CustomerId

	if customerId != "admin" && payment.CustomerID.String() != customerId {
		return &proto.GetPaymentResponse{
			Success: false,
			Error: &proto.Error{
				Type:    string(errors.AuthError),
				Message: "You are not authorized to view this payment",
			},
		}, nil
	}

	return &proto.GetPaymentResponse{
		Success:       true,
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
		if appErr, ok := errors.IsAppError(err); ok {
			return &proto.GetPaymentResponse{
				Success: false,
				Error: &proto.Error{
					Type:    string(appErr.Type),
					Message: appErr.Message,
					Details: utils.ConvertMapToKeyValuePairs(appErr.Details),
				},
			}, nil
		}

		return &proto.GetPaymentResponse{
			Success: false,
			Error: &proto.Error{
				Type:    string(errors.InternalError),
				Message: "An unexpected error occurred",
			},
		}, nil
	}

	customerId := req.CustomerId

	if customerId != "admin" && payment.CustomerID.String() != customerId {
		return &proto.GetPaymentResponse{
			Success: false,
			Error: &proto.Error{
				Type:    string(errors.AuthError),
				Message: "You are not authorized to view this payment",
			},
		}, nil
	}

	return &proto.GetPaymentResponse{
		Success:       true,
		PaymentId:     payment.ID.String(),
		OrderId:       payment.OrderID.String(),
		CustomerId:    payment.CustomerID.String(),
		TransactionId: payment.TransactionID,
		Amount:        payment.Amount,
		Status:        payment.Status,
	}, nil
}
