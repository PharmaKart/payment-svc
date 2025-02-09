package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/PharmaKart/payment-svc/internal/models"
	"github.com/PharmaKart/payment-svc/internal/proto"
	"github.com/PharmaKart/payment-svc/internal/repositories"
)

type PaymentService interface {
	ProcessPayment(payment *models.Payment) (string, error)
	RefundPayment(transactionId string) error
	GetPaymentByTransactionID(transactionID string) (*models.Payment, error)
	GetPayment(paymentID string) (*models.Payment, error)
	GetPaymentByOrderID(orderID string) (*models.Payment, error)
}

type paymentService struct {
	paymentRepo repositories.PaymentRepository
	orderClient proto.OrderServiceClient
}

func NewPaymentService(paymentRepo repositories.PaymentRepository, orderService *proto.OrderServiceClient) PaymentService {
	return &paymentService{
		paymentRepo: paymentRepo,
		orderClient: *orderService,
	}
}

func (s *paymentService) ProcessPayment(payment *models.Payment) (string, error) {
	payment, err := s.paymentRepo.GetPaymentByTransactionID(payment.TransactionID)
	if err == nil {
		return "", errors.New(fmt.Sprintf("payment already exists with transaction id %s", payment.TransactionID))
	}

	payment, err = s.paymentRepo.GetPaymentByOrderID(payment.OrderID.String())
	if err == nil {
		return "", errors.New(fmt.Sprintf("payment already exists for order id %s", payment.OrderID.String()))
	}

	transaction_id, err := s.paymentRepo.CreatePayment(payment)
	if err != nil {
		return "", err
	}
	return transaction_id, nil
}

func (s *paymentService) RefundPayment(transactionId string) error {
	payment, err := s.paymentRepo.GetPaymentByTransactionID(transactionId)
	if err != nil {
		return err
	}

	err = s.paymentRepo.UpdatePaymentStatus(payment.OrderID.String(), "refunded")
	if err != nil {
		return err
	}

	_, err = s.orderClient.UpdateOrderStatus(context.Background(), &proto.UpdateOrderStatusRequest{
		OrderId: payment.OrderID.String(),
		Status:  "refunded",
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *paymentService) GetPaymentByTransactionID(transactionID string) (*models.Payment, error) {
	payment, err := s.paymentRepo.GetPaymentByTransactionID(transactionID)
	if err != nil {
		return nil, err
	}
	return payment, nil
}

func (s *paymentService) GetPayment(paymentID string) (*models.Payment, error) {
	payment, err := s.paymentRepo.GetPayment(paymentID)
	if err != nil {
		return nil, err
	}
	return payment, nil
}

func (s *paymentService) GetPaymentByOrderID(orderID string) (*models.Payment, error) {
	payment, err := s.paymentRepo.GetPaymentByOrderID(orderID)
	if err != nil {
		return nil, err
	}
	return payment, nil
}
