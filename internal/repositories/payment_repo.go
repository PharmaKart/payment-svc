package repositories

import (
	"github.com/PharmaKart/payment-svc/internal/models"
	"gorm.io/gorm"
)

type PaymentRepository interface {
	CreatePayment(payment *models.Payment) (string, error)
	GetPaymentByOrderID(orderID string) (*models.Payment, error)
	GetPaymentByTransactionID(transactionID string) (*models.Payment, error)
	GetPayment(paymentID string) (*models.Payment, error)
	UpdatePaymentStatus(orderID string, status string) error
}

type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{db}
}

func (r *paymentRepository) CreatePayment(payment *models.Payment) (string, error) {
	// return r.db.Create(payment).Error
	if err := r.db.Create(payment).Error; err != nil {
		return "", err
	}

	return payment.TransactionID, nil
}

func (r *paymentRepository) GetPaymentByOrderID(orderID string) (*models.Payment, error) {
	var payment models.Payment
	err := r.db.Where("order_id = ?", orderID).First(&payment).Error
	return &payment, err
}

func (r *paymentRepository) GetPaymentByTransactionID(transactionID string) (*models.Payment, error) {
	var payment models.Payment
	err := r.db.Where("transaction_id = ?", transactionID).First(&payment).Error
	return &payment, err
}

func (r *paymentRepository) GetPayment(paymentID string) (*models.Payment, error) {
	var payment models.Payment
	err := r.db.Where("id = ?", paymentID).First(&payment).Error
	return &payment, err
}

func (r *paymentRepository) UpdatePaymentStatus(orderID string, status string) error {
	return r.db.Model(&models.Payment{}).Where("order_id = ?", orderID).Update("status", status).Error
}
