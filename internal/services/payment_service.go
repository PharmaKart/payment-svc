package services

import (
	"context"

	"github.com/PharmaKart/payment-svc/internal/models"
	"github.com/PharmaKart/payment-svc/internal/proto"
	"github.com/PharmaKart/payment-svc/internal/repositories"
	"github.com/PharmaKart/payment-svc/pkg/config"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/checkout/session"
)

type StripeResponse struct {
	URL string
}

type PaymentService interface {
	GeneratePaymentURL(orderId string, customerId string) (StripeResponse, error)
	StorePayment(payment *models.Payment) (string, error)
	RefundPayment(transactionId string) error
	GetPaymentByTransactionID(transactionID string) (*models.Payment, error)
	GetPayment(paymentID string) (*models.Payment, error)
	GetPaymentByOrderID(orderID string) (*models.Payment, error)
}

type paymentService struct {
	paymentRepo repositories.PaymentRepository
	orderClient proto.OrderServiceClient
	cfg         *config.Config
}

func NewPaymentService(paymentRepo repositories.PaymentRepository, orderService *proto.OrderServiceClient, cfg *config.Config) PaymentService {
	return &paymentService{
		paymentRepo: paymentRepo,
		orderClient: *orderService,
		cfg:         cfg,
	}
}

func (s *paymentService) GeneratePaymentURL(orderID string, customerID string) (StripeResponse, error) {
	stripe.Key = s.cfg.StripeSecretKey

	order, err := s.orderClient.GetOrder(context.Background(), &proto.GetOrderRequest{
		OrderId:    orderID,
		CustomerId: "admin",
	})
	if err != nil {
		return StripeResponse{}, err
	}

	lineItems := []*stripe.CheckoutSessionLineItemParams{}

	for _, item := range order.Items {
		lineItems = append(lineItems, &stripe.CheckoutSessionLineItemParams{
			PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
				Currency: stripe.String("cad"),
				ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
					Name: stripe.String(item.ProductName),
				},
				UnitAmount: stripe.Int64(int64(item.Price * 100)),
			},
			Quantity: stripe.Int64(int64(item.Quantity)),
		})
	}

	params := &stripe.CheckoutSessionParams{
		SuccessURL: stripe.String(s.cfg.GatewayURL + "/orders/" + orderID),
		LineItems:  lineItems,
		Metadata: map[string]string{
			"order_id":    orderID,
			"customer_id": customerID,
		},
		Mode: stripe.String(string(stripe.CheckoutSessionModePayment)),
	}

	session, err := session.New(params)
	if err != nil {
		return StripeResponse{}, err
	}

	return StripeResponse{
		URL: session.URL,
	}, nil
}

func (s *paymentService) StorePayment(payment *models.Payment) (string, error) {
	err := s.paymentRepo.StorePayment(payment)
	if err != nil {
		return "", err
	}

	var status string
	if payment.Status == "succeeded" {
		status = "paid"
	} else {
		status = "payment_failed"
	}

	_, err = s.orderClient.UpdateOrderStatus(context.Background(), &proto.UpdateOrderStatusRequest{
		OrderId: payment.OrderID.String(),
		Status:  status,
	})
	if err != nil {
		return "", err
	}
	return "Payment stored successfully", nil
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
